package semaphore_worker_pool

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestRunNonSemaphorePool(t *testing.T) {
	work := getCommonWork()
	maxWorkers := 20
	start := time.Now()

	// this will be faster because the workers just "go"
	RunNonSemaphorePool(maxWorkers, work)
	fmt.Printf("Non semaphore workers took: %v\n", time.Since(start).Seconds())
}

func TestRunSemaphorePool(t *testing.T) {
	work := getCommonWork()
	start := time.Now()
	maxWorkers, maxConcurrency := 20, 4

	// this will be slower because of the "barrier" but has better flexibility
	// on goroutine counts. You'll see work being performed in batches
	RunSemaphorePool(maxWorkers, maxConcurrency, work)
	fmt.Printf("Semaphore restricted workers took: %v\n", time.Since(start).Seconds())
}

// this returns some work for the workers to perform. To keep the work same
// for semaphore and non semaphore pools, we sleep for 1 second each
func getCommonWork() func(id int, ch chan<- string) {
	return func(id int, ch chan<- string) {
		time.Sleep(time.Second)
		ch <- fmt.Sprintf("Done #%d", id)
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
