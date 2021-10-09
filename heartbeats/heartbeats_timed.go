package heartbeats

import "time"

// Zen: Heartbeats are a way for concurrent processes to signal life to outside properties.
// They can occur on a timed interval (useful for concurrent code waiting for something else
// to happen for it to process a unit of work). Because we don't know when that work might
// come in, our goroutine might be sitting around waiting for something to happen. Thus, a
// heartbeat is a way to signal to its listeners that all is well and silence is expected.

// BasicHeartbeatAndResult is a function that provides a channel which is signalled every
// pulse interval seconds along with a result channel that is signalled at double the interval
func BasicHeartbeatAndResult(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
	heartbeatCh := make(chan interface{})
	resultCh := make(chan time.Time) // result channel could be on anything, we just send time

	go func() {
		defer close(heartbeatCh)
		defer close(resultCh)
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(pulseInterval * 2) // we choose twice the interval arbitrarily

		sendWorkResult := func(t time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse: // just like done, we also need to include a case for pulse
					sendPulse(heartbeatCh)
				case resultCh <- t: // signal the actual result and return, our work is done
					return
				}
			}
		}
		// we might be sending out multiple pulses while waiting to send results, hence, for loop.
		for {
			select {
			case <-done:
				return
			case <-pulse: // send a pulse when the pulse ticker signals
				sendPulse(heartbeatCh)
			case t := <-workGen: // send result when the result ticker signals
				sendWorkResult(t)
			}
		}
	}()
	return heartbeatCh, resultCh
}

func sendPulse(heartbeatCh chan<- interface{}) {
	select {
	case heartbeatCh <- struct{}{}:
	default:
		// we must guard against the fact that no one might be listening to our pulse
		// results emitted are critical, pulses are not
	}
}

// BasicHeartbeatAndResultFaulty is same as BasicHeartbeatAndResult but fails after two iterations
// and doesn't close its channel, which results in a panic. This will be detected as no pulse and
// the main goroutine can take appropriate action
func BasicHeartbeatAndResultFaulty(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
	heartbeatCh := make(chan interface{})
	resultCh := make(chan time.Time) // result channel could be on anything, we just send time

	go func() {
		// we forget to close the channel, resulting in a panic
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(pulseInterval * 2) // we choose twice the interval arbitrarily

		sendWorkResult := func(t time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse: // just like done, we also need to include a case for pulse
					sendPulse(heartbeatCh)
				case resultCh <- t: // signal the actual result and return, our work is done
					return
				}
			}
		}
		// we might be sending out multiple pulses while waiting to send results, hence, for loop.
		// but we want to simulate an error, hence we break after two iterations
		for i := 0; i < 2; i++ {
			select {
			case <-done:
				return
			case <-pulse: // send a pulse when the pulse ticker signals
				sendPulse(heartbeatCh)
			case t := <-workGen: // send result when the result ticker signals
				sendWorkResult(t)
			}
		}
	}()
	return heartbeatCh, resultCh
}
