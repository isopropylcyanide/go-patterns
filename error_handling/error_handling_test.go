package error_handling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHandlingThatIsNotAbleToPropagateValues(t *testing.T) {
	ErrorHandlingThatIsNotAbleToPropagateValues(getUrls()...)
}

func TestErrorHandlingThatIsAbleToPropagateValues(t *testing.T) {
	err := ErrorHandlingThatIsAbleToPropagateValues(getUrls()...)
	assert.EqualError(t, err, "[Informed] Error processing request https://bad-host: Get \"https://bad-host\": dial tcp: lookup bad-host: no such host")
}

func getUrls() []string {
	return []string{
		"https://www.google.com",
		"https://bad-host",
	}
}
