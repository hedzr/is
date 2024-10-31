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
	if runtime.GOOS == "darwin" {
		// var f *os.File
		rc, txt, err := exec.RunWithOutput("ps", "-ef", strconv.Itoa(ppid))
		if err == nil && rc == 0 {
			const title = `  UID   PID  PPID   C STIME   TTY           TIME ` // + `CMD`
			lines := strings.Split(txt, "\n")
			cmdline := lines[1][len(title):]
			if strings.HasSuffix(cmdline, "/dlv") {
				return true
			}
		}
	}
	return false
}
