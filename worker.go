package main

import "fmt"

func newWorker(target string, d dialer) func(ports <-chan int, opened chan<- int) {
	return func(ports <-chan int, opened chan<- int) {
		defer close(opened)
		for port := range ports {
			address := fmt.Sprintf("%s:%d", target, port)
			connection, err := d.Dial("tcp", address)
			if err != nil {
				continue
			}
			if connection.Close() != nil {
				panic(fmt.Sprintf("can't close opened connection for %s", address))
			}
			opened <- port
		}
	}
}
