package channel_patterns

// Zen: When you do not know the number of channels (to multiplex into one) in advance,
// you cannot use a static select statement. Instead, use something like an "Or channel"
// that recursively runs a select operation.

// Or Multiplexes multiple channels into one channel that closes if any of its component
// channels close. Useful, when you don't know the number of channels in advance
// the input ch and the output ch are read only
func Or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orChannel := make(chan interface{})
	go func() {
		defer close(orChannel)
		switch len(channels) {
		case 2:
			// this case block and switch can be removed. It is only
			// a minor optimization to avoid recursion overhead (always 2 channels)
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-Or(append(channels[3:], orChannel)...):
			}
		}
	}()
	return orChannel
}
