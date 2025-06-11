package layer4

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/harpf/go-netanalyzer/internal/utils"
	"github.com/spf13/cobra"
)

type TCPBannerResult struct {
	Host      string   `json:"host"`
	Port      int      `json:"port"`
	Protocol  string   `json:"protocol"`
	Banner    []string `json:"banner"`
	Raw       string   `json:"raw,omitempty"`
	Duration  int64    `json:"duration_ms"`
	Success   bool     `json:"success"`
	Error     string   `json:"error,omitempty"`
	Timestamp string   `json:"timestamp"`
}

func NewTCPBannerCommand() *cobra.Command {
	var protocol string
	var timeout time.Duration
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "tcpbanner [host] [port]",
		Short: "Grab a TCP service banner (Layer 4)",
		Long: `Attempts to connect to a TCP service and read its initial response.
This is useful for identifying services like HTTP, FTP, SMTP, etc.

Arguments:
  host  - Target hostname or IP
  port  - Port to connect to`,
		Example: `
  netanalyzer tcpbanner 192.168.1.1 80
  netanalyzer tcpbanner mail.server.com 25
  netanalyzer tcpbanner example.com 80 --protocol http --json`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			port, _ := strconv.Atoi(args[1])
			result := RunTCPBannerGrab(host, port, protocol, timeout)

			if jsonOutput {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				_ = enc.Encode(result)
			} else {
				if result.Success {
					fmt.Printf("[%s:%d] (%s) Banner (%d ms):\n", result.Host, result.Port, result.Protocol, result.Duration)
					for _, line := range result.Banner {
						fmt.Println(line)
					}
				} else {
					fmt.Printf("Connection to %s:%d failed: %s\n", result.Host, result.Port, result.Error)
				}
			}
		},
	}

	cmd.Flags().StringVar(&protocol, "protocol", "", "Protocol to simulate (e.g. http, smtp, ftp)")
	cmd.Flags().DurationVar(&timeout, "timeout", 3*time.Second, "Connection timeout duration")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output result in JSON format")
	return cmd
}

func RunTCPBannerGrab(host string, port int, protocol string, timeout time.Duration) TCPBannerResult {
	address := utils.FormatAddress(host, port)
	start := time.Now()

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return TCPBannerResult{
			Host:      host,
			Port:      port,
			Protocol:  protocol,
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(timeout))

	// Optionale Protokollsimulation
	switch strings.ToLower(protocol) {
	case "http":
		_, _ = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	case "smtp":
		// Server sendet zuerst
	case "ftp":
		// Server sendet zuerst
	default:
		// kein aktives Protokollverhalten
	}

	r := bufio.NewReader(conn)
	var banner []string
	for i := 0; i < 5; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		banner = append(banner, strings.TrimSpace(line))
	}

	duration := time.Since(start).Milliseconds()
	if len(banner) == 0 {
		return TCPBannerResult{
			Host:      host,
			Port:      port,
			Protocol:  protocol,
			Duration:  duration,
			Success:   false,
			Error:     "No banner received or timeout",
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	return TCPBannerResult{
		Host:      host,
		Port:      port,
		Protocol:  protocol,
		Banner:    banner,
		Raw:       strings.Join(banner, "\n"),
		Duration:  duration,
		Success:   true,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
