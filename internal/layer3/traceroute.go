// internal/layer3/traceroute.go
package layer3

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func NewTracerouteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "traceroute [host]",
		Short: "Perform a traceroute to a host (Layer 3)",
		Long: `Sends ICMP echo requests with increasing TTL to trace the path to the destination host.
Requires administrative privileges to open raw sockets.

Arguments:
  host  - Target IP address or hostname`,
		Example: `
  netanalyzer traceroute 8.8.8.8`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := args[0]
			RunTraceroute(target, 30)
		},
	}
}

func RunTraceroute(target string, maxHops int) {
	ipAddr, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		fmt.Println("Resolve error:", err)
		return
	}

	for ttl := 1; ttl <= maxHops; ttl++ {
		conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil {
			fmt.Println("Listen error:", err)
			return
		}
		pconn := conn.IPv4PacketConn()
		_ = pconn.SetTTL(ttl)
		_ = conn.SetDeadline(time.Now().Add(2 * time.Second))

		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
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

		if err != nil {
			fmt.Printf("%2d  * * *\n", ttl)
			_ = conn.Close()
			continue
		}

		resp, _ := icmp.ParseMessage(1, buffer[:n])
		fmt.Printf("%2d  %s  %v\n", ttl, peer.String(), duration)
		_ = conn.Close()

		if resp.Type == ipv4.ICMPTypeEchoReply {
			break
		}
	}
}
