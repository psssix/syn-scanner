package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
)

//
//func TestScannerWork(t *testing.T) {
//	tests := []struct {
//		target  string
//		threads int
//	}{
//		{"test.local", 8},
//		{"127.0.0.1", 16},
//		{"127.0.0.3", 32},
//	}
//
//	for _, test := range tests {
//		t.Run(fmt.Sprintf("scan %s with %d threads", test.target, test.threads), func(t *testing.T) {
//			isProducerRun := false
//			workerRunCount := 0
//			isReporterRun := false
//			var (
//				actualFrom, actualTo int
//				mutex                sync.Mutex
//				actualTargets        []string
//				actualPorts          []int
//				actualOpened         []int
//			)
//			testPorts := make([]int, test.threads)
//			for i := range testPorts {
//				testPorts[i] = rand.Int()
//			}
//			testOpened := make([]int, test.threads)
//			for i := range testOpened {
//				testOpened[i] = rand.Int()
//			}
//
//			wg := new(mocks.WaitGroup)
//			wg.On("Add", 1).Times(test.threads + 2)
//			wg.On("Done").Times(test.threads + 2)
//			wg.On("Wait").Once()
//
//			newScanner(
//				func(from, to int, ports chan<- int) {
//					defer close(ports)
//					isProducerRun = true
//					actualFrom = from
//					actualTo = to
//					for _, port := range testPorts {
//						ports <- port
//					}
//				},
//				func(target string, ports <-chan int, opened chan<- int) {
//					defer close(opened)
//					mutex.Lock()
//					workerRunCount++
//					actualTargets = append(actualTargets, target)
//					actualPorts = append(actualPorts, <-ports)
//					mutex.Unlock()
//					for _, port := range testOpened {
//						opened <- port
//					}
//				},
//				func(target string, opened <-chan int) {
//					isReporterRun = true
//					for port := range opened {
//						actualOpened = append(actualOpened, port)
//					}
//				},
//				wg,
//			)(test.target, test.threads)
//
//			assert.True(t, isProducerRun)
//			assert.Equal(t, minPortNumber, actualFrom)
//			assert.Equal(t, maxPortNumber, actualTo)
//			assert.Equal(t, test.threads, workerRunCount)
//			var expectedTargets []string
//			for i := 0; i < test.threads; i++ {
//				expectedTargets = append(expectedTargets, test.target)
//			}
//			assert.Equal(t, expectedTargets, actualTargets)
//			assert.Equal(t, testPorts, actualPorts)
//			assert.True(t, isReporterRun)
//			assert.Equal(t, testOpened, actualOpened)
//			wg.AssertExpectations(t)
//		})
//	}
//}

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
