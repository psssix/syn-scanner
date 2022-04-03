package mocks

import (
	"github.com/stretchr/testify/mock"
	"net"
	"time"
)

type Connection struct {
	mock.Mock
}

func (c *Connection) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (c *Connection) Write(b []byte) (n int, err error) {
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

func (c *Connection) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (c *Connection) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}
