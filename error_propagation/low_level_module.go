package error_propagation

import "os"

type LowLevelErr struct {
	error
}

func isGloballyExecutable(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		// here, the low level module returns a correct error at the boundary
		// however, the same cannot be said for the systems that call it
		return false, LowLevelErr{error: wrapError(err, err.Error())}
	}
	return stat.Mode().Perm()&0o100 == 0o100, nil
}
