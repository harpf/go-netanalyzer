package layer3

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

func NewDNSLookupCommand() *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "dnslookup [hostnames...]",
		Short: "Resolve one or more hostnames using DNS (Layer 3)",
		Long: `Resolves one or more hostnames using the system's DNS configuration.
Returns IP addresses for each hostname provided.

Arguments:
  hostnames - One or more hostnames to resolve (e.g. google.com example.org)`,
		Example: `
  netanalyzer dnslookup google.com
  netanalyzer dnslookup google.com github.com
  netanalyzer dnslookup example.org --json`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			RunDNSLookup(args, outputJSON)
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output result as JSON")
	return cmd
}

type DNSResult struct {
	Hostname    string   `json:"hostname"`
	ResolvedIPs []string `json:"resolved_ips"`
	Count       int      `json:"count"`
	Error       string   `json:"error,omitempty"`
}

func RunDNSLookup(hostnames []string, asJSON bool) {
	results := []DNSResult{}

	for _, hostname := range hostnames {
		ips, err := net.LookupHost(hostname)
		result := DNSResult{
			Hostname: hostname,
			Count:    len(ips),
		}
		if err != nil {
			result.Error = err.Error()
		} else {
			result.ResolvedIPs = ips
		}
		results = append(results, result)
	}

	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(results)
		return
	}

	for _, res := range results {
		fmt.Printf("Hostname: %s\n", res.Hostname)
		if res.Error != "" {
			fmt.Printf("  Error: %s\n", res.Error)
		} else {
			fmt.Println("  Resolved IPs:")
			for _, ip := range res.ResolvedIPs {
				fmt.Printf("    %s\n", ip)
			}
			fmt.Printf("  Count: %d\n", res.Count)
		}
		fmt.Println()
	}
}
