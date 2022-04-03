package workers

import (
	"fmt"
	"net"
)

type dialer interface {
	Dial(network, address string) (net.Conn, error)
}

func NewWorker(d dialer) func(target string, ports <-chan int, opened chan<- int) {
	return func(target string, ports <-chan int, opened chan<- int) {
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
