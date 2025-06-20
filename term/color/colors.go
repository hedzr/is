// Package color provides a wrapped standard output device like printf but with colored enhancements.
package color

import (
	"encoding/binary"
	"io"
	"strconv"
	"unicode/utf8"
)

// Color interface represents an ansi escaped
// sequences code in terminal/tty/console.
//
// The supported object includes 4-bit (16-colors),
// 8-bit (256 colors) and 16-bit (true-colors)
// encoders. See also [NewColor16], [NewColor256]
// and [NewColor16m].
//
// Control codes and more sequences will be
// supported soon.
type Color interface {
	Int() int
	Color() string
	ColorTo(out io.Writer)
}

var _ Color = (*Color16)(nil)
var _ Color = (*Color256)(nil)
var _ Color = (*Color16m)(nil)

// NewColor16m constrcuts a true-color object
// which can be serialized as ansi escaped
// sequences by calling [Color]() string.
func NewColor16m(r, g, b byte, isBg bool) Color16m {
	return Color16m{
		clr: [4]byte{r, g, b}, bg: isBg,
	}
}

// NewColor16m constrcuts a 8-bit object
// which can be serialized as ansi escaped
// sequences by calling [Color]() string.
func NewColor256(clr byte, isBg bool) Color256 {
	return Color256{
		clr: [4]byte{clr}, bg: isBg,
	}
}

// Color256table prints a 8-bit color table for testing
func Color256table(out io.Writer) {
	for row := range 16 {
		for col := range 16 {
			val := row*16 + col
			c := NewColor256(byte(val), true)
			fmt.Fprintf(out, "%s %3d %s", c, val, SGRdefaultBg)
		}
		fmt.Fprintln(out)
	}
}

// NewColor16 cast a clr to Color16.
//
// Valid Color16 color codes include:
// fore- and bg-color (eg [FgRed],
// [BgBlack], ...), and effects (such as bold/hilight -
// [BgBoldOrBright], italic - [BgItalic], dim -
// [BgDim], ...)
func NewColor16(clr Color16) Color16 {
	// _ = isBg //ignore it
	return Color16(clr)
}

// NewStyle creates a container of [Color] objects.
// All of these children will be bound and printed
// in a one sequences.
func NewStyle() Style {
	return Style{}
}

// NewControlCode return the given [ControlCode] code directly.
func NewControlCode(code ControlCode) ControlCode {
	return code
}

// NewFeCode return the given [FeCode] code directly.
func NewFeCode(code FeCode) FeCode {
	return code
}

type Color16 int // ANSI Escaped Sequences here

func (c Color16) Color() string {
	var sb = NewFmtBuf()
	if i := int(c); i >= 0 {
		_, _ = sb.WriteString(csi)
		_, _ = sb.WriteInt(i)
		_, _ = sb.WriteRune('m')
	}
	return sb.PutBack()
}

func (c Color16) ColorTo(out io.Writer) {
	if i := int(c); i >= 0 {
		wrString(out, csi)
		wrInt(out, i)
		wrRune(out, 'm')
	}
}

func wrString(out io.Writer, str string) {
	data := []byte(str)
	_, _ = out.Write(data)
}
func wrInt(out io.Writer, i int) {
	var buffer []byte
	buffer = strconv.AppendInt(buffer, int64(i), 10)
	_, _ = out.Write(buffer)
}
func wrRune(out io.Writer, r rune) {
	// n1 := len(s.buffer)
	var buffer []byte
	buffer = utf8.AppendRune(buffer, r)
	_, _ = out.Write(buffer)
	// return len(s.buffer) - n1, nil
}

func (c Color16) Int() int {
	return int(c)
}

type Color16m struct {
	// r, g, b, a byte
	clr [4]byte
	bg  bool
}

func (c Color16m) Color() string {
	var sb = NewFmtBuf()
	_, _ = sb.WriteString(csi)
	if c.bg {
		_, _ = sb.WriteString("48;2;")
	} else {
		_, _ = sb.WriteString("38;2;")
	}
	_, _ = sb.WriteInt(int(c.clr[0])) // r
	_, _ = sb.WriteRune(';')
	_, _ = sb.WriteInt(int(c.clr[1])) // g
	_, _ = sb.WriteRune(';')
	_, _ = sb.WriteInt(int(c.clr[2])) // b
	_, _ = sb.WriteRune('m')
	return sb.PutBack()
}

func (c Color16m) ColorTo(out io.Writer) {
	wrString(out, csi)
	if c.bg {
		wrString(out, "48;2;")
	} else {
		wrString(out, "38;2;")
	}
	wrInt(out, int(c.clr[0])) // r
	wrRune(out, ';')
	wrInt(out, int(c.clr[1])) // g
	wrRune(out, ';')
	wrInt(out, int(c.clr[2])) // b
	wrRune(out, 'm')
}

func (c Color16m) Int() (color int) {
	_, _ = binary.Decode(c.clr[0:4], binary.LittleEndian, &color)
	return
}

