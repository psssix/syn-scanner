package reporters

type printer interface { //nolint:forbidigo // linter false-positive
	// Print formats using the default formats for its operands and writes to standard output.
	// Spaces are added between operands when neither is a string.
	// It returns the number of bytes written and any write error encountered.
	Print(args ...interface{})

	// Printf formats according to a format specifier and writes to standard output.
	// It returns the number of bytes written and any write error encountered.
	Printf(format string, args ...interface{})
}
