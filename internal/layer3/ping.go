package layer3

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type PingResult struct {
	Target      string        `json:"target"`
	IP          string        `json:"ip"`
	Transmitted int           `json:"transmitted"`
	Received    int           `json:"received"`
	Loss        float64       `json:"loss_percent"`
	RTTMin      time.Duration `json:"rtt_min"`
	RTTAvg      time.Duration `json:"rtt_avg"`
	RTTMax      time.Duration `json:"rtt_max"`
	PerPacket   []PacketInfo  `json:"per_packet"`
}

type PacketInfo struct {
	Seq    int           `json:"seq"`
	RTT    time.Duration `json:"rtt"`
	Status string        `json:"status"`
}

func NewPingCommand() *cobra.Command {
	var count int
	var timeout int
	var interval int
	var ipv6Mode bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "ping [host]...",
		Short: "Send ICMP echo requests to one or more hosts (Layer 3)",
		Long: `Performs a network reachability test using ICMP echo requests (ping).
Supports both IPv4 and IPv6.
Multiple hosts can be specified and will be pinged in parallel.
Each response includes per-packet information, round-trip timing, and overall statistics.
Results can be returned in plain text or JSON format, making it suitable for scripting and cross-platform integration.

Arguments:
  host  - One or more IP addresses or hostnames to ping (space-separated)`,
		Example: `
  netanalyzer ping 8.8.8.8 --count 5 --json
  netanalyzer ping example.com --ipv6
  netanalyzer ping host1.com host2.com`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var wg sync.WaitGroup
			results := make([]PingResult, len(args))
			mutex := sync.Mutex{}

			wg.Add(len(args))
			for i, target := range args {
				go func(i int, target string) {
					defer wg.Done()
					r, _ := runPingSingle(target, count, timeout, interval, ipv6Mode, jsonOutput)
					mutex.Lock()
					results[i] = r
					mutex.Unlock()
				}(i, target)
			}
			wg.Wait()

			if jsonOutput {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(results)
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&count, "count", 4, "Number of ICMP echo requests")
	cmd.Flags().IntVar(&timeout, "timeout", 2000, "Timeout per request in milliseconds")
	cmd.Flags().IntVar(&interval, "interval", 1000, "Interval between pings in milliseconds")
	cmd.Flags().BoolVar(&ipv6Mode, "ipv6", false, "Use IPv6 instead of IPv4")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")

	return cmd
}

func runPingSingle(target string, count int, timeoutMs int, intervalMs int, ipv6Mode bool, jsonOutput bool) (PingResult, error) {
	var network, listenAddr string
	var protocol int
	if ipv6Mode {
		network = "ip6:ipv6-icmp"
		listenAddr = "::"
		protocol = 58
	} else {
		network = "ip4:icmp"
		listenAddr = "0.0.0.0"
		protocol = 1
	}

	addr, err := net.ResolveIPAddr("ip", target)
	if err != nil {
		return PingResult{Target: target}, fmt.Errorf("resolve error: %w", err)
	}

	conn, err := icmp.ListenPacket(network, listenAddr)
	if err != nil {
		return PingResult{Target: target}, fmt.Errorf("listen error: %w", err)
	}
	defer conn.Close()

	id := os.Getpid() & 0xffff
	sent := 0
	received := 0
	var rttTimes []time.Duration
	var packets []PacketInfo

	for i := 0; i < count; i++ {
		seq := i + 1
		msg := icmp.Message{
			Type: getICMPType(ipv6Mode),
			Code: 0,
			Body: &icmp.Echo{
				ID:   id,
				Seq:  seq,
				Data: []byte("NETANALYZER-PING"),
			},
		}

		data, err := msg.Marshal(nil)
		if err != nil {
			return PingResult{Target: target}, fmt.Errorf("marshal error: %w", err)
		}

		start := time.Now()
		_, err = conn.WriteTo(data, addr)
		sent++
		if err != nil {
			packets = append(packets, PacketInfo{Seq: seq, Status: "send error"})
			continue
		}

		reply := make([]byte, 1500)
		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutMs) * time.Millisecond))
		n, peer, err := conn.ReadFrom(reply)
		duration := time.Since(start)

		if err != nil {
			packets = append(packets, PacketInfo{Seq: seq, Status: "timeout"})
			if !jsonOutput {
				fmt.Printf("%s: Request timeout for icmp_seq %d\n", target, seq)
			}
		} else {
			rm, err := icmp.ParseMessage(protocol, reply[:n])
			if err != nil {
				packets = append(packets, PacketInfo{Seq: seq, Status: "parse error"})
				continue
			}

			switch body := rm.Body.(type) {
			case *icmp.Echo:
				if rm.Type == getICMPEchoReplyType(ipv6Mode) && body.ID == id {
					received++
					rttTimes = append(rttTimes, duration)
					packets = append(packets, PacketInfo{Seq: seq, RTT: duration, Status: "ok"})
					if !jsonOutput {
						fmt.Printf("%s: %d bytes from %s: icmp_seq=%d time=%v\n", target, n, peer.String(), seq, duration)
					}
				}
			default:
				packets = append(packets, PacketInfo{Seq: seq, Status: "unexpected reply"})
			}
		}
		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}

	loss := float64(sent-received) / float64(sent) * 100

	var min, max, total time.Duration
	if len(rttTimes) > 0 {
		min, max = rttTimes[0], rttTimes[0]
		for _, rtt := range rttTimes {
			total += rtt
			if rtt < min {
				min = rtt
			}
			if rtt > max {
				max = rtt
			}
		}
	}

	avg := time.Duration(0)
	if len(rttTimes) > 0 {
		avg = total / time.Duration(len(rttTimes))
	}

	res := PingResult{
		Target:      target,
		IP:          addr.String(),
		Transmitted: sent,
		Received:    received,
		Loss:        loss,
		RTTMin:      min,
		RTTAvg:      avg,
		RTTMax:      max,
		PerPacket:   packets,
	}

	if !jsonOutput {
		fmt.Printf("\n--- %s ping statistics ---\n", target)
		fmt.Printf("%d packets transmitted, %d received, %.1f%% packet loss\n", sent, received, loss)
		if len(rttTimes) > 0 {
			fmt.Printf("rtt min/avg/max = %v/%v/%v\n", min, avg, max)
		}
	}

	return res, nil
}

func getICMPType(ipv6Enabled bool) icmp.Type {
	if ipv6Enabled {
		return ipv6.ICMPTypeEchoRequest
	}
	return ipv4.ICMPTypeEcho
}

func getICMPEchoReplyType(ipv6Enabled bool) icmp.Type {
	if ipv6Enabled {
		return ipv6.ICMPTypeEchoReply
	}
	return ipv4.ICMPTypeEchoReply
}
