package scanners_test

import (
	"sync/atomic"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/psssix/syn-scanner/internal/scanners"
	"github.com/stretchr/testify/assert"
)

func TestScannerIntegrityWork(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, target string
		threads      int
	}{
		{name: "scan test.local with 8 threads", target: "test.local", threads: 8},
		{name: "scan 127.0.0.1 with 16 threads", target: "127.0.0.1", threads: 16},
		{name: "scan 127.0.0.3 with 32 threads", target: "127.0.0.3", threads: 32},
		{name: "scan test2.local with 100000 threads", target: "test2.local", threads: 100000},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var (
				producerRunCount, workerRunCount, reporterRunCount atomic.Uint64

				f          = faker.New()
				targetPort = f.Int()
				openedPort = f.Int()
			)

			scanners.NewScanner(
				func(ports chan<- int) {
					producerRunCount.Add(1)

					assert.Equal(t, test.threads, cap(ports))

					for i := 0; i < test.threads; i++ {
						ports <- targetPort
					}
				},
				func(target string, ports <-chan int, opened chan<- int) {
					workerRunCount.Add(1)

					assert.Equal(t, test.target, target)
					assert.Equal(t, test.threads, cap(ports))
					assert.Equal(t, test.threads, cap(opened))

					assert.Equal(t, targetPort, <-ports)

					opened <- openedPort
				},
				func(target string, opened <-chan int) {
					reporterRunCount.Add(1)

					assert.Equal(t, test.target, target)
					assert.Equal(t, test.threads, cap(opened))

					for port := range opened {
						assert.Equal(t, openedPort, port)
					}
				},
			)(test.target, test.threads)

			assert.EqualValues(t, 1, producerRunCount.Load())
			assert.EqualValues(t, test.threads, workerRunCount.Load())
			assert.EqualValues(t, 1, reporterRunCount.Load())
		})
	}
}
