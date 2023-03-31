package producers_test

import (
	"testing"

	"github.com/psssix/syn-scanner/pkg/producers"
	"github.com/stretchr/testify/assert"
)

func TestProducerGeneratesRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		from, to int
		want     []int
	}{
		{name: "produce for 1-1 range", from: 1, to: 1, want: []int{1}},
		{name: "produce for 1-2 range", from: 1, to: 2, want: []int{1, 2}},
		{name: "produce for 1-10 range", from: 1, to: 10, want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{
			name: "produce for 65526-65535 range",
			from: 65526, to: 65535,
			want: []int{65526, 65527, 65528, 65529, 65530, 65531, 65532, 65533, 65534, 65535},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var (
				ports  = make(chan int, test.to-test.from+1)
				actual = make([]int, 0, test.to-test.from+1)
			)

			producers.NewProducer(test.from, test.to)(ports)
			close(ports)

			for port := range ports {
				actual = append(actual, port)
			}

			assert.Equal(t, test.want, actual)
		})
	}
}

func TestNewProducerWithPanicsWhenUsingInvalidPorts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		from, to int
	}{
		{name: "produce with 'from' less that 1", from: 0, to: 1},
		{name: "produce with 'from' greater that 65535", from: 65536, to: 65535},
		{name: "produce with 'to' less that 1", from: 1, to: 0},
		{name: "produce with 'to' greater that 65535", from: 65535, to: 65536},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.PanicsWithValue(t, "invalid ports range, ports can be in range from 1 to 65535",
				func() {
					producers.NewProducer(test.from, test.to)
				},
			)
		})
	}
}

func TestNewProducerWithPanicsWhenUsingInvalidRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		from, to int
	}{
		{name: "produce for 10-1 range", from: 10, to: 1},
		{name: "produce for 65535-65526 range", from: 65535, to: 65526},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.PanicsWithValue(t, "'to' must be greater than 'from'",
				func() {
					producers.NewProducer(test.from, test.to)
				},
			)
		})
	}
}
