package healing_goroutines

import (
	"log"
	"patterns/channel_patterns"
	"time"

	"go.uber.org/atomic"
)

// Zen: In long-lived processes such as daemons, it's common to have a long-lived set of
// goroutines. These may be blocked, live locked or doing useful work. It may be useful
// to establish a mechanism that ensures the goroutines remain healthy and are restarted
// if they become unhealthy. This is called healing.
//
// Normally, separation of concerns
// dictates that it's the work of a supervisor process (as in Erlang) to heal. However,
// it's a good pattern to know for simpler architectures
// Here, the healing goroutine is called as "steward" and the goroutine it heals is a "ward".

// haltCounter tracks the number of times the ward goroutine is halted
var haltCounter atomic.Int64

// doIrresponsibleWork is simply waiting on the input channel. It's not progressing and
// nor is it sending any pulses. Running this simply will wait until the channel is closed
func doIrresponsibleWork(done <-chan interface{}, _ time.Duration) <-chan interface{} {
	out := make(chan interface{})
	log.Println("ward: I'm the ward. I can be irresponsible and become unhealthy")
	go func() {
		defer close(out)
		<-done
		log.Println("ward: Halting")
		haltCounter.Add(1)
	}()
	return out
}

// startFn is a signature of a function that can be started, closed and monitored
// It sends heart beats pulses at a duration specified by "pulseInterval"
type startFn func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{})

// NewSteward is a monitoring goroutine that takes a start function and a timeout.
// If the ward doesn't reply with a healthy heartbeat to the steward, it will time out and
// restart the goroutine. The steward itself returns a startFn, so it can be monitored too.
func newSteward(timeout time.Duration, f startFn) startFn {
	return func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
		heartbeat := make(chan interface{})
		go func() {
			defer close(heartbeat)
			var wardDone chan interface{}
			var wardHeartbeat <-chan interface{}
			// closure to ensure a consistent way to start the ward
			startWard := func() {
				// the wardDone channel becomes the only signal for the ward to continue
				wardDone = make(chan interface{})
				// using the multiplexed or channel pattern, we want the ward to halt if
				// the steward is halted or if the steward wants the ward to halt
				wardHeartbeat = f(channel_patterns.Or(wardDone, done), timeout/2) // extra tick for heartbeat to respond to
			}
			startWard()
			pulse := time.Tick(pulseInterval)

		monitorLoop:
			for { // forever monitoring loop
				timeoutSignal := time.After(timeout) // when steward decides it's enough
				for {
					select {
					case <-pulse: // wait for a heartbeat from a ward
						select {
						case heartbeat <- struct{}{}: // steward is healthy
						default:
						}
					case <-wardHeartbeat:
						log.Println("steward: got ward heartbeat, ")
						continue monitorLoop
					case <-timeoutSignal:
						log.Println("steward: ward unhealthy, restarting")
						close(wardDone) // tell ward to stop since it's not behaving properly
						startWard()
						continue monitorLoop
					case <-done:
						return
					}
				}
			}
		}()
		return heartbeat
	}
}
