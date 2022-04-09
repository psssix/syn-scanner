package workers_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/psssix/syn-scanner/pkg/mocks"
	"github.com/psssix/syn-scanner/pkg/workers"
	"github.com/stretchr/testify/assert"
)

func TestWorkerDials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		target string
		ports  []int
	}{
		{"", "test.local", []int{10}},
		{"", "127.0.0.1", []int{20, 30}},
		{"", "test2.local", []int{30, 31, 32, 33}},
	}

	for _, test := range tests {
		test := test
		test.name = fmt.Sprintf("work with target %s and ports %v", test.target, test.ports)
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ports := make(chan int, len(test.ports))
			opened := make(chan int, len(test.ports))
			dialer := new(mocks.Dialer)
			var connections []*mocks.Connection
			for _, port := range test.ports {
				ports <- port
				address := fmt.Sprintf("%s:%d", test.target, port)

				c := new(mocks.Connection)
				c.On("Close").Return(nil)
				dialer.On("Dial", "tcp", address).Return(c, nil)
				connections = append(connections, c)
			}
			close(ports)

			workers.NewWorker(dialer)(test.target, ports, opened)
			close(opened)

			actualOpened := make([]int, 0)
			for port := range opened {
				actualOpened = append(actualOpened, port)
			}
			assert.Equal(t, test.ports, actualOpened)
			dialer.AssertExpectations(t)
			for _, c := range connections {
				c.AssertExpectations(t)
			}
		})
	}
}

func TestWorkerDialsAndSomeConnectionIsNotOpen(t *testing.T) {
	t.Parallel()

	const target = "test.local"
	targetPorts := []struct {
		number     int
		canConnect bool
	}{
		{30, true},
		{31, false},
		{32, true},
	}

	ports := make(chan int, len(targetPorts))
	opened := make(chan int, len(targetPorts))
	var expectedOpened []int
	dialer := new(mocks.Dialer)
	var connections []*mocks.Connection
	for _, port := range targetPorts {
		ports <- port.number
		address := fmt.Sprintf("%s:%d", target, port.number)

		if port.canConnect {
			expectedOpened = append(expectedOpened, port.number)
			c := new(mocks.Connection)
			c.On("Close").Return(nil)
			dialer.On("Dial", "tcp", address).Return(c, nil)
			connections = append(connections, c)
		} else {
			err := &net.OpError{Op: "dial", Net: "tcp", Source: nil, Addr: nil, Err: nil}
			dialer.On("Dial", "tcp", address).Return(nil, err)
		}
	}
	close(ports)

	workers.NewWorker(dialer)(target, ports, opened)
	close(opened)

	actualOpened := make([]int, 0)
	for port := range opened {
		actualOpened = append(actualOpened, port)
	}
	assert.Equal(t, expectedOpened, actualOpened)
	dialer.AssertExpectations(t)
	for _, c := range connections {
		c.AssertExpectations(t)
	}
}

func TestWorkerPanicsWhenConnectionIsNotClose(t *testing.T) {
	t.Parallel()

	const target = "test.local"
	const port = 80

	ports := make(chan int, 1)
	opened := make(chan int)
	close(opened)

	ports <- port
	close(ports)

	address := fmt.Sprintf("%s:%d", target, port)
	err := &net.OpError{Op: "close", Net: "tcp", Source: nil, Addr: nil, Err: nil}
	c := new(mocks.Connection)
	c.On("Close").Return(err)
	dialer := new(mocks.Dialer)
	dialer.On("Dial", "tcp", address).Return(c, nil)

	assert.PanicsWithValue(t, "can't close opened connection for test.local:80",
		func() {
			workers.NewWorker(dialer)(target, ports, opened)
		})
	dialer.AssertExpectations(t)
	c.AssertExpectations(t)
}
