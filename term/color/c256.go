package color

import (
	"strconv"
	"sync/atomic"
)

type c256S struct {
	n       byte // FgXXX, FgLightXXX, BgXXX, BgLightXXX
	bg      bool
	written int32
	cS      *Cursor
}

func (s c256S) writecsicode() {
	if atomic.CompareAndSwapInt32(&s.written, 0, 1) {
		// color code
		_, _ = s.cS.sb.WriteString(csi)
		if s.bg {
			_, _ = s.cS.sb.WriteString("48;5;")
		} else {
			_, _ = s.cS.sb.WriteString("38;5;")
		}
		if s.n > 0 {
			_, _ = s.cS.sb.WriteString(strconv.Itoa(int(s.n)))
			_, _ = s.cS.sb.WriteRune('m')
		}
	}
}

// Echo prints contents into buffer for [Cursor.Build].
func (s c256S) Echo(args ...string) *Cursor {
	s.cS.Echo(args...)
	return s.cS
}

// Print prints contents into buffer for [Cursor.Build].
func (s c256S) Print(args ...any) *Cursor {
	s.cS.Print(args...)
	return s.cS
}

// Println prints contents into buffer for [Cursor.Build].
func (s c256S) Println(args ...any) *Cursor {
	s.writecsicode()
	s.cS.Println(args...)
	return s.cS
}

// Printf prints contents into buffer for [Cursor.Build].
func (s c256S) Printf(format string, args ...any) *Cursor {
	s.writecsicode()
	return s.cS.Printf(format, args...)
}

func (s c256S) ResetColor() *Cursor {
	s.cS.ResetColor()
	return s.cS
}

// Color256 starts a child builder for 256-colors foreground color.
// The `n` is in 0..255.
func (s *Cursor) Color256(n byte) c256S {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	return c256S{n, false, 0, s}
}

// Bg256 starts a child builder for 256-colors background color.
// The `n` is in 0..255.
func (s *Cursor) Bg256(n byte) c256S {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	return c256S{n, true, 0, s}
}
