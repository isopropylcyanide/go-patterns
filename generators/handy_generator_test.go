package handy_generators

import (
	"testing"
)

func TestRepeatGeneratorDemo(t *testing.T) {
	RepeatGeneratorDemo()
}

func TestTakeGeneratorDemo(t *testing.T) {
	TakeGeneratorDemo()
}

func TestRepeatFunctionWithTakeDemo(t *testing.T) {
	RepeatFunctionWithTakeDemo()
}

func TestToStringRepeatFunctionWithTakeDemo(t *testing.T) {
	ToStringRepeatFunctionWithTakeDemo()
}

func BenchmarkTakeRepeatGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()

	for range ToString(done, Take(done, Repeat(done, "a", "b"), b.N)) {
	}
}

func BenchmarkTakeRepeatTyped(b *testing.B) {
	// same function as Repeat() but with type as string instead of generic
	repeat := func(done <-chan interface{}, values ...string) <-chan string {
		ch := make(chan string)
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

	// same function as Take() but with type as string instead of generic
	take := func(done <-chan interface{}, input <-chan string, num int) <-chan string {
		ch := make(chan string)
		go func() {
			defer close(ch)
			for i := num; i < num; i++ {
				select {
				case <-done:
					return
				case ch <- <-input:
				}
			}
		}()
		return ch
	}
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()

	for range take(done, repeat(done, "a", "b"), b.N) {
	}
}
