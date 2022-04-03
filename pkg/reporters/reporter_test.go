package reporters_test

import (
	"fmt"
	"syn-scanner/pkg/mocks"
	"syn-scanner/pkg/reporters"
	"testing"
)

func TestReporterPrints(t *testing.T) {
	tests := []struct {
		target string
		ports  []int
	}{
		{"test.local", []int{10}},
		{"127.0.0.1", []int{20, 30}},
		{"test2.local", []int{30, 31, 32, 33}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("inform about target %s and ports %v", test.target, test.ports), func(t *testing.T) {
			opened := make(chan int, len(test.ports))
			printer := new(mocks.Printer)
			printer.On("Printf", "scanning: %s opened ports: ", []interface{}{test.target})
			printer.On("Print", []interface{}{"\ndone\n"})
			for at, port := range test.ports {
				opened <- port
				if at == 0 {
					printer.On("Printf", "%d", []interface{}{port})
				} else {
					printer.On("Printf", ", %d", []interface{}{port})
				}
			}
			close(opened)

			reporters.NewReporter(printer)(test.target, opened)

			printer.AssertExpectations(t)
		})
	}
}

func TestReporterPrintsWhenNoPortOpen(t *testing.T) {
	const target = "test.local"

	printer := new(mocks.Printer)
	printer.On("Printf", "scanning: %s opened ports: ", []interface{}{target})
	printer.On("Print", []interface{}{"none"})
	printer.On("Print", []interface{}{"\ndone\n"})
	opened := make(chan int)
	close(opened)

	reporters.NewReporter(printer)(target, opened)

	printer.AssertExpectations(t)
}