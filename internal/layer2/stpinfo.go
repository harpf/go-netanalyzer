package layer2

import (
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/spf13/cobra"
)

func NewStpInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stpinfo [host] [community]",
		Short: "Display STP port states via SNMP (Layer 2)",
		Long: `Queries the dot1dStpPortState OID (1.3.6.1.2.1.17.2.15) on SNMP-enabled devices
such as switches or bridges to return the spanning tree state of each port.

Common port state values:
  1 = Disabled
  2 = Blocking
  3 = Listening
  4 = Learning
  5 = Forwarding
  6 = Broken

Arguments:
  host       - IP address or hostname of the SNMP device
  community  - SNMP community string`,
		Example: `
  netanalyzer stpinfo 192.168.1.1 public
  netanalyzer stpinfo core-switch private`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			community := args[1]
			err := ReadStpInfo(host, community)
			if err != nil {
				fmt.Println("Error:", err)
			}
		},
	}
	return cmd
}

func ReadStpInfo(host, community string) error {
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

	results, err := params.WalkAll("1.3.6.1.2.1.17.2.15") // dot1dStpPortState
	if err != nil {
		return fmt.Errorf("SNMP walk error: %w", err)
	}

	fmt.Println("STP Port States:")
	for _, variable := range results {
		fmt.Printf("%s = %v\n", variable.Name, variable.Value)
	}
	return nil
}
