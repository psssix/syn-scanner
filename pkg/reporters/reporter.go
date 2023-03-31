package reporters

type Reporter func(target string, opened <-chan int)

func NewReporter(p printer) Reporter { //nolint:forbidigo // linter false-positive
	return func(target string, opened <-chan int) {
		p.Printf("scanning %q opened ports is: ", target)

		firstPrint := true
		for port := range opened {
			if firstPrint {
				p.Printf("%d", port)
				firstPrint = false //nolint:wsl // in this case it doesn't make sense
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
