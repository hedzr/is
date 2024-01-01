//go:build !aix && !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd && !zos && !windows && !solaris && !plan9

package term

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
	return
}

func isDoubleClickRun() bool { return false }
