package scanners_test

import (
	"fmt"
	"github.com/psssix/syn-scanner/internal/scanners"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/psssix/syn-scanner/pkg/producers"
)

type processCounter struct {
	m sync.Mutex
	c int
}

func (pc *processCounter) add() {
	pc.m.Lock()
	pc.c++
	pc.m.Unlock()
}

func (pc *processCounter) value() int {
	return pc.c
}

func TestScannerIntegrityWork(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		target  string
		threads int
	}{
		{"", "test.local", 8},
		{"", "127.0.0.1", 16},
		{"", "127.0.0.3", 32},
		{"", "test2.local", 100000},
	}

	for _, test := range tests {
		test := test
		test.name = fmt.Sprintf("scan %s with %d threads", test.target, test.threads)
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			producerRunCount, workerRunCount, reporterRunCount := processCounter{}, processCounter{}, processCounter{}
			targetPort := rand.Int()
			openedPort := rand.Int()

			scanners.NewScanner(
				func(from, to int, ports chan<- int) {
					producerRunCount.add()
					assert.Equal(t, producers.MinPortNumber, from)
					assert.Equal(t, producers.MaxPortNumber, to)
					for i := 0; i < test.threads; i++ {
						ports <- targetPort
					}
				},
				func(target string, ports <-chan int, opened chan<- int) {
					workerRunCount.add()
					assert.Equal(t, test.target, target)
					assert.Equal(t, targetPort, <-ports)
					opened <- openedPort
				},
				func(target string, opened <-chan int) {
					reporterRunCount.add()
					for port := range opened {
						assert.Equal(t, openedPort, port)
					}
				},
			)(test.target, test.threads)

			assert.Equal(t, 1, producerRunCount.value())
			assert.Equal(t, test.threads, workerRunCount.value())
			assert.Equal(t, 1, reporterRunCount.value())
		})
	}
}