// Style is an array of [Color] objects
type Style struct {
	Items []Color
}

func (c *Style) Add(colors ...Color) *Style {
	c.Items = append(c.Items, colors...)
	return c
}

func (c Style) String() string { return c.Color() }

func (c Style) Color() string {
	var sb = NewFmtBuf()
	for _, it := range c.Items {
		it.ColorTo(sb)
	}
	return sb.PutBack()
}

func (c Style) ColorTo(out io.Writer) {
	for _, it := range c.Items {
		it.ColorTo(out)
	}
}

func (c Style) Int() (color int) {
	for _, it := range c.Items {
		color = it.Int()
	}
	return
}

type ControlCode byte

func (c ControlCode) String() string { return c.Color() }

func (c ControlCode) Color() string {
	return string(byte(c))
}

func (c ControlCode) ColorTo(out io.Writer) {
	wrString(out, c.Color())
}

func (c ControlCode) Int() (color int) {
	color = int(byte(c))
	return
}

// See also these rune(s)
//
// const bell = '\x07'           // CTRL-G BEL, Makes an audible noise.
// const backspace = '\x08'      // CTRL-H BS, Moves the cursor left (but may "backwards wrap" if cursor is at start of line).
// const tabstop = '\x09'        // CTRL-I HT, Moves the cursor right to next tab stop.
// const linefeed = '\x0a'       // CTRL-J LF, Moves to next line, scrolls the display up if at bottom of the screen. Usually does not move horizontally, though programs should not rely on this.
// const formfeed = '\x0c'       // CTRL-L FF, Move a printer to top of next page. Usually does not move horizontally, though programs should not rely on this. Effect on video terminals varies.
// const carriagereturn = '\x0d' // CTRL-M CR, Moves the cursor to column zero.
// const escape = '\x1b'         // CTRL-[ ESC, Starts all the escape sequences
const (
	BEL ControlCode = bell           // CTRL-G BEL, Makes an audible noise.
	BS  ControlCode = backspace      // CTRL-H BS, Moves the cursor left (but may "backwards wrap" if cursor is at start of line).
	HT  ControlCode = tabstop        // CTRL-I HT, Moves the cursor right to next tab stop.
	LF  ControlCode = linefeed       // CTRL-J LF, Moves to next line, scrolls the display up if at bottom of the screen. Usually does not move horizontally, though programs should not rely on this.
	FF  ControlCode = formfeed       // CTRL-L FF, Move a printer to top of next page. Usually does not move horizontally, though programs should not rely on this. Effect on video terminals varies.
	CR  ControlCode = carriagereturn // CTRL-M CR, Moves the cursor to column zero.
	ESC ControlCode = escape         // CTRL-[ ESC, Starts all the escape sequences
)

// FeCode will be expanded as ESC + byte sequence.
// For example, CSI ('\x9B') will be expanded to '\x1B\x9B' (`ESC [`).
type FeCode byte

func (c FeCode) String() string { return c.Color() }

func (c FeCode) Color() string {
	var sb = NewFmtBuf()
	_ = sb.WriteByte(ESCAPE)
	_ = sb.WriteByte(byte(c))
	return sb.PutBack()
}

func (c FeCode) ColorTo(out io.Writer) {
	wrString(out, c.Color())
}

func (c FeCode) Int() (color int) {
	color = int(byte(c))
	return
}

const (
	SS2 FeCode = '\x8E' // ESC N
	SS3 FeCode = '\x8F' // ESC 0
	DCS FeCode = '\x90' // ESC P
	CSI FeCode = '\x9B' // ESC [
	ST  FeCode = '\x9c' // ESC \
	OSC FeCode = '\x9D' // ESC ]
	SOS FeCode = '\x98' // ESC X
	PM  FeCode = '\x9E' // ESC ^
	APC FeCode = '\x9F' // ESC _
)

}

func (c Color256) Color() string {
	var sb = NewFmtBuf()
	_, _ = sb.WriteString(csi)
	if c.bg {
		_, _ = sb.WriteString("48;5;")
	} else {
		_, _ = sb.WriteString("38;5;")
	}
	_, _ = sb.WriteInt(int(c.clr[0])) // r
	_, _ = sb.WriteRune('m')
	return sb.PutBack()
}

func (c Color256) ColorTo(out io.Writer) {
	wrString(out, csi)
	if c.bg {
		wrString(out, "48;5;")
	} else {
		wrString(out, "38;5;")
	}
	wrInt(out, int(c.clr[0])) // n
	wrRune(out, 'm')
}

func (c Color256) Int() (color int) {
	_, _ = binary.Decode(c.clr[0:4], binary.LittleEndian, &color)
	return
}

