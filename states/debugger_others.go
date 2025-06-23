//go:build !dragonfly && !freebsd && !netbsd && !openbsd && !darwin && !linux && !windows
// +build !dragonfly,!freebsd,!netbsd,!openbsd,!darwin,!linux,!windows

package states

func isDebuggerAttached() bool {
	return false
}
