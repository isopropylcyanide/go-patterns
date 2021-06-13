package main

import (
	"fmt"
	"time"
)

// or Multiplexes multiple channels into one channel that closes if any of its component
// channels close. Useful, when you don't know the number of channels in advance
// the input ch and the output ch are read only
func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orChannel := make(chan interface{})
	go func() {
		defer close(orChannel)
		switch len(channels) {
		case 2:
			// this case block and switch can be removed. It is only
			// a minor optimization to avoid recursion overhead (always 2 channels)
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orChannel)...):
			}
		}
	}()
	return orChannel
}

// signalAfter returns a channel after sleeping for given time units
func signalAfter(id int, after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		defer fmt.Printf("Closing channel[%d]\n", id)
		time.Sleep(after)
	}()
	return c
}

// demoOrChannelPlain tests the working of the or channel by passing
func demoOrChannelPlain() {
	start := time.Now()
	<-or(
		signalAfter(1, 1*time.Minute),
		signalAfter(2, 3*time.Second),
		signalAfter(3, 4*time.Second),
		signalAfter(4, 3*time.Hour),
		signalAfter(5, 1*time.Second),
	)
	fmt.Printf("or channel finished in %v\n", time.Since(start))
}

func main() {
	demoOrChannelPlain()
}
