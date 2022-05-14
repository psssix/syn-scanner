package cmd

import "github.com/spf13/cobra"

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan <target> with <scanner>.",
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
