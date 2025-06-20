package color

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/hedzr/is/states"
	"github.com/hedzr/is/term"
)

// New returns a *cS (Cursor) object so that you could
// render the colorful text with it.
//
// Each cursor object must have ended by `Build()`.
//
// For example:
//
//	// another colorful builfer
//	c = color.New()
//	fmt.Println(c.Color16(color.FgRed).
//
//	Printf("hello, %s.", "world").Println().Build())
//
// With *cS (Cursor) object, you could make Color16(),
// Color256(), or RGB() text.
//
// See the example codes.
func New() (s *Cursor) {
	s = poolBuilder.Get().(*Cursor)
	if atomic.CompareAndSwapInt32(&s.ncUpdated, 1, 0) {
		s.reinit(false)
	}
	return s
}

var poolBuilder = sync.Pool{
	New: func() any {
		return newBuilder(true)
	},
}

func newBuilder(first bool) (s *Cursor) {
	s = &Cursor{}
	s.reinit(first)
	return s
}

type Cursor struct {
	colorful  bool
	useColor  bool
	ncUpdated int32
	needReset int32
	sb        bytes.Buffer
	w         Writer
	sw        io.StringWriter

	// c16S
	// c88S
	// c256S
	// c16MS

	closers []func()
}

type Result struct {
	bytes.Buffer
}

func (s Result) String() string {
	return s.Buffer.String()
}

func (s *Cursor) reinit(first bool) {
	s.useColor = !states.Env().IsNoColorMode()
	s.w = os.Stdout
	s.sw = os.Stdout
	if first {
		states.Env().SetOnNoColorChanged(s.updateCandidated)
		s.colorful = term.IsColorful(os.Stdout)
	}
}

func (s *Cursor) updateCandidated(mod bool, level int) {
	atomic.CompareAndSwapInt32(&s.ncUpdated, 0, 1)
}

func (s *Cursor) Reset() {
	s.sb.Reset()
	s.useColor = !states.Env().IsNoColorMode()
	// // s.w = &s.sb
	// // s.sw = &s.sb
	// s.w = os.Stdout
	// s.sw = os.Stdout
	// return
}

func (s *Cursor) Build() (r string) {
	for i := len(s.closers) - 1; i >= 0; i-- {
		fn := s.closers[i]
		if fn != nil {
			fn()
		}
	}
	s.closers = s.closers[0:0]

	if atomic.CompareAndSwapInt32(&s.needReset, 1, 0) {
		s.echoResetColor()
	}

	r = s.sb.String()
	s.sb.Reset()
	poolBuilder.Put(s)
	return
}

func (s *Cursor) WithWriter(w io.Writer) *Cursor {
	s.colorful = term.IsColorful(w)
	s.useColor = !states.Env().IsNoColorMode()
	if sw, ok := w.(Writer); ok {
		s.w = sw
	} else {
		panic("want a color.Writer object as input, such as os.Stdout")
	}
	if sw, ok := w.(io.StringWriter); ok {
		s.sw = sw
	}
	return s
}

func (s *Cursor) String() string {
	if atomic.CompareAndSwapInt32(&s.needReset, 1, 0) {
		s.echoResetColor()
	}
	return s.sb.String()
}

func (s *Cursor) ResetColor() *Cursor {
	s.echoResetColor()
	return s
}

// Translate translates color tags from input string and make it colorful.
//
// The tags include:
//
//	<kbd></kbd>, <b></b>, ...
func (s *Cursor) Translate(str string) *Cursor {
	result := GetCPT().TranslateTo(str, FgDefault)
	_, _ = s.sb.WriteString(result)
	return s
}

// StripLeftTabsColorful strips the least left side tab chars from lines.
// It also strips html tags.
// At the end, StripLeftTabsC try translate color code in string.
func (s *Cursor) StripLeftTabsColorful(str string) *Cursor {
	result := cpt.stripLeftTabs(str)
	_, _ = s.sb.WriteString(result)
	return s
}

// StripLeftTabs strips the least left side tab chars from lines.
// It also strips html tags.
func (s *Cursor) StripLeftTabs(str string) *Cursor {
	result := cptNC.stripLeftTabs(str)
	_, _ = s.sb.WriteString(result)
	return s
}

// StripLeftTabsOnly strips the least left side tab chars from lines.
func (s *Cursor) StripLeftTabsOnly(str string) *Cursor {
	result := cptNC.stripLeftTabsOnly(str)
	_, _ = s.sb.WriteString(result)
	return s
}

// StripHTMLTags aggressively strips HTML tags from a string.
// It will only keep anything between `>` and `<`.
func (s *Cursor) StripHTMLTags(str string) *Cursor {
	result := cptNC.stripHTMLTags(str)
	_, _ = s.sb.WriteString(result)
	return s
}

