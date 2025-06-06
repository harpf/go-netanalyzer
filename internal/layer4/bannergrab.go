package layer4

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/harpf/go-netanalyzer/internal/utils"
	"github.com/spf13/cobra"
)

func NewTCPBannerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tcpbanner [host] [port]",
		Short: "Grab a TCP service banner (Layer 4)",
		Long: `Attempts to connect to a TCP service and read its initial response.
This is useful for identifying services like HTTP, FTP, SMTP, etc.

Arguments:
  host  - Target hostname or IP
  port  - Port to connect to`,
		Example: `
  netanalyzer tcpbanner 192.168.1.1 80
  netanalyzer tcpbanner mail.server.com 25`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			port, _ := strconv.Atoi(args[1])
			RunTCPBannerGrab(host, port)
		},
	}
}

func RunTCPBannerGrab(host string, port int) {
	address := utils.FormatAddress(host, port)

	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	r := bufio.NewReader(conn)
	banner, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("No banner received or timeout.")
		return
	}

	fmt.Printf("Banner from %s: %s\n", address, strings.TrimSpace(banner))
}
