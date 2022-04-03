package main

import (
	"flag"
	"net"
	"os"
	"syn-scanner/pkg/producers"
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
		newWorker(&net.Dialer{Timeout: dialerTimeout * time.Second}),
		newReporter(printerAdapter{}),
	)(*target, *threads)
}
