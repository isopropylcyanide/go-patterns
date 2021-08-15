package for_select_loop

import (
	"testing"

	"go.uber.org/goleak"
)

func TestSendIterationValuesOnChannel(t *testing.T) {
	SendIterationValuesOnChannel()
}

func TestInfiniteLooping(t *testing.T) {
	t.Skipf("Skipping test as this will loop infinitely ")
	InfiniteLooping()
}

func TestInfiniteLoopingII(t *testing.T) {
	t.Skipf("Skipping test as this will loop infinitely ")
	InfiniteLoopingII()
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
