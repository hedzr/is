//go:build aix || linux || solaris || zos

package term

import (
	"fmt"
	"log/slog"
	"os"
	"syscall"
	"unsafe"
)

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
	cols, rows, _ = getTtySize("")
	return
}

func getTtySize(device string) (cols, rows int, err error) {
	cols, rows, err = getFdSize(uintptr(syscall.Stdin))
	if err != nil {
		slog.Error("[GetTtySize] cannot get terminal size", "err", err)
	}
	return
}

func getDeviceSize(outf *os.File) (cols, rows int, err error) {
	out := outf.Fd()
	return getFdSize(out)
}

func getFdSize(fd uintptr) (cols, rows int, err error) {
	var sz struct {
		rows, cols, xPixels, yPixels uint16
	}

	res, _, errno := syscall.Syscall(syscall.SYS_IOCTL, //nolint:dogsled //like it
		fd,
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&sz)))
	if int(res) == -1 {
		err = fmt.Errorf("[GetTtySize] cannot get terminal size. err = %v, res = %v", errno.Error(), res)
	}
	cols, rows = int(sz.cols), int(sz.rows)
	return
}

func isDoubleClickRun() bool { return false }
