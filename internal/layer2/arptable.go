package layer2

import (
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/spf13/cobra"
)

func NewArpTableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "arptable [host] [community]",
		Short: "Display the ARP table via SNMP (Layer 2)",
		Long: `Performs an SNMP walk to retrieve the ARP table of a target device.
This command uses the ipNetToMediaPhysAddress OID (1.3.6.1.2.1.4.22.1.2) to
retrieve IP-to-MAC address mappings from routers, switches, or other SNMP-capable devices.

Arguments:
  host       - IP address or hostname of the SNMP device
  community  - SNMP community string (e.g., public)`,
		Example: `
  netanalyzer arptable 192.168.1.1 public
  netanalyzer arptable switch.local private`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			community := args[1]
			err := ReadArpTable(host, community)
			if err != nil {
				fmt.Println("Error:", err)
			}
		},
	}
	return cmd
}

func ReadArpTable(host, community string) error {
	params := &gosnmp.GoSNMP{
		Target:    host,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Retries:   2,
	}

	if err := params.Connect(); err != nil {
		return fmt.Errorf("SNMP connect error: %w", err)
	}
	defer params.Conn.Close()

	results, err := params.WalkAll("1.3.6.1.2.1.4.22.1.2") // ipNetToMediaPhysAddress
	if err != nil {
		return fmt.Errorf("SNMP walk error: %w", err)
	}

	fmt.Println("ARP Table Entries:")
	for _, variable := range results {
		fmt.Printf("%s = %v\n", variable.Name, variable.Value)
	}
	return nil
}
