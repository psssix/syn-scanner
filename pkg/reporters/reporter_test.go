package reporters

import (
	"testing"
)

func TestReporterPrints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, target string
		ports        []int
	}{
		{name: `inform about target "test.local" and ports 10`, target: "test.local", ports: []int{10}},
		{name: `inform about target "127.0.0.1" and ports 20 and 30`, target: "127.0.0.1", ports: []int{20, 30}},
		{
			name:   `inform about target "test2.local" and ports 30, 31, 32 and 33`,
			target: "test2.local",
			ports:  []int{30, 31, 32, 33},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var (
				opened  = make(chan int, len(test.ports))
				printer = new(printerMock) //nolint:forbidigo // linter false-positive
			)

			printer.On("Printf", "scanning %q opened ports is: ", []interface{}{test.target}) //nolint:forbidigo,lll // linter false-positive
			printer.On("Print", []interface{}{"\ndone\n"})                                    //nolint:forbidigo,lll // linter false-positive

			for i, port := range test.ports {
				opened <- port
				if i == 0 {
					printer.On("Printf", "%d", []interface{}{port}) //nolint:forbidigo // linter false-positive
				} else {
					printer.On("Printf", ", %d", []interface{}{port}) //nolint:forbidigo // linter false-positive
				}
			}
			close(opened)

			NewReporter(printer)(test.target, opened) //nolint:forbidigo // linter false-positive

			printer.AssertExpectations(t) //nolint:forbidigo // linter false-positive
		})
	}
}

func TestReporterPrintsWhenNoPortOpen(t *testing.T) {
	t.Parallel()

	var (
		target  = "test.local"
		opened  = make(chan int)
		printer = new(printerMock) //nolint:forbidigo // linter false-positive
	)

	printer.On("Printf", "scanning %q opened ports is: ", []interface{}{target}) //nolint:forbidigo,lll // linter false-positive
	printer.On("Print", []interface{}{"none"})                                   //nolint:forbidigo,lll // linter false-positive
	printer.On("Print", []interface{}{"\ndone\n"})                               //nolint:forbidigo,lll // linter false-positive

	close(opened)

	NewReporter(printer)(target, opened) //nolint:forbidigo // linter false-positive

	printer.AssertExpectations(t) //nolint:forbidigo // linter false-positive
}
