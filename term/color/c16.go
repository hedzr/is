package color

import "sync/atomic"

type c16S struct {
	clr Color16
	// written int32
	cS *Cursor
}

func (s c16S) writecsicode() {
	_, _ = s.cS.sb.WriteString(s.clr.Color())
	// // if atomic.CompareAndSwapInt32(&s.written, 0, 1) {
	// _, _ = s.cS.sb.WriteString(csi)
	// if int(s.clr) > 0 {
	// 	_, _ = s.cS.sb.WriteString(strconv.Itoa(int(s.clr)))
	// 	_, _ = s.cS.sb.WriteRune('m')
	// }
	// // }
}

// Echo prints contents into buffer for [Cursor.Build].
func (s c16S) Echo(args ...string) *Cursor {
	s.cS.Echo(args...)
	return s.cS
}

// Print prints contents into buffer for [Cursor.Build].
func (s c16S) Print(args ...any) *Cursor {
	s.cS.Print(args...)
	return s.cS
}

// Println prints contents into buffer for [Cursor.Build].
func (s c16S) Println(args ...any) *Cursor {
	s.writecsicode()
	s.cS.Println(args...)
	return s.cS
}

// Printf prints contents into buffer for [Cursor.Build].
func (s c16S) Printf(format string, args ...any) *Cursor {
	s.writecsicode()

	// if s.close {
	// 	defer s.cS.echoResetColor()
	// } else {
	// 	s.cS.closers = append(s.cS.closers, s.echoResetColor)
	// }
	return s.cS.Printf(format, args...)
}

func (s c16S) ResetColor() *Cursor {
	s.cS.ResetColor()
	return s.cS
}

//
//
//

func (s *Cursor) Color(clr Color, format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) Bg(bg Color, format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	s.bg(bg, format, args...)
	return s
}

func (s *Cursor) Effect(bg Color, format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	s.bg(bg, format, args...)
	return s
}

// Color16 make a csi 16-colors sequences.
//
// For example,
//
//	var c = color.New()
//
//	// don't close it, but Println() will close it automatically
//	c.Color16(color.FgRed, false).Printf("hello, %s", "world\n").Println()
//	// don;t close, but String() will close it automatically
//	c.Color16(color.FgGreen, false).Printf("hello, %s", "world\n")
//
//	t.Logf("%s", c.String())
func (s *Cursor) Color16(clr Color16) c16S {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	return c16S{clr, s}
}

// More Fg ...

func (s *Cursor) Black(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgBlack
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) Red(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgRed
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) Green(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgGreen
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) Yellow(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgYellow
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) Blue(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgBlue
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) Magenta(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgMagenta
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) Cyan(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgCyan
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) LightGray(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgLightGray
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) LightBlack(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgLightBlack
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) LightRed(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgLightRed
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) LightGreen(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgLightGreen
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) LightYellow(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgLightYellow
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) LightBlue(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgLightBlue
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) LightMagenta(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgLightMagenta
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) LightCyan(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgLightCyan
	s.fg(clr, format, args...)
	return s
}

func (s *Cursor) White(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = FgWhite
	s.fg(clr, format, args...)
	return s
}

// More Bg ...

func (s *Cursor) BgBlack(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgBlack
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgRed(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgRed
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgGreen(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgGreen
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgYellow(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgYellow
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgBlue(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgBlue
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgMagenta(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgMagenta
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgCyan(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgCyan
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgLightGray(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgLightGray
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgDarkGray(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgDarkGray
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgLightRed(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgLightRed
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgLightGreen(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgLightGreen
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgLightYellow(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgLightYellow
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgLightBlue(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgLightBlue
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgLightMagenta(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgLightMagenta
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgLightCyan(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgLightCyan
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) BgWhite(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgWhite
	s.bg(clr, format, args...)
	return s
}

// More Effects ...

func (s *Cursor) ENormal(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgNormal
	s.bg(clr, format, args...)
	return s
}
func (s *Cursor) EBold(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgBoldOrBright
	s.bg(clr, format, args...)
	return s
}
func (s *Cursor) EHighlight(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgBoldOrBright
	s.bg(clr, format, args...)
	return s
}
func (s *Cursor) EDim(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgDim
	s.bg(clr, format, args...)
	return s
}
func (s *Cursor) EItalic(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgItalic
	s.bg(clr, format, args...)
	return s
}
func (s *Cursor) EUnderline(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgUnderline
	s.bg(clr, format, args...)
	return s
}
func (s *Cursor) EBlink(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgBlink
	s.bg(clr, format, args...)
	return s
}
func (s *Cursor) ERapidBlink(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgRapidBlink
	s.bg(clr, format, args...)
	return s
}
func (s *Cursor) EInverse(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgInverse
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) EHidden(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgHidden
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) EStrikeout(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgStrikeout
	s.bg(clr, format, args...)
	return s
}

// Reset Effects ...

func (s *Cursor) EResetBold(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgResetBoldOrDoubleUnderLine
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) EResetDim(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgResetNormalColorAndBright
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) EResetItalic(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgResetItalic
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) EResetUnderline(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgResetUnderline
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) EResetBlink(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgResetBlink
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) EResetInverse(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgResetInverse
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) EResetHidden(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgResetHidden
	s.bg(clr, format, args...)
	return s
}

func (s *Cursor) EResetStrikeout(format string, args ...any) *Cursor {
	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	const clr = BgResetStrikeout
	s.bg(clr, format, args...)
	return s
}
