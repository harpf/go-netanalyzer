package layer2

import (
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/spf13/cobra"
)

func NewMacTableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mactable [host] [community]",
		Short: "Display the MAC address table via SNMP (Layer 2)",
		Long: `Performs an SNMP walk on the dot1dTpFdbPort OID (1.3.6.1.2.1.17.4.3.1.2) to retrieve
MAC address to port mappings from an SNMP-capable switch or bridge.

Useful for identifying which MAC addresses are learned on which switch ports.

Arguments:
  host       - IP address or hostname of the SNMP device
  community  - SNMP community string (e.g., public)`,
		Example: `
  netanalyzer mactable 192.168.1.1 public
  netanalyzer mactable core-switch private`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			community := args[1]
			err := ReadMacTable(host, community)
			if err != nil {
				fmt.Println("Error:", err)
			}
		},
	}
	return cmd
}

func ReadMacTable(host string, community string) error {
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

	results, err := params.WalkAll("1.3.6.1.2.1.17.4.3.1.2")
	if err != nil {
		return fmt.Errorf("SNMP walk error: %w", err)
	}

	fmt.Println("MAC Table Entries:")
	for _, variable := range results {
		fmt.Printf("%s = %v\n", variable.Name, variable.Value)
	}
	return nil
}
