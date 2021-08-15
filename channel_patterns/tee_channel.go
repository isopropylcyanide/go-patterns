package channel_patterns

// Zen: Use Tee to split an incoming value from a channel into multiple channels
// Usually used to carry out two independent units of work from an input

// Tee returns two separate channels from where a single input value can be read separately
// Tee behaves like a de-multiplexer
func Tee(done <-chan interface{}, input <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	teeA := make(chan interface{})
	teeB := make(chan interface{})

	go func() {
		defer close(teeA) // don't forget to close once the goroutine exits
		defer close(teeB)
		for v := range input {
			// iteration over input is coupled to both to writes in teeA and teeB
			// usually not a problem

			// the problem here is to send item to each tee once but select doesn't
			// ensure it, so we rely on the fact that sending to a nil channel blocks
			// once we write to a channel, say teeA, we make sure we set it to nil, so
			// that any write to it will block giving time for select to choose teeB
			teeA, teeB := teeA, teeB
			// we need to set their local copies to nil
			for i := 0; i < 2; i++ {
				// we loop twice so that select will work for each channel only once
				select {
				case <-done:
					return
				case teeA <- v:
					// if teeA is chosen at i = 0, we set teeA = nil so that at i = 1, teeB is chosen
					teeA = nil
				case teeB <- v:
					// if teeB is chosen at i = 0, we set teeB = nil so that at i = 1, teeA is chosen
					teeB = nil
				}
			}
		}
	}()
	return teeA, teeB
}
