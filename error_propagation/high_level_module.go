package error_propagation

import (
	"fmt"
	"log"
)

type HighLevelErr struct {
	error
}

func runJobNotClean(id string) error {
	executable, err := isGloballyExecutable(id)
	if err != nil {
		// not clean as we do not wrap the low level error
		return err
	}
	if !executable {
		return wrapError(nil, "binary not executable")
	}
	return nil
}

func runJobClean(id string) error {
	executable, err := isGloballyExecutable(id)
	if err != nil {
		// clean as we  wrap the low level error here
		return HighLevelErr{wrapError(err, "could not find binary %v", id)}
	}
	if !executable {
		return wrapError(nil, "binary not executable")
	}
	return nil
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logId: %v", key))
	// notice how complete errors are logged
	log.Printf("%#v", err)
	// notice how only the user-friendly message is written to out
	fmt.Printf("[%v] %v\n", key, message)
}
