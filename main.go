package main

import (
	"flag"
	"net"
	"os"
	"time"

	"github.com/psssix/syn-scanner/internal/adapters"
	"github.com/psssix/syn-scanner/pkg/producers"
	"github.com/psssix/syn-scanner/pkg/reporters"
	"github.com/psssix/syn-scanner/pkg/scanners"
	"github.com/psssix/syn-scanner/pkg/workers"
)

const (
	defaultThreadCount = 8
	dialerTimeout      = 3
)

func main() {
	target := flag.String("t", "", "target for scanning")
	threads := flag.Int("s", defaultThreadCount, "number of threads(streams) when scanning")
	flag.Parse()
	if *target == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	scanners.NewScanner(
		producers.NewProducer(),
		workers.NewWorker(&net.Dialer{Timeout: dialerTimeout * time.Second}),
		reporters.NewReporter(adapters.PrinterAdapter{}),
	)(*target, *threads)
}
