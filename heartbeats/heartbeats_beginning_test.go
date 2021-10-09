package heartbeats

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/goleak"

	"github.com/stretchr/testify/assert"
)

// TestBasicHeartbeatGenerateIntStream_BadTest is a bad test because it relies on a timeout
// and is flaky and to prove this we provide a value to sleep on, and we wait only a second.
// But in reality, we don't know what the sleep duration of the called routine would be
func TestBasicHeartbeatGenerateIntStream_BadTest(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	ints := []int{0, 1, 2, 3, 5}
	// our test is not going to wait any longer
	const testWaitingPeriod = time.Second

	cases := []struct {
		shouldTimeout bool
		sleepFor      time.Duration
	}{
		{
			shouldTimeout: false,
			sleepFor:      time.Millisecond * 100,
		},
		{
			shouldTimeout: true,
			sleepFor:      time.Second * 2,
		},
	}
	for _, c := range cases {
		_, intStream := HeartbeatGenerateIntStream(done, c.sleepFor, ints...)
		didTimeOut := false
		// notice how we completely ignore the heartbeat channel
		for i, expected := range ints {
			select {
			case v := <-intStream:
				if c.shouldTimeout {
					t.Fatal("we shouldn't be here as this test is supposed to fail")
				}
				if v != expected {
					t.Errorf("index: %v, expected: %v, got : %v", i, expected, v)
				}

			case <-time.After(testWaitingPeriod):
				// here we time out after what we think is a reasonable time for the test
				// to run, which is not a good way to write tests and is an awful position
				// because this is non-deterministic
				didTimeOut = true
				fmt.Println("we timed out as expected but we should use heartbeats instead")
				return
			}
		}
		assert.Equal(t, c.shouldTimeout, didTimeOut)
	}
}

// TestBasicHeartbeatGenerateIntStream_BadTest is a good test because it checks whether a
// heartbeat us
// and is flaky and to prove this we provide a value to sleep on, and we wait only a second.
// But in reality, we don't know what the sleep duration of the called routine would be
func TestBasicHeartbeatGenerateIntStream(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	ints := []int{0, 1, 2, 3, 5}

	pulses, intStream := HeartbeatGenerateIntStream(done, time.Second, ints...)
	<-pulses
	// no timeouts required as here we wait for the goroutine to signal that it is
	// beginning to process an iteration & now we can safely write test without timeouts.
	// The only risk is that if one of the iterations take extremely long, then our test
	// cannot detect it. It can only detect whether the goroutine has started processing.

	i := 0
	for v := range intStream {
		if v != ints[i] {
			t.Errorf("index: %v, expected: %v, got : %v", i, ints[i], v)
		}
		i++
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m,
		goleak.IgnoreTopFunction("time.Sleep"), // flaky tests results in a goroutine that is in deep sleep
	)
}
