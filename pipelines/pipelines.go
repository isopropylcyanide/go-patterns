package pipelines

import "fmt"

// Zen: By using a pipeline, you separate the concerns of each stage, which provides
// numerous benefits. You can modify stages independent of each other, process them
// concurrently, fan-out or even rate limit individual stages thereby improving flexibility
// Requirements: Each stage consumes and returns the same type
// Requirements: Stage must be reified by the language so that they can be passed around

// RudimentaryBatchPipeline A batch pipeline where we process inputs in chunks at once.
// Each function returns and consumes a slice of data and not discrete elements
func RudimentaryBatchPipeline() []int {
	add := func(list []int, additive int) []int {
		res := make([]int, len(list))
		// using range here means the memory footprint is high, but easier for the caller
		for i, v := range list {
			res[i] = v + additive
		}
		return res
	}
	multiply := func(list []int, multiplier int) []int {
		res := make([]int, len(list))
		for i, v := range list {
			res[i] = v * multiplier
		}
		return res
	}
	input := []int{1, 2, 3, 4}
	output := multiply(add(multiply(input, 2), 1), 2)
	fmt.Println("Batch ", output)
	return output
}

// RudimentaryStreamPipeline A stream pipeline where we process inputs one at a time
// Each function returns and consumes a discrete value
func RudimentaryStreamPipeline() []int {
	add := func(input int, additive int) int {
		return input + additive
	}
	multiply := func(input int, multiplier int) int {
		return input * multiplier
	}
	input := []int{1, 2, 3, 4}
	outputs := make([]int, 0)
	for _, v := range input {
		// this range loop limits our ability to scale and feed the pipeline
		// we are also making multiple function calls for each iteration
		output := multiply(add(multiply(v, 2), 1), 2)
		fmt.Println("Streaming ", output)
		outputs = append(outputs, output)
	}
	return outputs
}

// ChannelStreamPipeline Channels are suited to pipelining as they can receive and signal values,
// are safe to use concurrently, can be ranged over and are reified by Go
func ChannelStreamPipeline(done chan interface{}) <-chan int {
	// the done channel is given by callers which know when to stop the pipeline
	// if we create it here, we need to close it, which means caller will never see the output

	// a generator is used to convert input into a channel that signals input (write)
	// a read only done channel is used to know when to stop
	// we return a read only channel because callers are only going to read
	generator := func(done <-chan interface{}, input ...int) <-chan int {
		ch := make(chan int)
		go func() {
			// don't forget to close the channel
			defer close(ch)
			for _, v := range input {
				select {
				// we need to know when to stop
				case <-done:
					return
				case ch <- v:
				}
			}
		}()
		return ch
	}

	// add is now a function that operates on a read only channel of values (with done) and
	// returns its output on a write channel of values
	add := func(input <-chan int, done <-chan interface{}, additive int) <-chan int {
		res := make(chan int)
		go func() {
			defer close(res)
			for v := range input {
				select {
				case <-done:
					return
				case res <- v + additive:
				}
			}
		}()
		return res
	}

	// add is now a function that operates on a read only channel of values (with done) and
	// returns its output on a write channel of values
	multiply := func(input <-chan int, done <-chan interface{}, multiplier int) <-chan int {
		res := make(chan int)
		go func() {
			defer close(res)
			for v := range input {
				select {
				case <-done:
					return
				case res <- v * multiplier:
				}
			}
		}()
		return res
	}
	intStream := generator(done, 1, 2, 3, 4)
	return multiply(add(multiply(intStream, done, 2), done, 1), done, 2)
}
