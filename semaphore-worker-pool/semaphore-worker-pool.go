package main

import (
	"fmt"
	"time"
)

// Zen: When fanning out to a high number of goroutines that write on an unbuffered channel
// (to improve latency) make sure that the count is not too high. For such high fanouts,
// ensure the use of a semaphore that doesn't allow the creation of a small set of routines.

// This pool of workers fan out and create as many goroutines as there are workers
// Each request here created #workers number of goroutines. There's no restriction in place
func nonSemaphorePool(workers int, supplier func(id int, ch chan<- string)) {
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

// This pool of workers fan out but only allows a fixed number of goroutines running at
// any given instant of time. There's a latency hit but we restrict huge fanout
func semaphorePool(workers int, concurrent int, supplier func(id int, ch chan<- string)) {
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

func main() {
	work := func(id int, ch chan<- string) {
		//time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		time.Sleep(1000 * time.Millisecond)
		ch <- fmt.Sprintf("Done #%d", id)
	}
	maxWorkers := 20
	start := time.Now()
	nonSemaphorePool(maxWorkers, work)
	fmt.Printf("Non semaphore workers took: %v\n", time.Since(start).Seconds())

	start = time.Now()
	maxConcurrency := 5
	// this will be slower but has better flexibility on goroutines
	// you'll see work being performed in batches
	semaphorePool(maxWorkers, maxConcurrency, work)
	fmt.Printf("Semaphore restricted workers took: %v\n", time.Since(start).Seconds())
}
