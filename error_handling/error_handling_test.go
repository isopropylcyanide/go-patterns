package error_handling

import (
	"testing"
)

func TestErrorHandlingThatIsNotAbleToPropagateValues(t *testing.T) {
	ErrorHandlingThatIsNotAbleToPropagateValues()
}

func TestErrorHandlingThatIsAbleToPropagateValues(t *testing.T) {
	ErrorHandlingThatIsAbleToPropagateValues()
}
