package layer3

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func NewPingCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ping [host]",
		Short: "Send ICMP echo requests to a host (Layer 3)",
		Long: `Performs a basic network reachability test using raw ICMP echo requests (ping).
This test sends 4 pings to the specified host and returns basic timing information.

Requires administrative/root privileges due to raw socket usage.

Arguments:
  host  - Target IP address or hostname`,
		Example: `
  netanalyzer ping 8.8.8.8
  netanalyzer ping example.com`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := args[0]
			err := RunPingTest(target, 4)
			if err != nil {
				fmt.Println("Error:", err)
			}
		},
	}
}

func RunPingTest(target string, count int) error {
	addr, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		return fmt.Errorf("resolve error: %w", err)
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}
	defer conn.Close()

	for i := 0; i < count; i++ {
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  i,
				Data: []byte("HELLO-PING"),
			},
		}

		bytes, err := msg.Marshal(nil)
		if err != nil {
			return fmt.Errorf("marshal error: %w", err)
		}

		start := time.Now()
		_, err = conn.WriteTo(bytes, addr)
		if err != nil {
			return fmt.Errorf("write error: %w", err)
		}

		buffer := make([]byte, 1500)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, peer, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Request timeout")
			continue
		}

		duration := time.Since(start)
		resp, err := icmp.ParseMessage(1, buffer[:n])
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		if resp.Type == ipv4.ICMPTypeEchoReply {
			fmt.Printf("Reply from %s: time=%v\n", peer, duration)
		} else {
			fmt.Printf("Unexpected ICMP type: %v\n", resp.Type)
		}

		time.Sleep(1 * time.Second)
	}
	return nil
}
