package contexts

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrintGreetAndFarewellWithGreetDeadlineCancelsFarewell(t *testing.T) {
	greetResult, farewellResult := printGreetAndFarewellWith1SecGreetDeadlineCancelsFarewell(
		contextConfiguration{
			greetDelay:        time.Second,
			localeComputeTime: time.Minute,
			// this would ensure that locale would always be delayed past greet
		})
	assert.Equal(t, "cannot print greeting", greetResult.msg)
	assert.Equal(t, context.DeadlineExceeded, greetResult.err)

	// farewell result should indicate that it has been cancelled
	assert.Equal(t, "cannot print farewell", farewellResult.msg)
	assert.Equal(t, context.Canceled, farewellResult.err)
}

func TestPrintGreetAndFarewellWithNoGreetDeadlineRunsNormally(t *testing.T) {
	greetResult, farewellResult := printGreetAndFarewellWith1SecGreetDeadlineCancelsFarewell(
		contextConfiguration{
			greetDelay:        time.Second,
			localeComputeTime: time.Millisecond,
			// this would ensure that locale would always get a chance to run
		})
	assert.Equal(t, "hello world", greetResult.msg)
	assert.NoError(t, greetResult.err)

	assert.Equal(t, "goodbye world", farewellResult.msg)
	assert.NoError(t, farewellResult.err)
}
