package fan_out_fan_in

import (
	"fmt"
	"math"
	"math/rand"
	gen "patterns/generators"
	"runtime"
	"sync"
)

// Zen: Pipelines are elegant composable stages, but they can be slow, very slow.
// In some situations we fan out to process the input from the stage above in parallel.
// This improves runtime of the stage overall, and it is said to be fanned out.
// Requirements: The stage shouldn't rely on state/values that it has calculated before.
// Requirements: It takes a long time to run to warrant a fan-out.

// PrimeNumberFinderWithNoFanOut demos a stage in which we try to find first 10 primes of a stream
// of random integers. Bound to be slow as this stage is processing it sequentially.
func PrimeNumberFinderWithNoFanOut(count, max, min int) []int {
	// this generates "count" number of prime numbers until
	done := make(chan interface{})
	defer close(done)
	// random number generator isn't order dependent but its way too fast to fan out
	infiniteRandomNumbers := infiniteNumbersStream(done, max, min)
	out := make([]int, 0)

	// now this is the part that is slow is that we process primes in sequence off of the channel
	// while we could be distributing work amongst multiple workers for the prime stage
	primeStream := primeNumbersStream(done, infiniteRandomNumbers, 0)

	for v := range gen.Take(done, primeStream, count) {
		out = append(out, v.(int))
	}
	return out
}

// PrimeNumberFinderWithFanOut is the same as slow but instead spawns (fan out) multiple prime
// number streams based off the parallelism. It then fans them in to create one channel
func PrimeNumberFinderWithFanOut(count, max, min int) []int {
	done := make(chan interface{})
	defer close(done)
	infiniteRandomNumbers := infiniteNumbersStream(done, max, min)
	out := make([]int, 0)

	// now this is the part that is fanned out, we just launch mutiple stages
	numCPU := runtime.NumCPU()
	primeStreamers := make([]<-chan interface{}, numCPU)
	for i := 0; i < numCPU; i++ {
		primeStreamers[i] = primeNumbersStream(done, infiniteRandomNumbers, i)
	}
	// we now need a way to fan in, basically to multiplex multiple prime streams
	// into one channel so that the rest of the streams continues unabated
	primeStreamerFannedIn := FanIn(done, primeStreamers)

	for v := range gen.Take(done, primeStreamerFannedIn, count) {
		out = append(out, v.(int))
	}
	return out
}

// FanIn multiplexes multiple channels into one by draining them concurrently
func FanIn(done <-chan interface{}, channels []<-chan interface{}) chan interface{} {
	multiplexed := make(chan interface{})
	// this wait group is required to drain all channels
	var wg sync.WaitGroup
	multiplex := func(ch <-chan interface{}) {
		// this drains a single channel only
		defer wg.Done()
		for v := range ch {
			select {
			case <-done:
				return
			default:
				multiplexed <- v
			}
		}
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	// to not block the stage, we cannot simply wg.Wait() in this thread
	// we need to wait in a new goroutine
	go func() {
		wg.Wait()
		// remember to close the multiplexed channel
		close(multiplexed)
	}()
	return multiplexed
}

// this stage just returns random numbers within a range using the RepeatWithFn primitive
func infiniteNumbersStream(done chan interface{}, max int, min int) <-chan int {
	return gen.ToInt(done, gen.RepeatWithFn(done, func() interface{} {
		return rand.Intn(max-min) + min
	}))
}

// this stage receives input from a channel and sequentially tries to filter
// prime numbers off of the list
func primeNumbersStream(done chan interface{}, nums <-chan int, streamId int) <-chan interface{} {
	fmt.Printf("Prime Streamer [%d] started\n", streamId)
	primes := make(chan interface{})
	go func() {
		defer close(primes)
		for v := range nums {
			select {
			case <-done:
				return
			default:
				if isPrime(v) {
					primes <- v
				}
			}
		}
	}()
	return primes
}

// helper function that does a naive test using the repeated square root division
// Although a more performant "big.NewInt(num).ProbablyPrime(0)" is available, we avoid it
func isPrime(n int) bool {
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
