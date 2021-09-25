package error_propagation

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Zen: It's better to create wrappers over low level / intermediate alerts at module
// boundaries and return an easy-to-digest error at the source. The error ideally should
// have a unique identifier which can then be referenced to find the full error log.
// Since errors in go, by default do not answer this, it's best to create a custom wrapper.

func TestErrorPropagation(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	dir := t.TempDir()

	cases := []struct {
		name            string
		behaviour       func(string) error
		readableMessage string
		binary          string
		noError         bool
	}{
		{
			name:            "high level module does not propagate module error correctly",
			behaviour:       runJobNotClean,
			readableMessage: "", // no readable message from the high level module
			binary:          "bin123",
		},
		{
			name:            "high level module propagates module error correctly",
			behaviour:       runJobClean,
			readableMessage: "could not find binary bin123",
			binary:          "bin123",
		},
		{
			name:      "high level module has no error",
			behaviour: runJobClean,
			binary:    dir, // root file should exist on most systems
			noError:   true,
		},
	}

	for _, tc := range cases {
		err := tc.behaviour(tc.binary)
		if tc.noError {
			assert.NoError(t, err)
		}
		if err != nil {
			msg := "There was an unexpected issue; please report this as a bug"
			if _, ok := err.(HighLevelErr); ok {
				// notice that only if it is a high level error, the message becomes the wrapped
				// if for some reason, the high level module sends unwrapped errors, it becomes a
				// problem as it may or may not be fit for human consumption and delay debugging
				msg = err.Error()
				assert.Equal(t, tc.readableMessage, msg)
			}
			// if it wasn't a high level error, the message remains the blanket one
			handleError(1, err, msg)
		}
	}
}
