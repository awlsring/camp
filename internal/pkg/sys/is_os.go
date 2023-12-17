package sys

import "runtime"

func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsFreeBSD() bool {
	return runtime.GOOS == "freebsd"
}

func IsOpenBSD() bool {
	return runtime.GOOS == "openbsd"
}
