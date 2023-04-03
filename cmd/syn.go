package cmd

import (
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/psssix/syn-scanner/internal/adapters"
	"github.com/psssix/syn-scanner/internal/scanners"
	"github.com/psssix/syn-scanner/pkg/producers"
	"github.com/psssix/syn-scanner/pkg/reporters"
	"github.com/psssix/syn-scanner/pkg/workers"
	"github.com/spf13/cobra"
)

var ErrEmptyTarget = errors.New("target is not specified for scanner")

const (
	defaultThreadCount = 64
	dialerTimeout      = 15
)

var synScanCmd = &cobra.Command{
	Use:   "syn <target>",
	Short: "Scan <target> using syn/ack-scanner",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args[0]) == 0 {
			return ErrEmptyTarget
		} else {
			return nil
		}
	},
	RunE: scan,
}

func init() {
	scanCmd.AddCommand(synScanCmd)
	synScanCmd.Flags().IntP("threads", "t", defaultThreadCount, "number of scan threads")
}

func scan(cmd *cobra.Command, args []string) error {
	target := args[0]
	threads, err := cmd.Flags().GetInt("threads")
	if err != nil {
		return errors.Wrap(err, "error while parsing threads flag")
	}

	scanners.NewScanner(
		producers.NewProducer(producers.MinPortNumber, producers.MinPortNumber),
		workers.NewSynAckScanner(&net.Dialer{Timeout: dialerTimeout * time.Second}),
		reporters.NewReporter(adapters.Printer{}),
	)(target, threads)

	return nil
}
