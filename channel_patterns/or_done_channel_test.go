package channel_patterns

import (
	"fmt"
	handy_generators "patterns/generators"
	"testing"
	"time"
)

func TestOrDone(t *testing.T) {
	done := signalAfter(100 * time.Microsecond)
	items := handy_generators.Repeat(done, 1, 2, 3)

	for val := range OrDone(done, items) {
		fmt.Printf("Received %v\n", val)
	}
}

func TestOrDoneNaive(t *testing.T) {
	done := signalAfter(100 * time.Microsecond)
	items := handy_generators.Repeat(done, 1, 2, 3)

	for {
		// this is boilerplate and will be unwieldy in a nested loop scenario
		select {
		case <-done:
			return
		case val, ok := <-items:
			if !ok {
				return
			}
			fmt.Printf("Received %v\n", val)
		}
	}
}

// signalAfter returns a channel after sleeping for given time units
func signalAfter(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		defer fmt.Printf("Closing done channel\n")
		time.Sleep(after)
	}()
	return c
}
