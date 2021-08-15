package for_select_loop

import (
	"fmt"
	"time"
)

// Zen: A ubiquitous pattern that does a select on the channel and proceeds based on
// whichever channel were to react first. Can also be used to block forever

func SendIterationValuesOnChannel() {
	done := make(chan struct{})
	stringStream := make(chan string, 3)

	go func() {
		for _, v := range []string{"A", "B", "C"} {
			select {
			case <-done:
				return
			case stringStream <- v:
			}
		}
		close(stringStream)
	}()

	for str := range stringStream {
		fmt.Println("Got ", str)
	}
}

func InfiniteLooping() {
	done := make(chan struct{})
	for {
		select {
		case <-done:
			// until done is passed a value, this will exit the select block and loop
			return
		default:
		}
		// do work here
		fmt.Println("Looping")
	}
}

func InfiniteLoopingII() {
	done := make(chan struct{})
	for {
		select {
		case <-done:
			return
		case <-time.After(50 * time.Millisecond):
			fmt.Println("okay")
			return
		default:
			// do work here in default (also an option) however aim for less indentation
			fmt.Println("Looping II")
		}
	}
}
