package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProducer(t *testing.T) {
	tests := []struct {
		from   int
		to     int
		result []int
	}{
		{1, 1, []int{1}},
		{1, 2, []int{1, 2}},
		{1, 10, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{65526, 65535, []int{65526, 65527, 65528, 65529, 65530, 65531, 65532, 65533, 65534, 65535}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("create producer for %d, %d range", test.from, test.to), func(t *testing.T) {
			var (
				chanLength = test.to - test.from + 1
				ports      = make(chan int, chanLength)
				result     []int
			)

			newProducer(test.from, test.to)(ports)

			assert.Equal(t, chanLength, cap(ports))
			assert.Equal(t, chanLength, len(ports))

			for port := range ports {
				result = append(result, port)
			}
			assert.Equal(t, test.result, result)
		})
	}
}

func TestProducerPanicWhenUsingInvalidPorts(t *testing.T) {
	tests := []struct {
		from int
		to   int
	}{
		{0, 1},
		{1, 0},
		{65535, 65536},
		{65536, 65535},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("create producer for %d, %d range", test.from, test.to), func(t *testing.T) {
			var ports = make(chan int, 2)

			assert.Panicsf(t, func() {
				newProducer(test.from, test.to)(ports)
			}, "invalid ports range, ports can be in range from 1 to 65535")
		})
	}
}

// to must be greater than from
func TestProducerPanicWhenUsingInvalidRange(t *testing.T) {
	tests := []struct {
		from int
		to   int
	}{
		{10, 1},
		{65535, 65536},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("create producer for %d, %d range", test.from, test.to), func(t *testing.T) {
			var ports = make(chan int, 10)

			assert.Panicsf(t, func() {
				newProducer(test.from, test.to)(ports)
			}, "to must be greater than from")
		})
	}
}
