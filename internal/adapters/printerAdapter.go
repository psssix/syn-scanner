package adapters

import "fmt"

type PrinterAdapter struct{}

func (p PrinterAdapter) Print(args ...interface{}) {
	fmt.Print(args...) //nolint:forbidigo // this structure is a wrapper, ignore the use of the fmt
}

func (p PrinterAdapter) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...) //nolint:forbidigo // this structure is a wrapper, ignore the use of the fmt
}
