package main

import (
	"github.com/harpf/go-netanalyzer/cmd"
	"github.com/harpf/go-netanalyzer/internal/layer1"
	"github.com/harpf/go-netanalyzer/internal/layer2"
	"github.com/harpf/go-netanalyzer/internal/layer3"
	"github.com/harpf/go-netanalyzer/internal/layer4"
)

func main() {
	// Layer 1 Commands
	cmd.AddSubCommand(layer1.NewLinkStatusCommand())
	cmd.AddSubCommand(layer1.NewInterfaceSpeedCommand())
	cmd.AddSubCommand(layer1.NewHighSpeedCommand())

	// Layer 2 Commands
	cmd.AddSubCommand(layer2.NewMacTableCommand())
	cmd.AddSubCommand(layer2.NewArpTableCommand())
	cmd.AddSubCommand(layer2.NewStpInfoCommand())

	// Layer 3 Commands
	cmd.AddSubCommand(layer3.NewPingCommand())
	cmd.AddSubCommand(layer3.NewTracerouteCommand())
	cmd.AddSubCommand(layer3.NewDNSLookupCommand())
	cmd.AddSubCommand(layer3.NewIPInfoCommand())

	// Layer 4 Commands
	cmd.AddSubCommand(layer4.NewTCPBannerCommand())
	cmd.AddSubCommand(layer4.NewTCPScanCommand())
	cmd.AddSubCommand(layer4.NewUDPCheckCommand())
	cmd.AddSubCommand(layer4.NewIperfCommand())
	cmd.AddSubCommand(layer4.NewTCPServiceCheckCommand())

	cmd.Execute()
}
