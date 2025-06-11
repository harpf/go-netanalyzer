package layer3

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

func NewIPInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ipinfo [ip-addresses...]",
		Short: "Inspect details of one or more IP addresses (Layer 3)",
		Long: `Parses and displays detailed information about one or more IP addresses.
Supports IPv4 and IPv6. Output includes address family, loopback, multicast, private, and global unicast flags.

Arguments:
  ip-address - One or more IPv4 or IPv6 addresses`,
		Example: `
  netanalyzer ipinfo 8.8.8.8 192.168.1.1
  netanalyzer ipinfo ::1 fe80::1`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			RunIPInfoMultiple(args)
		},
	}
}

type IPInfo struct {
	IP              string `json:"ip"`
	AddressFamily   string `json:"address_family"`
	IsLoopback      bool   `json:"is_loopback"`
	IsMulticast     bool   `json:"is_multicast"`
	IsUnspecified   bool   `json:"is_unspecified"`
	IsGlobalUnicast bool   `json:"is_global_unicast"`
	IsPrivate       bool   `json:"is_private"`
}

func RunIPInfoMultiple(inputs []string) {
	var results []IPInfo

	for _, input := range inputs {
		ip := net.ParseIP(input)
		if ip == nil {
			fmt.Fprintf(os.Stderr, "Invalid IP address: %s\n", input)
			continue
		}

		result := IPInfo{
			IP:              ip.String(),
			AddressFamily:   getIPFamily(ip),
			IsLoopback:      ip.IsLoopback(),
			IsMulticast:     ip.IsMulticast(),
			IsUnspecified:   ip.IsUnspecified(),
			IsGlobalUnicast: ip.IsGlobalUnicast(),
			IsPrivate:       isPrivateIP(ip),
		}
		results = append(results, result)
	}

	jsonBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("Error serializing results:", err)
		return
	}

	fmt.Println(string(jsonBytes))
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
