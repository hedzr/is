package chk

import (
	"os"

	"golang.org/x/sys/windows"
)

func readBytesTill(fd int, delim byte) ([]byte, bool, error) {
	var st uint32
	if err := windows.GetConsoleMode(windows.Handle(fd), &st); err != nil {
		return nil, false, err
	}
	old := st

	st &^= (windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT)
	st |= (windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_PROCESSED_INPUT)
	if err := windows.SetConsoleMode(windows.Handle(fd), st); err != nil {
		return nil, false, err
	}

	defer windows.SetConsoleMode(windows.Handle(fd), old)

	var h windows.Handle
	p, _ := windows.GetCurrentProcess()
	if err := windows.DuplicateHandle(p, windows.Handle(fd), p, &h, 0, false, windows.DUPLICATE_SAME_ACCESS); err != nil {
		return nil, false, err
	}

	f := os.NewFile(uintptr(h), "stdin")
	defer f.Close()
	return readNoEchoTill(f, delim)
}
