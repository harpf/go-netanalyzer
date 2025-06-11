package layer4

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

type IperfResult struct {
	Start     interface{} `json:"start"`
	End       interface{} `json:"end"`
	Success   bool        `json:"success"`
	Error     string      `json:"error,omitempty"`
	RawOutput string      `json:"-"`
}

func NewIperfCommand() *cobra.Command {
	var port int
	var duration int
	var udp bool
	var parallel int
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "iperf [host]",
		Short: "Run an iperf3 bandwidth test against a target (Layer 4)",
		Long: `Runs an iperf3 test using the embedded binary to test network throughput.
Requires an iperf3 server to be running on the target system.`,
		Example: `
  netanalyzer iperf 192.168.1.1
  netanalyzer iperf example.com --port 5201 --udp --duration 10 --parallel 4 --json`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			result := RunIperfTest(host, port, duration, udp, parallel)

			if jsonOutput {
				out, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(out))
			} else {
				if result.Success {
					fmt.Println("iPerf3 test succeeded")
				} else {
					fmt.Printf("iPerf3 test failed: %s\n", result.Error)
				}
			}
		},
	}

	cmd.Flags().IntVar(&port, "port", 5201, "Port to use for iPerf3")
	cmd.Flags().IntVar(&duration, "duration", 10, "Test duration in seconds")
	cmd.Flags().BoolVar(&udp, "udp", false, "Use UDP instead of TCP")
	cmd.Flags().IntVar(&parallel, "parallel", 1, "Number of parallel client streams")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output result in JSON format")
	return cmd
}

func RunIperfTest(host string, port int, duration int, udp bool, parallel int) IperfResult {
	exe := "win-iperf3.exe"
	if runtime.GOOS != "windows" {
		exe = "iperf3"
	}
	path := filepath.Join("bin", exe)

	args := []string{"-c", host, "-p", fmt.Sprint(port), "-t", fmt.Sprint(duration), "-P", fmt.Sprint(parallel), "-J"}
	if udp {
		args = append(args, "-u")
	}

	cmd := exec.Command(path, args...)
	start := time.Now()
	output, err := cmd.CombinedOutput()
	elapsed := time.Since(start)

	result := IperfResult{
		RawOutput: string(output),
		Success:   err == nil,
	}
	if err != nil {
		result.Error = err.Error()
		return result
	}

	err = json.Unmarshal(output, &result)
	if err != nil {
		result.Error = "Failed to parse iPerf3 JSON output"
		result.Success = false
	}
	fmt.Printf("Execution time: %v\n", elapsed)
	return result
}
