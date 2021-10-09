package heartbeats

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestHeartbeatWithResult checks a happy case where the worker goroutine is never unhealthy
// which is determined by the lack of it sending a heartbeat after a certain grace period.
func TestHeartbeatWithResult(t *testing.T) {
	t.Parallel()
	done := make(chan interface{})
	// close done channel after 6 seconds to give some room for the routine to work
	time.AfterFunc(4*time.Second, func() { close(done) })

	// we are waiting for upto two seconds until we pronounce the goroutine as unhealthy
	const timeout = time.Second * 1

	// when there are no results, we are at least guaranteed to get a heartbeat every t/2
	// if we do not receive it, something is wrong with the goroutine
	pulses, results := HeartbeatAndResult(done, timeout/2)

	for {
		select {
		case _, ok := <-pulses:
			if !ok {
				// no more heartbeats, we can return
				return
			}
			fmt.Println("Got pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("Got result %v\n", r)
		case <-time.After(timeout):
			t.Fatal("Worker goroutine unhealthy")
		}
	}
}

// TestHeartbeatWithResultUnhealthyIsDetected is same as TestBasicHeartbeatWithResultUnhealthy,
// but we detect a failure (if we see no heartbeat) this way we avoid a deadlock and do not
// have to rely on a longer timeout
func TestHeartbeatWithResultUnhealthyIsDetected(t *testing.T) {
	t.Parallel()
	//
	done := make(chan interface{})
	time.AfterFunc(4*time.Second, func() { close(done) })
	const timeout = time.Second * 1

	pulses, results := HeartbeatAndResultFaulty(done, timeout/2)
	failureDetected := false
L:
	for {
		select {
		case _, ok := <-pulses:
			if !ok {
				// no more heartbeats, we can return
				return
			}
			fmt.Println("Got pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("Got result %v\n", r)
		case <-time.After(timeout):
			failureDetected = true
			// this is detected immediately, note we didn't have to depend on done channel
			// which would close after 6 seconds, we just return after two seconds
			fmt.Println("Worker goroutine is not healthy")
			// break is important otherwise we'll detect the failure but not do anything about it
			break L
		}
	}
	assert.True(t, failureDetected)
}
