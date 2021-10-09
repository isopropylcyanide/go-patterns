package heartbeats

import "time"

// Zen: Heartbeats are a way for concurrent processes to signal life to outside properties.
// They can occur at the beginning of a unit of work. These are extremely useful for tests.
// They let us know that long-running goroutines remain up, and are only slow (but will run)

// HeartbeatGenerateIntStream is a function that provides a channel which is signalled every
// time a unit of work is done. Unit of work here is to generate nums to a stream on a channel
func HeartbeatGenerateIntStream(done <-chan interface{}, sleep time.Duration, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeatCh := make(chan interface{}, 1)
	// ensure at least one pulse is sent even if no one is listening in time for the event to occur
	intStream := make(chan int)

	go func() {
		defer close(heartbeatCh)
		defer close(intStream)

		// we simulate some delay, which in reality could be anything
		// this makes the life of a test very hard as it now needs to choose to wait a time
		// if it is too high, failures will take a longer time, if too less, it is flaky
		time.Sleep(sleep)

		for _, n := range nums {
			select {
			// A good test will rely on our heartbeat to see if it should wait or not as we signal
			// that our work has begun. It may take much longer to process the iteration but at
			// least we have started processing
			case heartbeatCh <- struct{}{}:
			// we send a pulse before work, everytime
			// this is not in the same select as int stream because if the receiver
			// isn't ready for result, we don't want to send a pulse in which case
			// thr result would be lost
			default:
				// fallthrough if no one is listening, but because it is s a buffered
				// channel, even if someone missed the first pulse, they'll get notified
			}
			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()
	return heartbeatCh, intStream
}
