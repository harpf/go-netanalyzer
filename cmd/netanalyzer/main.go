package main

import (
	"github.com/harpf/go-netanalyzer/cmd"
	"github.com/harpf/go-netanalyzer/internal/layer2"
	"github.com/harpf/go-netanalyzer/internal/layer3"
)

func main() {
	cmd.AddSubCommand(layer2.NewMacTableCommand())
	cmd.AddSubCommand(layer3.NewPingCommand())
	cmd.Execute()
}
