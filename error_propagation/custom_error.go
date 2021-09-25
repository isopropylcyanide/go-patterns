package error_propagation

import (
	"fmt"
	"runtime/debug"
)

// CustomError is a wrapper. An ideal system error should let us answer the following
// 	- What happened
// 	- When and where it happened
// 	- Friendly user message
// 	- How the user can get more information
type CustomError struct {
	// Inner captures the underlying error
	Inner error
	// Message is a human-readable message
	Message string
	// Stacktrace when the error was created
	Stacktrace string
	// Misc is a catch-all bag for storing other useful diagnostics for the underlying error
	Misc map[string]interface{}
}

func (err CustomError) Error() string {
	return err.Message
}

func wrapError(err error, message string, msgArgs ...interface{}) CustomError {
	return CustomError{
		Inner:      err,
		Message:    fmt.Sprintf(message, msgArgs...),
		Stacktrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}
