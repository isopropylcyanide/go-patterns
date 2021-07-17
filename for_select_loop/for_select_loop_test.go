package for_select_loop

import "testing"

func TestSendIterationValuesOnChannel(t *testing.T) {
	SendIterationValuesOnChannel()
}

func TestInfiniteLooping(t *testing.T) {
	// this will loop infinitely
	InfiniteLooping()
}

func TestInfiniteLoopingII(t *testing.T) {
	// this will loop infinitely as well
	InfiniteLoopingII()
}
