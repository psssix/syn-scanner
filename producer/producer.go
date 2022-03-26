package producer

import "fmt"

const (
	minPortNumber = 1
	maxPortNumber = 65535
)

func NewProducer(from, to int) func(ports chan<- int) {
	return func(ports chan<- int) {
		defer close(ports)
		if from < minPortNumber || from > maxPortNumber || to < minPortNumber || to > maxPortNumber {
			panic(fmt.Sprintf("invalid ports range, ports can be in range from %d to %d",
				minPortNumber, maxPortNumber))
		}
		if from > to {
			panic("to must be greater than from")
		}

		for i := from; i <= to; i++ {
			ports <- i
		}
	}
}
