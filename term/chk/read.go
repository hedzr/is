package chk

import (
	"io"
	"runtime"
)

// ReadTill reads a line or a string ending with delim of input
// from a terminal without local echo.
func ReadTill(fd int, delim byte) (line string, delimMatched bool, err error) {
	var nig []byte
	nig, delimMatched, err = readBytesTill(fd, delim)
	line = string(nig)
	return
}

func readNoEchoTill(reader io.Reader, delim byte) (ret []byte, delimMatched bool, err error) {
	var buf [1]byte
	var n int
	for {
		n, err = reader.Read(buf[:])
		if n > 0 {
			switch buf[0] {
			case delim:
				delimMatched = true
				return
			case '\b':
				if len(ret) > 0 {
					ret = ret[:len(ret)-1]
				}
			case '\n':
				if runtime.GOOS != "windows" {
					return ret, false, nil
				}
				// otherwise ignore \n
			case '\r':
				if runtime.GOOS == "windows" {
					return ret, false, nil
				}
				// otherwise ignore \r
			default:
				ret = append(ret, buf[0])
			}
			continue
		}
		if err != nil {
			if err == io.EOF && len(ret) > 0 {
				return ret, false, nil
			}
			return ret, false, err
		}
	}
}
