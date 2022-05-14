package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use:   executableName(),
	Short: "Simple SYN/ACK scanner",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleError(err)
	}
}

func executableName() string {
	n, err := os.Executable()
	if err != nil {
		panic("can't get current executable name")
	}
	return filepath.Base(n)
}

func handleError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
