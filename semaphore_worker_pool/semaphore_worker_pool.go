package semaphore_worker_pool

import (
	"fmt"
)

// Zen: When fanning out to a high number of goroutines that write on an unbuffered channel
// (to improve latency) make sure that the count is not too high. For such high fanouts,
// ensure the use of a semaphore that doesn't allow the creation of a small set of routines.

// RunNonSemaphorePool This pool of workers fan out and create as many goroutines as there are workers
// Each request here created #workers number of goroutines. There's no restriction in place
func RunNonSemaphorePool(workers int, supplier func(id int, ch chan<- string)) {
	ch := make(chan string, workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			supplier(id, ch)
		}(i)
	}
	for i := 0; i < workers; i++ {
		val := <-ch
		fmt.Printf("Received signal: %v\n", val)
	}
}

// RunSemaphorePool This pool of workers fan out but only allows a fixed number of goroutines running at
// any given instant of time. There's a latency hit but we restrict huge fanout
func RunSemaphorePool(workers int, concurrent int, supplier func(id int, ch chan<- string)) {
	if workers < concurrent {
		return
	}
	ch := make(chan string, workers)
	// we create a barrier channel (unbuffered to the concurrency) that is only used
	// to signal "done" status hence we use struct{} instead of any concrete value type
	barrier := make(chan struct{}, concurrent)
	for i := 0; i < workers; i++ {
		go func(id int) {
			// each goroutine checks to see if it can write (initially "concurrent" unrestricted)
			barrier <- struct{}{}
			{
				supplier(id, ch)
			}
			// we reset the barrier that we used up.
			<-barrier
		}(i)
	}
	for i := 0; i < workers; i++ {
		val := <-ch
		fmt.Printf("Received signal: %v\n", val)
	}
}
