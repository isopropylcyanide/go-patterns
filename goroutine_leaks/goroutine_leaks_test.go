package goroutine_leaks

import (
	"testing"

	"go.uber.org/goleak"
)

func TestLeakGoRoutineBlockedOnReading(t *testing.T) {
	// this will leak and hence has been added to the go-leak ignore list
	LeakGoRoutineBlockedOnReading()
}

func TestLeakGoRoutineBlockedOnReadingDeadlock(t *testing.T) {
	// Uncomment this to see causing the deadlock
	t.Skipf("Skipping this test as it will cause the deadlock")
	LeakGoRoutineBlockedOnReadingDeadlock()
}

func TestLeakGoRoutineBlockedOnReadingFixedUsingDoneChannel(t *testing.T) {
	// this should not leak (as captured by go-leak)
	LeakGoRoutineBlockedOnReadingFixedUsingDoneChannel()
}

func TestLeakGoRoutineBlockedOnWriting(t *testing.T) {
	LeakGoRoutineBlockedOnWriting()
}

func TestLeakGoRoutineBlockedOnWritingFixedUsingDoneChannel(t *testing.T) {
	LeakGoRoutineBlockedOnWritingFixedUsingDoneChannel()
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m,
		goleak.IgnoreTopFunction("patterns/goroutine_leaks.LeakGoRoutineBlockedOnReading.func1.1"), // read goroutine leak assert
		goleak.IgnoreTopFunction("patterns/goroutine_leaks.LeakGoRoutineBlockedOnWriting.func1.1"), // write goroutine leak assert
	)
}
