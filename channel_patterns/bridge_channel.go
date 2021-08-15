package channel_patterns

// Zen: Sometimes we may have to consume from a channel of channels. This could be due
// to pipelines stage outputs etc. The consumer might not care about the fact that the
// values are coming from a chan of channels. This destructuring is called as bridging

// Bridge helps present an abstraction of a single channel from a channel of channels
// Bridge behaves like a multiplexer except that the input is not multiple channels
// as in the case of Or Channel pattern rather a channel of channels.
func Bridge(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	bridge := make(chan interface{})

	// wrapper routine to do most of the work in a separate thread
	go func() {
		defer close(bridge) // don't forget to close once the goroutine exits
		for {
			var stream <-chan interface{}
			// we need to work on one channel at a time, so we keep a variable
			// that points to the current "channel"
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			}
			// now that we have the channel, let's read values from it
			// this maintains the order across each channel
			for v := range OrDone(done, stream) {
				select {
				case <-done:
					return
				case bridge <- v:
					// here we add the value to the bridged channel
				}
			}
		}
	}()
	return bridge
}
