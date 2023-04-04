package workers_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/psssix/syn-scanner/pkg/workers"
	"github.com/stretchr/testify/assert"
)

func TestSynScanDials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, target string
		ports        []int
	}{
		{name: `work with target test.local and ports 10`, target: "test.local", ports: []int{10}},
		{name: `work with target 127.0.0.1 and ports 20 and 30`, target: "127.0.0.1", ports: []int{20, 30}},
		{
			name:   `work with target test2.local and ports 30, 31, 32 and 33`,
			target: "test2.local",
			ports:  []int{30, 31, 32, 33},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var (
				ports  = make(chan int, len(test.ports))
				opened = make(chan int, len(test.ports))

				connections = make([]*workers.Connection, 0, len(test.ports))
				dialer      = new(workers.Dialer)
			)

			for _, port := range test.ports {
				ports <- port

				conn := new(workers.Connection)
				conn.On("Close").Return(nil)
				connections = append(connections, conn)

				address := fmt.Sprintf("%s:%d", test.target, port)
				dialer.On("Dial", "tcp", address).Return(conn, nil)
			}
			close(ports)

			workers.NewSynScan(dialer)(test.target, ports, opened)
			close(opened)

			assertOpenedPorts(t, test.ports, opened)

			for _, c := range connections {
				c.AssertExpectations(t)
			}
			dialer.AssertExpectations(t)
		})
	}
}

func assertOpenedPorts(t *testing.T, expected []int, actual chan int) {
	t.Helper()

	actualPorts := make([]int, 0, len(expected))
	for port := range actual {
		actualPorts = append(actualPorts, port)
	}

	assert.Equal(t, expected, actualPorts)
}

func TestSynScanDialsAndSomeConnectionIsNotOpen(t *testing.T) {
	t.Parallel()

	var (
		target      = "test.local"
		targetPorts = []struct {
			number     int
			canConnect bool
		}{
			{number: 30, canConnect: true},
			{number: 31, canConnect: false},
			{number: 32, canConnect: true},
		}

		ports          = make(chan int, len(targetPorts))
		opened         = make(chan int, len(targetPorts))
		expectedOpened = make([]int, 0, len(targetPorts))

		connections = make([]*workers.Connection, 0, len(targetPorts))
		dialer      = new(workers.Dialer)

		dialerErr = &net.OpError{Op: "dial", Net: "tcp", Source: nil, Addr: nil, Err: nil}
	)

	for _, port := range targetPorts {
		ports <- port.number
		address := fmt.Sprintf("%s:%d", target, port.number)

		if port.canConnect {
			expectedOpened = append(expectedOpened, port.number)

			conn := new(workers.Connection)
			conn.On("Close").Return(nil)
			connections = append(connections, conn)

			dialer.On("Dial", "tcp", address).Return(conn, nil)
		} else {
			dialer.On("Dial", "tcp", address).Return(nil, dialerErr)
		}
	}

	close(ports)

	workers.NewSynScan(dialer)(target, ports, opened)
	close(opened)

	assertOpenedPorts(t, expectedOpened, opened)

	for _, c := range connections {
		c.AssertExpectations(t)
	}
	dialer.AssertExpectations(t)

}

func TestSynScanPanicsWhenConnectionIsNotClose(t *testing.T) {
	t.Parallel()

	var (
		target = "test.local"
		port   = 80

		ports  = make(chan int, 1)
		opened = make(chan int)

		address = fmt.Sprintf("%s:%d", target, port)

		conn   = new(workers.Connection)
		dialer = new(workers.Dialer)

		connectionErr = &net.OpError{Op: "close", Net: "tcp", Source: nil, Addr: nil, Err: nil}
	)

	close(opened)

	ports <- port
	close(ports)

	conn.On("Close").Return(connectionErr)
	dialer.On("Dial", "tcp", address).Return(conn, nil)

	assert.PanicsWithValue(t, "can't close opened connection for test.local:80", func() {
		workers.NewSynScan(dialer)(target, ports, opened)
	})

	conn.AssertExpectations(t)
	dialer.AssertExpectations(t)
}
