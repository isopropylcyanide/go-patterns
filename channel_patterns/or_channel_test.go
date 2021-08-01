package channel_patterns

import (
	"fmt"
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	start := time.Now()
	<-Or(
		signalIdAfter(1, 1*time.Minute),
		signalIdAfter(2, 3*time.Second),
		signalIdAfter(3, 4*time.Second),
		signalIdAfter(4, 3*time.Hour),
		signalIdAfter(5, 1*time.Second),
	)
	fmt.Printf("Or channel finished in %v\n", time.Since(start))
}

// signalIdAfter returns a channel after sleeping for given time units
func signalIdAfter(id int, after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		defer fmt.Printf("Closing channel[%d]\n", id)
		time.Sleep(after)
	}()
	return c
}
