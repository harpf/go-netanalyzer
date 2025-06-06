package layer4

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/harpf/go-netanalyzer/internal/utils"
	"github.com/spf13/cobra"
)

func NewTCPServiceCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tcpservice [host] [port]",
		Short: "Check if a TCP service is responsive (Layer 4)",
		Long: `Attempts to connect to a TCP service and confirms responsiveness.
Useful for checking if a service is accepting connections without reading banners.

Arguments:
  host  - Target hostname or IP
  port  - TCP port to check`,
		Example: `
  netanalyzer tcpservice 192.168.1.1 443
  netanalyzer tcpservice example.com 22`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			port, _ := strconv.Atoi(args[1])
			RunTCPServiceCheck(host, port)
		},
	}
}

func RunTCPServiceCheck(host string, port int) {
	address := utils.FormatAddress(host, port)

	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		fmt.Printf("TCP service at %s is not available: %v\n", address, err)
		return
	}
	defer conn.Close()

	fmt.Printf("TCP service at %s is responsive.\n", address)
}
