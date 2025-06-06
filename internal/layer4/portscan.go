package layer4

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/harpf/go-netanalyzer/internal/utils"
	"github.com/spf13/cobra"
)

func NewTCPScanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tcpscan [host] [start-port] [end-port]",
		Short: "Perform a TCP port scan on a host (Layer 4)",
		Long: `Attempts to connect to TCP ports within the given range on the specified host.
Reports open ports based on successful connection attempts.

Arguments:
  host        - Target IP or hostname
  start-port  - Starting port number (e.g. 20)
  end-port    - Ending port number (e.g. 1024)`,
		Example: `
  netanalyzer tcpscan 192.168.1.1 20 80
  netanalyzer tcpscan [2001:db8::1] 80 443`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			start, _ := strconv.Atoi(args[1])
			end, _ := strconv.Atoi(args[2])
			RunTCPScan(host, start, end)
		},
	}
}

func RunTCPScan(host string, startPort, endPort int) {
	fmt.Printf("Scanning TCP ports %d-%d on %s...\n", startPort, endPort, host)
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			address := utils.FormatAddress(host, p)
			conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
			if err == nil {
				fmt.Printf("Port %d open\n", p)
				_ = conn.Close()
			}
		}(port)
	}
	wg.Wait()
}
