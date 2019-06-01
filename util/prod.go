package util

import "runtime"

func IsRunningInProd() bool {
	return runtime.GOOS == "linux"
}
