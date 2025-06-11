package layer4

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/harpf/go-netanalyzer/internal/utils"
	"github.com/spf13/cobra"
)

type PortScanResult struct {
	Host     string        `json:"host"`
	Port     int           `json:"port"`
	Open     bool          `json:"open"`
	Duration time.Duration `json:"duration_ms"`
	Error    string        `json:"error,omitempty"`
}

func NewTCPScanCommand() *cobra.Command {
	var timeout time.Duration
	var concurrency int
	var jsonOutput bool
	var showAll bool

	cmd := &cobra.Command{
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
  netanalyzer tcpscan [2001:db8::1] 80 443 --json
  netanalyzer tcpscan example.com 1 1024 --timeout 1s --concurrency 100`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			start, _ := strconv.Atoi(args[1])
			end, _ := strconv.Atoi(args[2])
			results := RunTCPScan(host, start, end, timeout, concurrency)

			if jsonOutput {
				filtered := results
				if !showAll {
					filtered = []PortScanResult{}
					for _, r := range results {
						if r.Open {
							filtered = append(filtered, r)
						}
					}
				}
				out, _ := json.MarshalIndent(filtered, "", "  ")
				fmt.Println(string(out))
			} else {
				for _, r := range results {
					if r.Open {
						fmt.Printf("Port %d open (%v)\n", r.Port, r.Duration)
					} else if showAll {
						fmt.Printf("Port %d closed (%v)\n", r.Port, r.Duration)
					}
				}
			}
		},
	}

	cmd.Flags().DurationVar(&timeout, "timeout", 500*time.Millisecond, "Connection timeout")
	cmd.Flags().IntVar(&concurrency, "concurrency", 100, "Max number of parallel scans")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	cmd.Flags().BoolVar(&showAll, "all", false, "Show closed ports as well")
	return cmd
}

func RunTCPScan(host string, startPort, endPort int, timeout time.Duration, concurrency int) []PortScanResult {
	fmt.Printf("Scanning TCP ports %d-%d on %s...\n", startPort, endPort, host)
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	results := make([]PortScanResult, endPort-startPort+1)
	resultLock := sync.Mutex{}

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(p int) {
			defer func() {
				<-sem
				wg.Done()
			}()
			address := utils.FormatAddress(host, p)
			start := time.Now()
			conn, err := net.DialTimeout("tcp", address, timeout)
			duration := time.Since(start)

			result := PortScanResult{
				Host:     host,
				Port:     p,
				Open:     err == nil,
				Duration: duration / time.Millisecond,
			}
			if err != nil {
				result.Error = err.Error()
			} else {
				_ = conn.Close()
			}

			resultLock.Lock()
			results[p-startPort] = result
			resultLock.Unlock()
		}(port)
	}
	wg.Wait()
	return results
}
