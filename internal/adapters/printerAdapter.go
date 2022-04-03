package adapters

import "fmt"

type PrinterAdapter struct{}

func (p PrinterAdapter) Print(args ...interface{}) {
	fmt.Print(args...)
}

func (p PrinterAdapter) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
