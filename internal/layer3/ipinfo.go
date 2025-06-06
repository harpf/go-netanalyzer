// internal/layer3/ipinfo.go
package layer3

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

func NewIPInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ipinfo [ip-address]",
		Short: "Inspect details of an IP address (Layer 3)",
		Long: `Parses and displays detailed information about an IP address.
Includes address family, whether it's loopback, multicast, private, and global unicast.

Arguments:
  ip-address - The IPv4 or IPv6 address to inspect`,
		Example: `
  netanalyzer ipinfo 192.168.1.1
  netanalyzer ipinfo 8.8.8.8`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			RunIPInfo(args[0])
		},
	}
}

func RunIPInfo(input string) {
	ip := net.ParseIP(input)
	if ip == nil {
		fmt.Println("Invalid IP address")
		return
	}

	fmt.Printf("Analyzing IP: %s\n", ip)
	fmt.Printf("  Address Family: %s\n", getIPFamily(ip))
	fmt.Printf("  Is Loopback: %v\n", ip.IsLoopback())
	fmt.Printf("  Is Multicast: %v\n", ip.IsMulticast())
	fmt.Printf("  Is Unspecified: %v\n", ip.IsUnspecified())
	fmt.Printf("  Is Global Unicast: %v\n", ip.IsGlobalUnicast())
	fmt.Printf("  Is Private: %v\n", isPrivateIP(ip))
}

func getIPFamily(ip net.IP) string {
	if ip.To4() != nil {
		return "IPv4"
	}
	return "IPv6"
}

func isPrivateIP(ip net.IP) bool {
	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7",
	}
	for _, cidr := range privateCIDRs {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
