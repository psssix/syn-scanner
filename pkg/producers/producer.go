package producers

import "fmt"

const (
	MinPortNumber = 1
	MaxPortNumber = 65535
)

func NewProducer() func(from, to int, ports chan<- int) {
	return func(from, to int, ports chan<- int) {
		if from < MinPortNumber || from > MaxPortNumber || to < MinPortNumber || to > MaxPortNumber {
			panic(fmt.Sprintf("invalid ports range, ports can be in range from %d to %d",
				MinPortNumber, MaxPortNumber))
		}
		if from > to {
			panic("'to' must be greater than 'from'")
		}

		for i := from; i <= to; i++ {
			ports <- i
		}
	}
}
