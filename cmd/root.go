package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netanalyzer",
	Short: "A CLI tool to analyze all OSI layers",
	Long:  "NetAnalyzer is a diagnostic tool for performing network analysis across all OSI layers.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func AddSubCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
