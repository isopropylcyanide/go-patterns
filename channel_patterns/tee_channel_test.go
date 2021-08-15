package channel_patterns

import (
	"fmt"
	g "patterns/generators"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTee(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	// we create the input using the reusable generators
	input := g.Take(done, g.Repeat(done, "A", "B", "C"), 8)
	teeA, teeB := Tee(done, input)

	// now we loop over both channels (sort of zip) to see both have the same input
	for v1 := range teeA {
		v2 := <-teeB
		assert.Equal(t, v1, v2)
		fmt.Printf("TeeA received [%v] and TeeB received [%v]\n", v1, v2)
	}
}
