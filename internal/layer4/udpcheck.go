package layer4

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func NewUDPCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "udpcheck [host] [port]",
		Short: "Send a UDP packet to a host and port (Layer 4)",
		Long: `Attempts to send a simple UDP message to a specific port.
Since UDP is connectionless, a successful send does not confirm service availability.

Arguments:
  host  - Target IP address or hostname
  port  - Destination UDP port number`,
		Example: `
  netanalyzer udpcheck 192.168.1.1 161
  netanalyzer udpcheck example.com 53`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			port, _ := strconv.Atoi(args[1])
			RunUDPCheck(host, port)
		},
	}
}

func RunUDPCheck(host string, port int) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("udp", addr, 2*time.Second)
	if err != nil {
		fmt.Println("UDP dial error:", err)
		return
	}
	defer conn.Close()

	message := []byte("ping")
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("UDP write error:", err)
	} else {
		fmt.Printf("UDP packet sent to %s\n", addr)
	}
}
