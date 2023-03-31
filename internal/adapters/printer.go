package adapters

import "fmt"

type Printer struct{}

// Print formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (p Printer) Print(args ...interface{}) {
	fmt.Print(args...) //nolint:forbidigo // this structure is a wrapper, ignore the use of the fmt
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func (p Printer) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...) //nolint:forbidigo // this structure is a wrapper, ignore the use of the fmt
}
