package reporters

type printer interface {
	// Print appends args to the message output.
	Print(args ...interface{})

	// Printf writes a formatted string.
	Printf(format string, args ...interface{})
}

func NewReporter(p printer) func(target string, opened <-chan int) {
	return func(target string, opened <-chan int) {
		p.Printf("scanning: %s opened ports: ", target)
		firstPrint := true
		for port := range opened {
			if firstPrint {
				p.Printf("%d", port)
				firstPrint = false
			} else {
				p.Printf(", %d", port)
			}
		}
		if firstPrint {
			p.Print("none")
		}
		p.Print("\ndone\n")
	}
}
