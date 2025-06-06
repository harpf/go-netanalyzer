// internal/layer4/iperf.go
package layer4

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

func NewIperfCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "iperf [host]",
		Short: "Run an iperf3 bandwidth test against a target (Layer 4)",
		Long: `Runs an iperf3 test using the embedded binary to test network throughput.
Requires an iperf3 server to be running on the target system.

Arguments:
  host - Target hostname or IP`,
		Example: `
  netanalyzer iperf 192.168.1.1
  netanalyzer iperf example.com`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			RunIperfTest(host)
		},
	}
}

func RunIperfTest(host string) {
	exe := "win-iperf3.exe"
	if runtime.GOOS != "windows" {
		exe = "iperf3"
	}

	path := filepath.Join("bin", exe)
	cmd := exec.Command(path, "-c", host, "-J")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("iperf3 execution failed: %v\n", err)
		return
	}

	fmt.Println(string(output))
}
