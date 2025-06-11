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

type HopResult struct {
	TTL        int           `json:"ttl"`
	Address    string        `json:"address"`
	Hostname   string        `json:"hostname"`
	Duration   time.Duration `json:"duration"`
	Success    bool          `json:"success"`
	Type       string        `json:"type"`
	ICMPType   string        `json:"icmp_type"`
	RawMessage []byte        `json:"raw_message,omitempty"`
}

type TracerouteResult struct {
	Target string      `json:"target"`
	Hops   []HopResult `json:"hops"`
}

func NewTracerouteCommand() *cobra.Command {
	var maxHops int
	var ipv6Mode bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "traceroute [host]...",
		Short: "Perform a traceroute to one or more hosts (Layer 3)",
		Long: `Sends ICMP echo requests with increasing TTL to trace the path to the destination host.
Supports both IPv4 and IPv6. Requires administrative privileges to open raw sockets.

Each host is traced in parallel. Results can be printed in plain text or JSON format.
Each hop includes hostname, RTT, ICMP type and response details.`,
		Example: `
  netanalyzer traceroute 8.8.8.8
  netanalyzer traceroute example.com --ipv6
  netanalyzer traceroute 1.1.1.1 8.8.8.8 --json`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			results := make([]TracerouteResult, len(args))
			mutex := sync.Mutex{}

			wg.Add(len(args))
			for i, target := range args {
				go func(i int, target string) {
					defer wg.Done()
					hops := RunTraceroute(target, maxHops, ipv6Mode)
					mutex.Lock()
					results[i] = TracerouteResult{Target: target, Hops: hops}
					mutex.Unlock()
				}(i, target)
			}
			wg.Wait()

			if jsonOutput {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				_ = enc.Encode(results)
			} else {
				for _, res := range results {
					fmt.Printf("\nTraceroute to %s:\n", res.Target)
					for _, hop := range res.Hops {
						if hop.Success {
							fmt.Printf("%2d  %-40s  %v\n", hop.TTL, fmt.Sprintf("%s (%s)", hop.Hostname, hop.Address), hop.Duration)
						} else {
							fmt.Printf("%2d  * * *\n", hop.TTL)
						}
					}
				}
			}
		},
	}

	cmd.Flags().IntVar(&maxHops, "maxhops", 30, "Maximum number of hops")
	cmd.Flags().BoolVar(&ipv6Mode, "ipv6", false, "Use IPv6 instead of IPv4")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")
	return cmd
}

func RunTraceroute(target string, maxHops int, ipv6Mode bool) []HopResult {
	var (
		icmpType     icmp.Type
		replyType    icmp.Type
		network      string
		listenAddr   string
		proto        int
		protocolType string
	)

	if ipv6Mode {
		icmpType = ipv6.ICMPTypeEchoRequest
		replyType = ipv6.ICMPTypeEchoReply
		network = "ip6:ipv6-icmp"
		listenAddr = "::"
		proto = 58
		protocolType = "IPv6"
	} else {
		icmpType = ipv4.ICMPTypeEcho
		replyType = ipv4.ICMPTypeEchoReply
		network = "ip4:icmp"
		listenAddr = "0.0.0.0"
		proto = 1
		protocolType = "IPv4"
	}

	ipAddr, err := net.ResolveIPAddr("ip", target)
	if err != nil {
		return []HopResult{{TTL: 1, Address: "resolve error", Success: false, Type: protocolType}}
	}

	var results []HopResult
	for ttl := 1; ttl <= maxHops; ttl++ {
		conn, err := icmp.ListenPacket(network, listenAddr)
		if err != nil {
			results = append(results, HopResult{TTL: ttl, Address: "listen error", Success: false, Type: protocolType})
			continue
		}

		if ipv6Mode {
			_ = conn.IPv6PacketConn().SetHopLimit(ttl)
		} else {
			_ = conn.IPv4PacketConn().SetTTL(ttl)
		}
		_ = conn.SetDeadline(time.Now().Add(2 * time.Second))

		msg := icmp.Message{
			Type: icmpType,
			Code: 0,
			Body: &icmp.Echo{
				ID:   ttl,
				Seq:  ttl,
				Data: []byte("traceroute"),
			},
		}

		bytes, _ := msg.Marshal(nil)
		start := time.Now()
		_, _ = conn.WriteTo(bytes, ipAddr)
		buffer := make([]byte, 1500)
		n, peer, err := conn.ReadFrom(buffer)
		duration := time.Since(start)
		_ = conn.Close()

		hostname := ""
		if err == nil {
			addrs, _ := net.LookupAddr(peer.String())
			if len(addrs) > 0 {
				hostname = addrs[0]
			}
		}

		if err != nil {
			results = append(results, HopResult{TTL: ttl, Address: "*", Hostname: "", Success: false, Type: protocolType})
			continue
		}

		resp, _ := icmp.ParseMessage(proto, buffer[:n])
		results = append(results, HopResult{
			TTL:      ttl,
			Address:  peer.String(),
			Hostname: hostname,
			Duration: duration,
			Success:  true,
			Type:     protocolType,
			ICMPType: func() string {
				switch t := resp.Type.(type) {
				case ipv4.ICMPType:
					return t.String()
				case ipv6.ICMPType:
					return t.String()
				default:
					return fmt.Sprintf("%v", t)
				}
			}(),
			RawMessage: buffer[:n],
		})

		if resp.Type == replyType {
			break
		}
	}
	return results
}
