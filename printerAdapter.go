package main

import "fmt"

type printerAdapter struct{}

func (p printerAdapter) Print(args ...interface{}) {
	fmt.Print(args...)
}

func (p printerAdapter) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
