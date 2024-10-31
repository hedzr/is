package states

import (
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func isDebuggerAttached() bool {
	ppid := os.Getppid()
	if runtime.GOOS == "linux" {
		f, err := os.Open("/proc/" + strconv.Itoa(ppid) + "/cmdline")
		if err != nil {
			return false
		}
		defer f.Close()

		var buf bytes.Buffer
		_, err = buf.ReadFrom(f)
		txt := buf.String()

		re := regexp.MustCompile(`\000`)
		a := re.Split(txt, -1)
		if strings.HasSuffix(a[0], "/dlv") {
			return true
		}
	}
	return false
}
