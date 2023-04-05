package workers

import (
	"fmt"
	"net"
)

type dialer interface {
	Dial(network, address string) (net.Conn, error)
}

func NewSYNScan(d dialer) func(target string, ports <-chan int, opened chan<- int) {
	return func(target string, ports <-chan int, opened chan<- int) {
		for port := range ports {
			connection, err := d.Dial("tcp", fmt.Sprintf("%s:%d", target, port))
			if err != nil {
				continue
			}

			if connection.Close() != nil {
				panic(fmt.Sprintf("can't close opened connection for %s", fmt.Sprintf("%s:%d", target, port)))
			}

			opened <- port
		}
	}
}
