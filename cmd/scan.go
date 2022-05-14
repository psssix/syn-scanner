package cmd

import "github.com/spf13/cobra"

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan <scanner> <target>.",
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
