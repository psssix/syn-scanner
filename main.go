package main

import (
	"flag"
	"net"
	"os"
	"syn-scanner/pkg/producers"
	"syn-scanner/pkg/reporters"
	"syn-scanner/pkg/workers"
	"time"
)

const defaultThreadCount = 8
const dialerTimeout = 3

func main() {
	target := flag.String("t", "", "target for scanning")
	threads := flag.Int("s", defaultThreadCount, "number of threads(streams) when scanning")
	flag.Parse()
	if *target == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	newScanner(
		producers.NewProducer(),
		workers.NewWorker(&net.Dialer{Timeout: dialerTimeout * time.Second}),
		reporters.NewReporter(printerAdapter{}),
	)(*target, *threads)
}
