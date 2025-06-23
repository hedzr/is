package color

import (
	"strconv"
	"sync/atomic"
)

type c16MS struct {
	r, g, b int
	bg      bool
	written int32
	cS      *Cursor
}

func (s c16MS) writecsicode() {
	if atomic.CompareAndSwapInt32(&s.written, 0, 1) {
		_, _ = s.cS.sb.WriteString(csi)
		if s.bg {
			_, _ = s.cS.sb.WriteString("48;2;")
		} else {
			_, _ = s.cS.sb.WriteString("38;2;")
		}

		_, _ = s.cS.sb.WriteString(strconv.Itoa(int(s.r)))
		_, _ = s.cS.sb.WriteRune(';')
		_, _ = s.cS.sb.WriteString(strconv.Itoa(int(s.g)))
		_, _ = s.cS.sb.WriteRune(';')
		_, _ = s.cS.sb.WriteString(strconv.Itoa(int(s.b)))
		_, _ = s.cS.sb.WriteRune('m')
	}
}

// Echo prints contents into buffer for [Cursor.Build].
func (s c16MS) Echo(args ...string) *Cursor {
	s.cS.Echo(args...)
	return s.cS
}

// Print prints contents into buffer for [Cursor.Build].
func (s c16MS) Print(args ...any) *Cursor {
	s.cS.Print(args...)
	return s.cS
}

// Println prints CSI color sequences and prints the formatted text.
func (s c16MS) Println(args ...any) *Cursor {
	s.writecsicode()
	s.cS.Println(args...)
	return s.cS
}

// Printf prints CSI color sequences and prints the formatted text.
func (s c16MS) Printf(format string, args ...any) *Cursor {
	s.writecsicode()
	return s.cS.Printf(format, args...)
}

func (s c16MS) ResetColor() *Cursor {
	s.cS.ResetColor()
	return s.cS
}

// RGB starts a child builder for true-colors foreground color.
// The `r`, `g`, and `b` are a 0..255 number.
func (s *Cursor) RGB(r, g, b int) c16MS {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	return c16MS{r, g, b, false, 0, s}
}

// BgRGB starts a child builder for true-colors background color.
// The `r`, `g`, and `b` are a 0..255 number.
func (s *Cursor) BgRGB(r, g, b int) c16MS {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	return c16MS{r, g, b, true, 0, s}
}
