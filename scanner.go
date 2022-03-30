package main

import "sync"

func newScanner(
	producer func(from, to int, ports chan<- int),
	worker func(target string, ports <-chan int, opened chan<- int),
	reporter func(target string, opened <-chan int),
) func(address string, threads int) {
	return func(target string, threads int) {
		waitScanner := sync.WaitGroup{}
		ports := make(chan int, threads)
		opened := make(chan int, threads)

		waitScanner.Add(1)
		go func() {
			defer func() {
				close(ports)
				waitScanner.Done()
			}()
			producer(minPortNumber, maxPortNumber, ports)
		}()

		waitScanner.Add(1)
		go func() {
			defer func() {
				close(opened)
				waitScanner.Done()
			}()
			waitWorkers := sync.WaitGroup{}
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
