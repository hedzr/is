package chk

import (
	"fmt"
	"runtime"
)

func readBytesTill(fd int, delim byte) ([]byte, bool, error) {
	_, _ = fd, delim
	return nil, false, fmt.Errorf("terminal: ReadTill not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}
