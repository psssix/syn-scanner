package mocks

import (
	"github.com/stretchr/testify/mock"
	"net"
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
	return c.(net.Conn), nil
}
