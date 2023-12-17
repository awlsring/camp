package sys

import "runtime"

func IsX86() bool {
	return runtime.GOARCH == "amd64"
}

func IsARM() bool {
	return runtime.GOARCH == "arm"
}

func IsARM64() bool {
	return runtime.GOARCH == "arm64"
}

func IsRiscV() bool {
	return runtime.GOARCH == "riscv64"
}
