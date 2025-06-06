package layer3

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

func NewDNSLookupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "dnslookup [hostname]",
		Short: "Resolve a hostname using DNS (Layer 3)",
		Long: `Resolves the specified hostname using the system's DNS configuration.
Returns one or more IP addresses associated with the hostname.

Arguments:
  hostname - The hostname to resolve (e.g. google.com)`,
		Example: `
  netanalyzer dnslookup google.com
  netanalyzer dnslookup example.org`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hostname := args[0]
			RunDNSLookup(hostname)
		},
	}
}

func RunDNSLookup(hostname string) {
	ips, err := net.LookupHost(hostname)
	if err != nil {
		fmt.Println("DNS lookup failed:", err)
		return
	}

	fmt.Printf("Resolved IPs for %s:\n", hostname)
	for _, ip := range ips {
		fmt.Printf("  %s\n", ip)
	}
}
