package producers

import (
	"fmt"

	"github.com/psssix/syn-scanner/internal/scanners"
)

const (
	MinPortNumber = 1
	MaxPortNumber = 65535
)

func NewProducer(from, to int) scanners.Producer {
	if from < MinPortNumber || from > MaxPortNumber || to < MinPortNumber || to > MaxPortNumber {
		panic(fmt.Sprintf(
			"invalid ports range, ports can be in range from %d to %d", MinPortNumber, MaxPortNumber,
		))
	}

	if from > to {
		panic("'to' must be greater than 'from'")
	}

	return func(ports chan<- int) {
		for i := from; i <= to; i++ {
			ports <- i
		}
	}
}
