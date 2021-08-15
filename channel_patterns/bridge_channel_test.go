package channel_patterns

import (
	"fmt"
	g "patterns/generators"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBridge(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	// we create the input using the reusable generators
	inputA := g.Take(done, g.Repeat(done, "A", "B", "C"), 5)
	inputB := g.Take(done, g.Repeat(done, true, false), 4)
	inputC := g.Take(done, g.Repeat(done, 1, 2, 3), 3)

	unAbridged := chanStream(inputA, inputB, inputC)
	// expecting to see 12 items in the output of the bridge channel (5 + 4 + 3)
	// order should be maintained in the resulting channel (all of A, all of B, ...)
	bridged := Bridge(done, unAbridged)

	count := 0
	for v := range bridged {
		count += 1
		fmt.Printf("Bridged channel received [%v]\n", v)
	}
	assert.Equal(t, 12, count)
}

// chanStream helps generate a channel of channels from multiple individual channels
// note that we couldn't use the Or Channel because of the signature required for Bridge
func chanStream(input ...<-chan interface{}) <-chan <-chan interface{} {
	out := make(chan (<-chan interface{})) // braces are important in nested directional channel
	go func() {
		defer close(out)
		for _, ch := range input {
			out <- ch
		}
	}()

	return out
}
