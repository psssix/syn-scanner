package mocks

import (
	"net"
	"time"

	"github.com/stretchr/testify/mock"
)

type Connection struct {
	mock.Mock
}

func (c *Connection) Read([]byte) (int, error) {
	panic("implement me")
}

func (c *Connection) Write([]byte) (int, error) {
	panic("implement me")
}

func (c *Connection) Close() error {
	args := c.Called()
	return args.Error(0)
}

func (c *Connection) LocalAddr() net.Addr {
	panic("implement me")
}

func (c *Connection) RemoteAddr() net.Addr {
	panic("implement me")
}

func (c *Connection) SetDeadline(time.Time) error {
	panic("implement me")
}

func (c *Connection) SetReadDeadline(time.Time) error {
	panic("implement me")
}

func (c *Connection) SetWriteDeadline(time.Time) error {
	panic("implement me")
}
