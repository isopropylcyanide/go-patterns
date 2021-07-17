package fan_out_fan_in

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrimeNumberFinderWithNoFanOut(t *testing.T) {
	start := time.Now()
	count := 20000
	out := PrimeNumberFinderWithNoFanOut(count, 130000000, 3)
	assert.Equal(t, count, len(out))
	fmt.Printf("No fanout primes finished in %v\n", time.Since(start))
}

func TestPrimeNumberFinderWithFanOut(t *testing.T) {
	start := time.Now()
	count := 20000
	out := PrimeNumberFinderWithFanOut(count, 130000000, 3)
	assert.Equal(t, count, len(out))
	fmt.Printf("With fanout primes finished in %v\n", time.Since(start))
}
