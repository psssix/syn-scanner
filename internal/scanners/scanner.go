package scanners

import (
	"sync"

	"github.com/psssix/syn-scanner/pkg/producers"
)

type (
	Producer func(from, to int, ports chan<- int)
	Worker   func(target string, ports <-chan int, opened chan<- int)
	Reporter func(target string, opened <-chan int)
	Scanner  func(target string, threads int)
)

func NewScanner(producer Producer, worker Worker, reporter Reporter) Scanner {
	return func(target string, threads int) {
		var (
			waitScanner = sync.WaitGroup{}
			waitWorkers = sync.WaitGroup{}

			ports  = make(chan int, threads)
			opened = make(chan int, threads)
		)

		waitScanner.Add(1)

		go func() {
			defer func() {
				close(ports)
				waitScanner.Done()
			}()
			producer(producers.MinPortNumber, producers.MaxPortNumber, ports)
		}()

		waitScanner.Add(1)

		go func() {
			defer func() {
				close(opened)
				waitScanner.Done()
			}()

			for i := 0; i < threads; i++ {
				waitWorkers.Add(1)

				go func() {
					defer waitWorkers.Done()
					worker(target, ports, opened)
				}()
			}

			waitWorkers.Wait()
		}()

		waitScanner.Add(1)

		go func() {
			defer waitScanner.Done()
			reporter(target, opened)
		}()

		waitScanner.Wait()
	}
}
