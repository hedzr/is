//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package chk

import "golang.org/x/sys/unix"

func readBytesTill(fd int, delim byte) ([]byte, bool, error) {
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		// return nil, false, err
		return readNoEchoTill(disrectReader(fd), delim)
	}

	newState := *termios
	newState.Lflag &^= (unix.ECHO | unix.ICANON)
	newState.Lflag |= unix.ISIG
	// newState.Iflag |= unix.ICRNL
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil {
		return nil, false, err
	}

	defer unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)

	return readNoEchoTill(disrectReader(fd), delim)
}

// disrectReader is an io.Reader that reads from a specific file descriptor.
type disrectReader int

func (r disrectReader) Read(buf []byte) (int, error) {
	return unix.Read(int(r), buf)
}
