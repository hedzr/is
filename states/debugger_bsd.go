//go:build dragonfly || freebsd || netbsd || openbsd
// +build dragonfly freebsd netbsd openbsd

package states

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/hedzr/is/exec"
)

func isDebuggerAttached() bool {
	ppid := os.Getppid()
	if runtime.GOOS == "freebsd" {
		// var f *os.File
		rc, txt, err := exec.RunWithOutput("ps", "-ef", strconv.Itoa(ppid))
		if err == nil && rc == 0 {
			const title = `  UID   PID  PPID   C STIME   TTY           TIME ` // + `CMD`
			lines := strings.Split(txt, "\n")
			cmdline := lines[1][len(title):]
			parts := strings.Split(cmdline, " ")
			if strings.HasSuffix(parts[0], "/dlv") || strings.HasSuffix(parts[0], "debugserver") {
				return true
			}
		}
	}
	return false
}
