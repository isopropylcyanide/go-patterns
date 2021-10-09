package replicated_requests

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Zen: Sometimes, to service a request that can be served through multiple ways, we want to
// return the fastest way the result can be made possible. Generally, there may be different
// ways to do it, but for single process, we can call multiple goroutines and gather the first
// result. This pattern is called replicated requests.

// Note that for this to work, the resources required to service the request should be replicated
// too or else the handlers will not have an equal opportunity to service the requests.
// This technique is expensive. If all paths are uniform (with no major probability of outliers),
// replicating is wasteful. At the cost of speed and fault tolerance, we are trading resource
// utilization. Hence, it is better to use when there are multiple access patterns / paths to service
// the resource

// Note that this is different from the "fan out" pattern wherein we fan out to multiple routines
// and coalesce their results later, as here we only care about the first result, always

// DoWork processes a request with a given id with a random delay. This is equivalent to a
// handler. The calling code will spawn multiple instances of the handler
func DoWork(done <-chan interface{}, id int, wg *sync.WaitGroup, result chan<- int) {
	started := time.Now()
	defer wg.Done()

	// simulate random delay
	delay := time.Duration(1+rand.Intn(5)) * time.Second

	select {
	case <-done:
		fmt.Printf("handler %v has been cancelled pre sleep of %v seconds\n", id, delay)
		return // we have been cancelled
	case <-time.After(delay):
	}

	select {
	case <-done: // we might be cancelled, even after sleep
		fmt.Printf("handler %v has been cancelled post sleep of %v seconds\n", id, delay)
	case result <- id:
	}

	took := time.Since(started)
	if took < delay {
		took = delay
	}
	fmt.Printf("handler %v took %v\n", id, took)
}
