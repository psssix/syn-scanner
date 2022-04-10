package mocks

import (
	"fmt"
	"net"

	"github.com/stretchr/testify/mock"
)

type Dialer struct {
	mock.Mock
}

func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	args := d.Called(network, address)
	c := args.Get(0)
	if c == nil {
		return nil, args.Error(1)
	}
	conn, ok := c.(net.Conn)
	if !ok {
		panic(fmt.Sprintf("interface conversion: interface {} is %T, not net.Conn", c))
	}
	return conn, nil
}
