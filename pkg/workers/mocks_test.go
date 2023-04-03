package workers

import (
	"net"

	"github.com/stretchr/testify/mock"
)

type (
	Dialer struct {
		mock.Mock
	}
	Connection struct {
		net.Conn
		mock.Mock
	}
)

func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	args := d.Called(network, address)

	conn := args.Get(0)
	if conn == nil {
		return nil, args.Error(1)
	}

	return conn.(net.Conn), nil
}

func (c *Connection) Close() error {
	args := c.Called()
	return args.Error(0)
}