const (
	// https://en.wikipedia.org/wiki/ANSI_escape_code
	// https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97

	// FgBlack terminal color code
	FgBlack = Color16(30)
	// FgRed terminal color code
	FgRed = Color16(31)
	// FgGreen terminal color code
	FgGreen = Color16(32)
	// FgYellow terminal color code
	FgYellow = Color16(33)
	// FgBlue terminal color code
	FgBlue = Color16(34)
	// FgMagenta terminal color code
	FgMagenta = Color16(35)
	// FgCyan terminal color code
	FgCyan = Color16(36)
	// FgLightGray terminal color code (White)
	FgLightGray = Color16(37)

	// FgDarkGray terminal color code (Gray, Light Black).
	//
	// A highlight/bright black color, maybe 50% gray. See FgLightBlack.
	FgDarkGray = Color16(90)
	// FgLightBlack terminal color code (Gray, Light Black).
	//
	// A highlight/bright black color, maybe 50% gray. See FgDarkGray.
	FgLightBlack = Color16(90)
	// FgLightRed terminal color code
	FgLightRed = Color16(91)
	// FgLightGreen terminal color code
	FgLightGreen = Color16(92)
	// FgLightYellow terminal color code
	FgLightYellow = Color16(93)
	// FgLightBlue terminal color code
	FgLightBlue = Color16(94)
	// FgLightMagenta terminal color code
	FgLightMagenta = Color16(95)
	// FgLightCyan terminal color code
	FgLightCyan = Color16(96)
	// FgWhite terminal color code (Light White)
	FgWhite = Color16(97)

	// BgBlack terminal color code
	BgBlack = Color16(40)
	// BgRed terminal color code
	BgRed = Color16(41)
	// BgGreen terminal color code
	BgGreen = Color16(42)
	// BgYellow terminal color code
	BgYellow = Color16(43)
	// BgBlue terminal color code
	BgBlue = Color16(44)
	// BgMagenta terminal color code
	BgMagenta = Color16(45)
	// BgCyan terminal color code
	BgCyan = Color16(46)
	// BgLightGray terminal color code
	BgLightGray = Color16(47)

	// BgDarkGray terminal color code
	BgDarkGray = Color16(100)
	// BgLightRed terminal color code
	BgLightRed = Color16(101)
	// BgLightGreen terminal color code
	BgLightGreen = Color16(102)
	// BgLightYellow terminal color code
	BgLightYellow = Color16(103)
	// BgLightBlue terminal color code
	BgLightBlue = Color16(104)
	// BgLightMagenta terminal color code
	BgLightMagenta = Color16(105)
	// BgLightCyan terminal color code
	BgLightCyan = Color16(106)
	// BgWhite terminal color code
	BgWhite = Color16(107)

	// BgNormal terminal color code.
	//
	// All attributes become turned off.
	BgNormal = Color16(0)
	// BgBoldOrBright terminal color code
	//
	// Bold or increased intensity
	BgBoldOrBright = Color16(1)
	// BgDim terminal color code.
	//
	// Faint, decreased intensity, or dim.
	// May be implemented as a light font weight like bold.
	BgDim = Color16(2)
	// BgItalic terminal color code.
	//
	// Not widely supported. Sometimes treated as inverse or blink
	BgItalic = Color16(3)
	// BgUnderline terminal color code.
	//
	// Style extensions exist for Kitty, VTE, mintty, iTerm2 and Konsole.
	BgUnderline = Color16(4)
	// BgBlink terminal color code.
	//
	// Slow blink, Sets blinking to less than 150 times per minute.
	// But in many tty it's no effect.
	//
	// Sometimes it can be used for switching to 'normal' bg state without
	// reset all fg and bg settings (if using bg code 0)
	BgBlink = Color16(5)
	// BgRapidBlink terminal color code.
	//
	// MS-DOS ANSI.SYS, 150+ per minute; not widely supported.
	//
	// Sometimes it can be used for switching to 'normal' bg state without
	// reset all fg and bg settings (if using bg code 0)
	BgRapidBlink = Color16(6)
	// BgInverse terminal color code.
	//
	// Swap foreground and background colors; inconsistent emulation
	BgInverse = Color16(7)
	// BgHidden terminal color code.
	//
	// not widely supported.
	BgHidden = Color16(8)
	// BgStrikeout terminal color code.
	//
	// marked as if for deletion.
	BgStrikeout = Color16(9)

	BgResetBoldOrDoubleUnderLine = Color16(21)
	BgResetNormalColorAndBright  = Color16(22) // = BgResetDim
	BgResetItalic                = Color16(23)
	BgResetUnderline             = Color16(24)
	BgResetBlink                 = Color16(25)
	BgResetInverse               = Color16(27)
	BgResetHidden                = Color16(28)
	BgResetStrikeout             = Color16(29)

	FgDarkColor = FgLightGray

	FgDefault = Color16(39)
	BgDefault = Color16(49)

	ResetToNormalColor = Color16(0)

	// NoColor is not a declared ansi code but we can use it for identifying
	// a variable isn't initializing yet.
	NoColor = Color16(-1)
)
