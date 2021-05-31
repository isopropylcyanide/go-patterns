package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Zen: If a goroutine creates another goroutine, it's also responsible for ensuring
// that it can stop the started goroutine

// Shows a go func that never ends and is essentially leaked
func leakGoRoutineBlockedOnReading() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("Do Work Exited")
			defer close(completed)
			for s := range strings {
				fmt.Println("DoWork: ", s)
				completed <- s
			}
		}()
		return completed
	}
	// by passing nil, the go func() will essentially block during read
	doWork(nil)
}

// this is the same as leakGoRoutineBlockedOnReading but proves that the goroutine runs forever
// when there's a deadlock in the main waiting for the go func to complete
func leakGoRoutineBlockedOnReadingDeadlock() {
	//var wg sync.WaitGroup
	//wg.Add(1)
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("Do Work Exited")
			defer close(completed)
			for s := range strings {
				fmt.Println("DoWork: ", s)
				completed <- s
			}
		}()
		return completed
	}
	results := doWork(make(chan string))
	<-results
}

// We fix the leak by using a done channel that fixes the leak by satisfying the
// goroutine's select call in which case it returns after seeing a done value
func leakGoRoutineBlockedOnReadingFixedUsingDoneChannel() {
	// done is a read only channel
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		results := make(chan interface{})
		go func() {
			defer fmt.Println("Do Work Exited")
			defer close(results)
			for {
				select {
				case <-done:
					return
				case s := <-strings:
					fmt.Println("DoWork: ", s)
					results <- s
				}
			}
		}()
		return results
	}
	// by passing nil, the go func() will essentially block during read
	done := make(chan interface{})
	result := doWork(done, nil)

	// in a separate routine, we wait for a few seconds before deciding it is enough
	// and then we cancel the done channel signalling the goroutine above to return
	go func(tolerance int) {
		fmt.Println("will wait for a maximum of ", tolerance, " seconds")
		time.Sleep(time.Duration(tolerance) * time.Second)
		fmt.Println("cancelling the go routine that would have leaked otherwise")
		close(done)
	}(1)

	<-result
	fmt.Println("done")
}

// This go routine writes something on the channel but no one closes it
// after the main has read some items off the channel
func leakGoRoutineBlockedOnWriting() {
	doWork := func() <-chan int {
		results := make(chan int)
		go func() {
			// you will see that no one closes the channel
			defer close(results)
			defer fmt.Println("closing the infinite generation")
			for {
				results <- rand.Int()
			}
		}()
		return results
	}
	results := doWork()
	for i := 0; i < 3; i++ {
		fmt.Printf("main receives %d\n", <-results)
	}
	fmt.Println("Done")
}

// Same as leakGoRoutineBlockedOnWriting but there's no leak here because the main closes
// the channel and the random generator goroutine knows its time to stop
func leakGoRoutineBlockedOnWritingFixedUsingDoneChannel() {
	doWork := func(done <-chan interface{}) <-chan int {
		results := make(chan int)
		go func() {
			// you will see that main closes the channel using done
			defer close(results)
			defer fmt.Println("closing the infinite generation")
			for {
				select {
				case <-done:
					return
				case results <- rand.Int():
				}
			}
		}()
		return results
	}
	done := make(chan interface{})
	results := doWork(done)
	for i := 0; i < 3; i++ {
		fmt.Printf("main receives %d\n", <-results)
	}
	fmt.Println("done. closing channel which should close the child channel")
	close(done)
	// wait for some time so that you see the child defer
	time.Sleep(100 * time.Millisecond)
}

func main() {
	leakGoRoutineBlockedOnReading()
	// Uncomment below to see the deadlock
	// leakGoRoutineBlockedOnReadingDeadlock()
	leakGoRoutineBlockedOnReadingFixedUsingDoneChannel()

	leakGoRoutineBlockedOnWriting()
	leakGoRoutineBlockedOnWritingFixedUsingDoneChannel()
}
