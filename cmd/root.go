package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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
	executable, err := os.Executable()
	if err != nil {
		panic("can't get current executable name")
	}

	return filepath.Base(executable)
}

func handleError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
