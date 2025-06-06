package layer3

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

func NewPingCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ping [host]",
		Short: "Run Layer 3 ICMP Ping test",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := args[0]
			err := RunPingTest(target, 4)
			if err != nil {
				fmt.Println("Error:", err)
			}
		},
	}
}

func RunPingTest(target string, count int) error {
	pinger, err := ping.NewPinger(target)
	if err != nil {
		return fmt.Errorf("ping init failed: %w", err)
	}
	pinger.Count = count
	pinger.SetPrivileged(true)
	fmt.Printf("Pinging %s...\n", target)
	err = pinger.Run()
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	stats := pinger.Statistics()
	fmt.Printf("Packets: Sent = %d, Received = %d, Loss = %.2f%%\n",
		stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	return nil
}
