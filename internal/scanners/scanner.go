package scanners

import (
	"sync"

	"github.com/psssix/syn-scanner/pkg/producers"
	"github.com/psssix/syn-scanner/pkg/reporters"
)

type (
	Worker  func(target string, ports <-chan int, opened chan<- int)
	Scanner func(target string, threads int)
)

func NewScanner(producer producers.Producer, worker Worker, reporter reporters.Reporter) Scanner {
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
			producer(ports)
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
