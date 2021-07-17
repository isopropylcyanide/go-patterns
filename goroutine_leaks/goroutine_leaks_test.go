package goroutine_leaks

import "testing"

func TestLeakGoRoutineBlockedOnReading(t *testing.T) {
	LeakGoRoutineBlockedOnReading()
}

func TestLeakGoRoutineBlockedOnReadingDeadlock(t *testing.T) {
	// Comment this to avoid causing the deadlock
	LeakGoRoutineBlockedOnReadingDeadlock()
}

func TestLeakGoRoutineBlockedOnWriting(t *testing.T) {
	LeakGoRoutineBlockedOnWriting()
}

func TestLeakGoRoutineBlockedOnWritingFixedUsingDoneChannel(t *testing.T) {
	LeakGoRoutineBlockedOnWritingFixedUsingDoneChannel()
}
