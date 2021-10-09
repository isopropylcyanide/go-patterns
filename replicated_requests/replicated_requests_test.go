package replicated_requests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/goleak"
)

func TestDoWork(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	result := make(chan int)
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go DoWork(ctx, i, &wg, result)
	}

	r := <-result // select the first result
	fmt.Printf("first completed by handler: %v\n", r)
	assert.True(t, r < 10)

	// no need for the remaining handlers, cancel them
	cancelFunc()
	wg.Wait()
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
