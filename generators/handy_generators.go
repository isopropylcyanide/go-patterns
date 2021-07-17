package handy_generators

// Zen: A generator for a pipeline is any function that converts a set of discrete set of values
// into a stream of values on a channel. Using channels / done idiom, we can generate efficient
// generators

// Repeat repeats the values you pass to it indefinitely
func Repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case ch <- v:
				}
			}
		}
	}()
	return ch
}

// RepeatWithFn repeats the values indefinitely after applying a function
func RepeatWithFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case ch <- fn():
			}
		}
	}()
	return ch
}

// Take takes a finite set of values from a given channel represented by the number
// or returns all the elements in the channel if the count is lesser
func Take(done <-chan interface{}, input <-chan interface{}, num int) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case ch <- <-input:
			}
		}
	}()
	return ch
}

// ToString takes an input channel of type interface and converts the values into its
// string type using cast
func ToString(done <-chan interface{}, input <-chan interface{}) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for v := range input {
			select {
			case <-done:
				return
			case ch <- v.(string):
			}
		}
	}()
	return ch
}

// ToInt takes an input channel of type interface and converts the values into its
// int type using cast
func ToInt(done <-chan interface{}, input <-chan interface{}) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for v := range input {
			select {
			case <-done:
				return
			case ch <- v.(int):
			}
		}
	}()
	return ch
}
