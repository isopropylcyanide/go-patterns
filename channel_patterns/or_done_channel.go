package channel_patterns

// Zen: Unlike pipelines, you can't make any assertions about how a channel would behave
// when code you're working with is cancelled via its done channel. You don't know if the
// fact that your goroutine has been cancelled means the channel you're reading from will
// have been cancelled. Hence, we need to wrap our loop around "done" channel. To make
// sure stuff isn't verbose for every loop, we'll create a dedicated or_done_channel

// OrDone provides a pattern to read from a channel until done. Refer the naive version
// in tests to see the usage without this pattern
func OrDone(done, c <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case val, ok := <-c:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
				}
			}
		}
	}()
	return out
}
