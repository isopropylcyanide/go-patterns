package handy_generators

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.uber.org/goleak"

	"github.com/stretchr/testify/assert"
)

func TestRepeatGeneratorDemo(t *testing.T) {
	done := make(chan interface{})
	// here we would Repeat forever. to curb this, lets close channel in an another goroutine
	// that will let the main run for sometime until "it" (not main) closes channel
	go func() {
		time.Sleep(100 * time.Microsecond)
		close(done)
	}()

	actualInts := make([]int, 0)
	for v := range Repeat(done, 1, 2, 3, 4) {
		fmt.Printf("Repeat %v -> \n", v)
		actualInts = append(actualInts, v.(int))
	}
	assert.True(t, len(actualInts) > 0)
}

func TestTakeGeneratorDemo(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	actualInts := make([]int, 0)
	expectedInts := []int{1, 2, 3, 4, 1, 2, 3, 4, 1, 2}

	for v := range Take(done, Repeat(done, 1, 2, 3, 4), 10) {
		fmt.Printf("RepeatTake %v -> \n", v)
		actualInts = append(actualInts, v.(int))
	}
	assert.True(t, reflect.DeepEqual(actualInts, expectedInts))
}

func TestRepeatFunctionWithTakeDemo(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	sum, count := 0, 0
	fn := func() interface{} {
		sum += count
		count += 1
		return sum
	}
	actualInts := make([]int, 0)
	expectedInts := []int{0, 1, 3, 6, 10}

	for v := range Take(done, RepeatWithFn(done, fn), 5) {
		fmt.Printf("RepeatTakeFn %v -> \n", v)
		actualInts = append(actualInts, v.(int))
	}
	assert.True(t, reflect.DeepEqual(actualInts, expectedInts))
}

func TestToStringRepeatFunctionWithTakeDemo(t *testing.T) {
	// This demonstrates the usage of ToString applied on a list of input. We'll also
	// write a benchmark to prove that adding to string doesn't add lot of overhead
	done := make(chan interface{})
	defer close(done)

	actualStrings := make([]string, 0)
	expectedStrings := []string{"H1", "H3", "H4", "H1", "H3"}

	for v := range ToString(done, Take(done, Repeat(done, "H1", "H3", "H4"), 5)) {
		fmt.Printf("ToStringRepeatTake %v -> \n", v)
		actualStrings = append(actualStrings, v)
	}
	assert.True(t, reflect.DeepEqual(actualStrings, expectedStrings))
}

func TestToIntRepeatFunctionWithTakeDemo(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	actualInts := make([]int, 0)
	expectedInts := []int{3, 7, 5, 3, 7}

	for v := range ToInt(done, Take(done, Repeat(done, 3, 7, 5), 5)) {
		fmt.Printf("ToIntRepeatTake %v -> \n", v)
		actualInts = append(actualInts, v)
	}
	assert.True(t, reflect.DeepEqual(actualInts, expectedInts))
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

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
