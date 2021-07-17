package handy_generators

import (
	"fmt"
	"math/rand"
	"time"
)

// Zen: A generator for a pipeline is any function that converts a set of discrete set of values
// into a stream of values on a channel. Using channels / done idiom, we can generate efficient
// generators

// Repeat repeats the values you pass to it indefinitely
func Repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case ch <- v:
				}
			}
		}
	}()
	return ch
}

// RepeatWithFn repeats the values indefinitely after applying a function
func RepeatWithFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case ch <- fn():
			}
		}
	}()
	return ch
}

// Take takes a finite set of values from a given channel represented by the number
// or returns all the elements in the channel if the count is lesser
func Take(done <-chan interface{}, input <-chan interface{}, num int) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case ch <- <-input:
			}
		}
	}()
	return ch
}

// ToString takes an input channel of type interface and converts the values into its
// string type using cast
func ToString(done <-chan interface{}, input <-chan interface{}) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for v := range input {
			select {
			case <-done:
				return
			case ch <- v.(string):
			}
		}
	}()
	return ch
}

// RepeatGeneratorDemo demonstrates the usage of Repeat
func RepeatGeneratorDemo() {
	done := make(chan interface{})
	// here we would Repeat forever. to curb this, lets close channel in an another goroutine
	// that will let the main run for sometime until "it" (not main) closes channel
	go func() {
		time.Sleep(100 * time.Microsecond)
		close(done)
	}()

	for v := range Repeat(done, 1, 2, 3, 4) {
		fmt.Printf("Repeat %v -> \n", v)
	}
}

// TakeGeneratorDemo demonstrates the usage of Take
func TakeGeneratorDemo() {
	done := make(chan interface{})
	defer close(done)

	for v := range Take(done, Repeat(done, 1, 2, 3, 4), 10) {
		fmt.Printf("Take %v -> \n", v)
	}
}

// RepeatFunctionWithTakeDemo demonstrates the usage of Take on a RepeatFWithFn
func RepeatFunctionWithTakeDemo() {
	done := make(chan interface{})
	defer close(done)

	fn := func() interface{} {
		return rand.Int()
	}
	for v := range Take(done, RepeatWithFn(done, fn), 4) {
		fmt.Printf("TakeRepeatFn %v -> \n", v)
	}
}

// ToStringRepeatFunctionWithTakeDemo demonstrates the usage of ToString applied on a list of input
// We'll also write a benchmark to prove that adding to string doesn't add lot of overhead
func ToStringRepeatFunctionWithTakeDemo() {
	done := make(chan interface{})
	defer close(done)

	for v := range ToString(done, Take(done, Repeat(done, "H1", "H3", "H4"), 5)) {
		fmt.Printf("ToStringTakeRepeat %v -> \n", v)
	}
}
