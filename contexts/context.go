package contexts

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Zen: Context provides an API for cancelling the branches of a function's call graph
// It also provides a data bag for transporting request scoped data through the graph
// For cancellation of a goroutine, it's parent may want to cancel it or it may want to
// cancel all of its children. Any blocking op within the goroutine needs to be pre-emptable
// so that it may be cancelled. Context helps manage all of these

// prints a greeting and farewell. To print, we have generators for each. To generate, we
// make a call to locale(). Let's say, genGreeting only wants to wait 1 second for locale()
// and if printGreeting() isn't successful, we discard the whole print farewell thing as well

type result struct {
	err error
	msg string
}

// use context values only for request scoped data that transits process and API
// boundaries, not for passing optional parameters to functions
// A: data should not be generated in process memory
// B: data should be immutable
// C: data should trend towards simple types
// D: data should be data and not types with methods
// E: data should decorate functions and not drives them
// this struct satisfies A, B, D but not C, E and is used for demonstration
type contextConfiguration struct {
	greetDelay        time.Duration
	localeComputeTime time.Duration
}

// good practice defining a key type for your context baggage and then explicit methods
// to get keys out of that baggage
type ctxKey string

const ctxConfigKey = "config"

// type safe way to get greet delay duration from the context baggage
func greetDelay(ctx context.Context) time.Duration {
	return ctx.Value(ctxConfigKey).(contextConfiguration).greetDelay
}

// type safe way to get locale cost duration from the context baggage
func localeComputeTime(ctx context.Context) time.Duration {
	return ctx.Value(ctxConfigKey).(contextConfiguration).localeComputeTime
}

func printGreetAndFarewellWith1SecGreetDeadlineCancelsFarewell(config contextConfiguration) (*result, *result) {
	ctx, cancel := context.WithCancel(context.Background())
	// the key you use in the context bag must satisfy comparability and must be safe to
	// access from multiple routines
	ctx = context.WithValue(ctx, ctxConfigKey, config)
	var wg sync.WaitGroup
	farewellResult := &result{}
	greetResult := &result{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if msg, err := genGreeting(ctx); err != nil {
			fmt.Println("cannot print greeting. Cancelling context. Err:", err)
			greetResult.msg = "cannot print greeting"
			greetResult.err = err
			// call cancel() so that any part of the graph depending on this ctx is cancelled
			cancel()
		} else {
			greetResult.msg = msg
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if msg, err := genFarewell(ctx); err != nil {
			fmt.Println("cannot print farewell. Cancelling context. Err:", err)
			farewellResult.msg = "cannot print farewell"
			farewellResult.err = err
			// no need to call cancel here
		} else {
			farewellResult.msg = msg
		}
	}()
	wg.Wait()
	return greetResult, farewellResult
}

func genGreeting(ctx context.Context) (string, error) {
	// we want the call to locale to not take more than a second
	ctx, cancel := context.WithTimeout(ctx, greetDelay(ctx))
	defer cancel()

	// we are passing caller just for debug, don't do it on production
	switch locale, err := locale(ctx, "greeting"); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello world", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(ctx context.Context) (string, error) {
	// no specific bounds in this method on the context apart from what is passed
	switch locale, err := locale(ctx, "farewell"); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye world", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

// locale is a function that can take an arbitrary amount of time to run
// based on the value set in its context
func locale(ctx context.Context, caller string) (string, error) {
	// as an optimization we can early check whether we are going to meet the deadline
	// received in the context vs what the function is set to take (if it is set)
	localeCostTime := localeComputeTime(ctx)
	if deadline, ok := ctx.Deadline(); ok {
		// there is a deadline set here
		if deadline.Sub(time.Now().Add(localeCostTime)) <= 0 {
			fmt.Println("There's no point in continuing locale as we'd exceed the deadline set by", caller)
			return "", context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		fmt.Println("context is closed inside locale() called by", caller)
		return "", ctx.Err()
	case <-time.After(localeCostTime):
	}
	return "EN/US", nil
}
