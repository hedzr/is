package color

import (
	"strconv"
	"sync"
	"unicode/utf8"
)

var _ CWriter = (*fmtbufS)(nil)
var _ FmtBuf = (*fmtbufS)(nil)

type FmtBuf interface {
	CWriter
	PutBack() (str string)
}

type CWriter interface {
	io.Writer
	WriteString(str string) (n int, err error)
	WriteInt(int) (n int, err error)
	WriteRune(rune) (n int, err error)
	WriteByte(byte) (err error)
}

var colorFormatBufPool = sync.Pool{
	New: func() any {
		return newColorFormatBufPool().reset()
	},
}

func newColorFormatBufPool() *fmtbufS {
	return &fmtbufS{
		buffer: make([]byte, 32),
	}
}

type fmtbufS struct {
	buffer []byte
}

func NewFmtBuf() (s *fmtbufS) {
	return colorFormatBufPool.Get().(*fmtbufS)
}

func (s *fmtbufS) PutBack() (str string) {
	str = string(s.buffer)
	colorFormatBufPool.Put(s.reset())
	return
}

func (s *fmtbufS) reset() *fmtbufS {
	s.buffer = s.buffer[0:0]
	return s
}

func (s *fmtbufS) Write(data []byte) (n int, err error) {
	n = len(data)
	if n > 0 {
		s.buffer = append(s.buffer, data...)
	}
	return
}

func (s *fmtbufS) WriteInt(i int) (n int, err error) {
	n1 := len(s.buffer)
	s.buffer = strconv.AppendInt(s.buffer, int64(i), 10)
	n = len(s.buffer) - n1
	return
}

func (s *fmtbufS) WriteString(str string) (n int, err error) {
	data := []byte(str)
	return s.Write(data)
}

func (s *fmtbufS) WriteRune(r rune) (n int, err error) {
	n1 := len(s.buffer)
	s.buffer = utf8.AppendRune(s.buffer, r)
	return len(s.buffer) - n1, nil
}

func (s *fmtbufS) WriteByte(b byte) (err error) {
	s.buffer = append(s.buffer, b)
	return nil
}