func (s *Cursor) Echo(args ...string) *Cursor {
	for _, z := range args {
		_, _ = s.sb.WriteString(z)
	}
	return s
}

func (s *Cursor) Print(args ...any) *Cursor {
	_, _ = fmt.Fprint(&s.sb, args...)
	return s
}

func (s *Cursor) Println(args ...any) *Cursor {
	_, _ = fmt.Fprint(&s.sb, args...)
	if atomic.CompareAndSwapInt32(&s.needReset, 1, 0) {
		s.echoResetColor()
	}
	_, _ = fmt.Fprintln(&s.sb)
	return s
}

func (s *Cursor) Printf(format string, args ...any) *Cursor {
	s.print(format, args...)
	return s
}

func (s *Cursor) print(format string, args ...any) {
	if len(args) > 0 {
		_, _ = fmt.Fprintf(&s.sb, format, args...)
		// fmt.Appendf(s.sb, format, args...)
		// strconv.AppendFormat()
		// s.sb.append
	} else {
		_, _ = s.sb.WriteString(format)
	}
}

func (s *Cursor) fg(clr Color, format string, args ...any) {
	s.echoColor(clr)
	if len(args) > 0 {
		_, _ = fmt.Fprintf(&s.sb, format, args...)
		// fmt.Appendf(s.sb, format, args...)
		// strconv.AppendFormat()
		// s.sb.append
	} else {
		_, _ = s.sb.WriteString(format)
	}
	// atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	// if close {
	// 	// defer s.echoResetColor()
	// 	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	// } else {
	// 	s.closers = append(s.closers, s.echoResetColor)
	// }
}

func (s *Cursor) bg(clr Color, format string, args ...any) {
	s.echoBg(clr)
	if len(args) > 0 {
		_, _ = fmt.Fprintf(&s.sb, format, args...)
	} else {
		_, _ = s.sb.WriteString(format)
	}
	// atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	// if close {
	// 	// defer s.echoResetColor()
	// 	atomic.CompareAndSwapInt32(&s.needReset, 0, 1)
	// } else {
	// 	s.closers = append(s.closers, s.echoResetColor)
	// }
}

func (s *Cursor) echoColor(clr Color) {
	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
	if clr != NoColor {
		_, _ = s.sb.WriteString(clr.Color())
		// _, _ = s.sb.WriteString(csi)
		// _, _ = s.sb.Write([]byte(strconv.Itoa(int(clr))))
		// _, _ = s.sb.WriteRune('m')
	}
}

func (s *Cursor) echoColorAndBg(clr, bg Color) {
	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
	if clr != NoColor {
		_, _ = s.sb.WriteString(clr.Color())
		// _, _ = s.sb.WriteString(csi)
		// _, _ = s.sb.Write([]byte(strconv.Itoa(int(clr))))
		// _, _ = s.sb.WriteRune('m')
	}
	s.echoBg(bg)
}

func (s *Cursor) echoBg(bg Color) {
	if bg != NoColor {
		_, _ = s.sb.WriteString(bg.Color())
		// _, _ = s.sb.WriteString(csi)
		// _, _ = s.sb.Write([]byte(strconv.Itoa(int(bg))))
		// _, _ = s.sb.WriteRune('m')
	}
}

func (s *Cursor) echoResetColor() { //nolint:unused //no
	// _, _ = fmt.Fprint(os.Stdout, "\x1b[0m")
	_, _ = s.sb.WriteString(ResetToNormalColor.Color())
	// _, _ = s.sb.WriteString(csi)
	// _, _ = s.sb.WriteRune('0')
	// _, _ = s.sb.WriteRune('m')
}

type Writer interface {
	io.Writer
	Fd() uintptr
}

const bell = '\x07'           // CTRL-G BEL, Makes an audible noise.
const backspace = '\x08'      // CTRL-H BS, Moves the cursor left (but may "backwards wrap" if cursor is at start of line).
const tabstop = '\x09'        // CTRL-I HT, Moves the cursor right to next tab stop.
const linefeed = '\x0a'       // CTRL-J LF, Moves to next line, scrolls the display up if at bottom of the screen. Usually does not move horizontally, though programs should not rely on this.
const formfeed = '\x0c'       // CTRL-L FF, Move a printer to top of next page. Usually does not move horizontally, though programs should not rely on this. Effect on video terminals varies.
const carriagereturn = '\x0d' // CTRL-M CR, Moves the cursor to column zero.
const escape = '\x1b'         // CTRL-[ ESC, Starts all the escape sequences
const csi = "\x1b["
const ESCAPE = '\x1b'

const (
	Reset Color16 = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

const (
	ResetBold Color16 = iota + 22
	ResetItalic
	ResetUnderline
	ResetBlinking
	_
	ResetReversed
	ResetConcealed
	ResetCrossedOut
)
