package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProducerGenerateRange(t *testing.T) {
	tests := []struct {
		from     int
		to       int
		expected []int
	}{
		{1, 1, []int{1}},
		{1, 2, []int{1, 2}},
		{1, 10, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{65526, 65535, []int{65526, 65527, 65528, 65529, 65530, 65531, 65532, 65533, 65534, 65535}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("producer for %d-%d range", test.from, test.to), func(t *testing.T) {
			ports := make(chan int, test.to-test.from+1)

			newProducer(test.from, test.to)(ports)

			var actual []int
			for port := range ports {
				actual = append(actual, port)
			}
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestProducerPanicWhenUsingInvalidPorts(t *testing.T) {
	tests := []struct {
		name string
		from int
		to   int
	}{
		{"producer with from less that 1", 0, 1},
		{"producer with from greater that 65535", 65536, 65535},
		{"producer with to less that 1", 1, 0},
		{"producer with to greater that 65535", 65535, 65536},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ports := make(chan int)
			assert.PanicsWithValue(t, "invalid ports range, ports can be in range from 1 to 65535",
				func() {
					newProducer(test.from, test.to)(ports)
				})
		})
	}
}

func TestProducerPanicWhenUsingInvalidRange(t *testing.T) {
	tests := []struct {
		from int
		to   int
	}{
		{10, 1},
		{65535, 65526},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("producer for %d-%d range", test.from, test.to), func(t *testing.T) {
			ports := make(chan int)
			assert.PanicsWithValue(t, "'to' must be greater than 'from'",
				func() {
					newProducer(test.from, test.to)(ports)
				})
		})
	}
}
