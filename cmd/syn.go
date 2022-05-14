package cmd

import (
	"errors"
	"github.com/psssix/syn-scanner/internal/adapters"
	"github.com/psssix/syn-scanner/internal/scanners"
	"github.com/psssix/syn-scanner/pkg/producers"
	"github.com/psssix/syn-scanner/pkg/reporters"
	"github.com/psssix/syn-scanner/pkg/workers"
	"github.com/spf13/cobra"
	"net"
	"time"
)

var synScanCmd = &cobra.Command{
	Use:   "syn <target>",
	Short: "Scan <target> using syn/ack-scanner",
	RunE:  scan,
}

func init() {
	scanCmd.AddCommand(synScanCmd)

}

var ErrEmptyTarget = errors.New("target is not specified for scanner")

const (
	defaultThreadCount = 64
	dialerTimeout      = 15
)

func scan(cmd *cobra.Command, args []string) error {
	var target string

	//target = flag.String("t", "", "target for scanning")
	switch len(args) {
	case 0:
		return ErrEmptyTarget
	case 1:
	default:
		target = args[0]
	}

	threads := defaultThreadCount
	//threads := flag.Int("s", defaultThreadCount, "number of threads(streams) when scanning")

	scanners.NewScanner(
		producers.NewProducer(),
		workers.NewWorker(&net.Dialer{Timeout: dialerTimeout * time.Second}),
		reporters.NewReporter(adapters.PrinterAdapter{}),
	)(target, threads)

	return nil
}
