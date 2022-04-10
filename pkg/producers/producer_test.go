package producers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/psssix/syn-scanner/pkg/producers"
)

func TestProducerGeneratesRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		from     int
		to       int
		expected []int
	}{
		{"", 1, 1, []int{1}},
		{"", 1, 2, []int{1, 2}},
		{"", 1, 10, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{"", 65526, 65535, []int{65526, 65527, 65528, 65529, 65530, 65531, 65532, 65533, 65534, 65535}},
	}

	for _, test := range tests {
		test := test
		test.name = fmt.Sprintf("produce for %d-%d range", test.from, test.to)
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ports := make(chan int, test.to-test.from+1)

			producers.NewProducer()(test.from, test.to, ports)
			close(ports)

			var actual []int
			for port := range ports {
				actual = append(actual, port)
			}
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestProducerPanicsWhenUsingInvalidPorts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		from int
		to   int
	}{
		{"produce with 'from' less that 1", 0, 1},
		{"produce with 'from' greater that 65535", 65536, 65535},
		{"produce with 'to' less that 1", 1, 0},
		{"produce with 'to' greater that 65535", 65535, 65536},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ports := make(chan int)
			close(ports)
			assert.PanicsWithValue(t, "invalid ports range, ports can be in range from 1 to 65535",
				func() {
					producers.NewProducer()(test.from, test.to, ports)
				})
		})
	}
}

func TestProducerPanicsWhenUsingInvalidRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		from int
		to   int
	}{
		{"", 10, 1},
		{"", 65535, 65526},
	}

	for _, test := range tests {
		test := test
		test.name = fmt.Sprintf("produce for %d-%d range", test.from, test.to)
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ports := make(chan int)
			close(ports)
			assert.PanicsWithValue(t, "'to' must be greater than 'from'",
				func() {
					producers.NewProducer()(test.from, test.to, ports)
				})
		})
	}
}
