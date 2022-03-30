package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
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
	tests := []struct {
		target  string
		threads int
	}{
		{"test.local", 8},
		{"127.0.0.1", 16},
		{"127.0.0.3", 32},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("scan %s with %d threads", test.target, test.threads), func(t *testing.T) {
			producerRunCount, workerRunCount, reporterRunCount := processCounter{}, processCounter{}, processCounter{}
			targetPort := rand.Int()
			openedPort := rand.Int()

			newScanner(
				func(from, to int, ports chan<- int) {
					producerRunCount.add()
					assert.Equal(t, minPortNumber, from)
					assert.Equal(t, maxPortNumber, to)
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
