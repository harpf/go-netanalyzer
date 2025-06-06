package main

import (
	"github.com/yourname/go-netanalyzer/cmd"
	"github.com/yourname/go-netanalyzer/internal/layer3"
)

func main() {
	cmd.AddSubCommand(layer3.NewPingCommand())
	cmd.Execute()
}
