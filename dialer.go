package main

import (
	"net"
)

type dialer interface {
	Dial(network, address string) (net.Conn, error)
}
