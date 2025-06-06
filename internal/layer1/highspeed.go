package layer1

import (
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/spf13/cobra"
)

func NewHighSpeedCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "highspeed [host] [community] [ifIndex]",
		Short: "Check high-speed interface via SNMP (Layer 1)",
		Long: `Queries the SNMP OID ifHighSpeed (1.3.6.1.2.1.31.1.1.1.15.X) for the given interface index.

Returns the interface speed in megabits per second (Mbps), allowing values greater than 4 Gbps.

Arguments:
  host       - IP address or hostname of the SNMP device
  community  - SNMP community string
  ifIndex    - Interface index (e.g., 1, 2, 3...)`,
		Example: `
  netanalyzer highspeed 192.168.1.1 public 2`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			community := args[1]
			ifIndex := args[2]
			err := CheckHighSpeed(host, community, ifIndex)
			if err != nil {
				fmt.Println("Error:", err)
			}
		},
	}
	return cmd
}

func CheckHighSpeed(host string, community string, ifIndex string) error {
	oid := fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.15.%s", ifIndex) // ifHighSpeed

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

	result, err := params.Get([]string{oid})
	if err != nil {
		return fmt.Errorf("SNMP get error: %w", err)
	}

	if len(result.Variables) == 0 {
		return fmt.Errorf("no result returned for OID %s", oid)
	}

	speed := result.Variables[0].Value
	fmt.Printf("High Speed for interface %s: %v Mbit/s\n", ifIndex, speed)
	return nil
}
