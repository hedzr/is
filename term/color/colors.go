// Package color provides a wrapped standard output device like printf but with colored enhancements.
//
// The main types are [Cursor] and [Translator].
//
// [Cursor] allows formatting colorful text and moving cursor to another coordinate.
//
// [New] will return a [Cursor] object.
//
// [RowsBlock] is another cursor controller, which can treat the current line and following lines as a block and updating these lines repeatedly. This feature will help the progressbar writers or the continuous lines updater.
//
// [Translator] is a text and tiny HTML tags translator to convert these markup text into colorful console text sequences.
// [GetCPT] can return a smart translator which translate colorful text or strip the ansi escaped sequence from result text if `states.Env().IsNoColorMode()` is true.
//
// [Color] is an interface type to represent a terminal color object, which can be serialized to ansi escaped sequence directly by [Color.Color].
//
// To create a [Color] object, there are several ways:
//
//   - by [NewColor16], or use [Color16] constants directly like [FgBlack], [BgGreen], ...
//   - by [NewColor256] to make a 8-bit 256-colors object
//   - by [NewColor16m] to make a true-color object
//   - by [NewControlCode] or [ControlCode] constants
//   - by [NewFeCode] or [FeCode] constants
//   - by [NewSGR] or use [CSIsgr] constants directly like [SGRdim], [SGRstrike], ...
//   - by [NewStyle] to make a compounded object
//   - ...
//
// See also [docsite].
//
// [docsite]: https://docs.hedzr.com/docs/is
package color

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"

	_ "github.com/hedzr/is/states"
)

// Color interface represents an ansi escaped
// sequences code in terminal/tty/console.
//
// The supported object includes 4-bit (16-colors),
// 8-bit (256 colors) and 24-bit (true-colors)
// encoders. See also [NewColor16], [NewColor256]
// and [NewColor16m].
//
// These objects are Color (s):
//
//   - Style (multiple Color items)
//   - ControlCode (eg: ESC, CR, FF, LF, HT, BS, BEL)
//   - FeCode (eg: SS2, SS3, DCS, ...)
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
var _ Color = (*Style)(nil)
var _ Color = (*ControlCode)(nil)
var _ Color = (*FeCode)(nil)

// var _ = states.Env().IsNoColorMode()

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

// NewSGR return the given [CSIsgr] code directly.
//
// A [CSIsgr] code is a `ESC n m` sequence to
// represent bold, dim, fg color, bg colors, ...
func NewSGR(code CSIsgr) CSIsgr {
	return code
}

func CSIAddCode(code CSIsuffix) (c CSICodes) {
	c.Items = append(c.Items, code.Code())
	return
}

func CSIAddCode1(code CSIsuffix, n int) (c CSICodes) {
	c.Items = append(c.Items, code.Code1(n))
	return
}

func CSIAddCode2(code CSIsuffix, n, m int) (c CSICodes) {
	c.Items = append(c.Items, code.Code2(n, m))
	return
}

type Color16 int // ANSI Escaped Sequences here

func (c Color16) String() string { return c.Color() }

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

func (c Color16) Int() int {
	return int(c)
}

type Color256 struct {
	// r, g, b, a byte
	clr [4]byte
	bg  bool
}

func (c Color256) String() string { return c.Color() }

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

type Color16m struct {
	// r, g, b, a byte
	clr [4]byte
	bg  bool
}

func (c Color16m) String() string { return c.Color() }

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

// CSICodes wraps several csiCode together.
//
// The best way for creating a CSI sequence is using
// [CSIsuffix.Code], [CSIsuffix.Code1] and
// [CSIsuffix.Code2].
//
// Or you can use thees function: [CSIAddCode],
// [CSIAddCode1] and [CSIAddCode2].
//
// For example,
//
//	c := CSICursorUp.Code1(7)
//	fmt.Printf("%s", c) // move cursor up 7 lines
//	c := CSICursorPosition.Code2(2, 3) // move cursor to row 2 col 3
//	fmt.Printf("%s", c)
//
//	c := color.CSICursorUp.Code1(7)
//	fmt.Printf("%s", c) // move cursor up 7 lines
//	d := color.CSICursorDown.Code1(7)
//	fmt.Printf("%s", d) // move cursor down 7 lines
//
//	fmt.Printf("%s", color.New().SavePos().Build())
//
//	var cx color.CSICodes
//	cx.AddCode2(color.CSICursorPosition, 2, 3) // move cursor to row 2 col 3
//	fmt.Printf("%s", cx)
//
//	fmt.Printf("%s", color.New().RestorePos().Build())
type CSICodes struct {
	Items []csiCode
}

func (c *CSICodes) AddCode(code CSIsuffix) {
	c.Items = append(c.Items, code.Code())
}

func (c *CSICodes) AddCode1(code CSIsuffix, n int) {
	c.Items = append(c.Items, code.Code1(n))
}

func (c *CSICodes) AddCode2(code CSIsuffix, n, m int) {
	c.Items = append(c.Items, code.Code2(n, m))
}

func (c CSICodes) core(out CWriter) {
	for i, it := range c.Items {
		if i > 0 {
			_ = out.WriteByte(';')
		}
		it.core(out)
	}
}

func (c CSICodes) Color() string {
	var sb = NewFmtBuf()
	if len(c.Items) > 0 {
		cc := c.Items[0]
		cc.prologue(sb)
		cc.core(sb)
		cc.epilogue(sb)
	}
	return sb.PutBack()
}

func (c CSICodes) ColorTo(out io.Writer) {
	wrString(out, c.Color())
}

func (c CSICodes) Int() (color int) {
	return
}

func (c CSICodes) String() string { return c.Color() }

const (
	CSICursorUp       CSIsuffix = 'A' // ESC n A   - CUU - Cursor Up. move the cursor n (default 1) cells in the given direction
	CSICursorDown     CSIsuffix = 'B' // ESC n B   - CUD - Cursor Down. move the cursor n (default 1) cells in the given direction
	CSICursorForward  CSIsuffix = 'C' // ESC n C   - CUF - Cursor Forward. move the cursor n (default 1) cells in the given direction
	CSICursorBack     CSIsuffix = 'D' // ESC n D   - CUB - Cursor Back. move the cursor n (default 1) cells in the given direction
	CSICursorNextLine CSIsuffix = 'E' // ESC n E   - CNL - Cursor Next Line. move to beginning of the line n lines down
	CSICursorPrevLine CSIsuffix = 'F' // ESC n F   - CPL - Cursor Previous Line. move to beginning of the line n lines up
	CSICursorHorzAbs  CSIsuffix = 'G' // ESC n G   - CHA - Cursor Horizontal Absolute. move to column n
	CSICursorPosition CSIsuffix = 'H' // ESC n;m H - CUP - Cursor Position. move the cursor to row n, column m (1-based)
	CSIEraseInDisplay CSIsuffix = 'J' // ESC n J   - ED  - Erase in Screen. clears part of the screen. if n is 0, clear from cursor to end of scrren. 1 to beginning, 2 for entire screen and move to up-left corner, 3 for all screen and clear them in the scrollback buffer.
	CSIEraseInLine    CSIsuffix = 'K' // ESC n K   - EL  - Erase in Line. clears part of the line. if n is 0, clear from cursor to the end of the line. 1 to begining, 2 for entire line. cursor position does not change.
	CSIScrollUp       CSIsuffix = 'S' // ESC n S   - SU  - Scroll UP. scroll whole page up by n lines.
	CSIScrollDown     CSIsuffix = 'T' // ESC n T   - SD  - Scroll DOWN. scroll whole page down by n lines.

	CSIHorzVertPosition   CSIsuffix = 'f' // ESC n;m f - HVP - Horizontal Vertical Position. Same as CUP, but counts as a format effector function (lick CR or LF) rather than an editor function (like CUD or CNL).
	CSIAuxPortOn          CSIsuffix = '5' // ESC 5i    -     - Enable aux serial port usually for local serial printer
	CSIAuxPortOff         CSIsuffix = '4' // ESC 4i    -     - Disable aux serial port
	CSIDeviceStatusReport CSIsuffix = '6' // ESC 6n    -     - Reports the cursor postion (CPR) by transmitting `ESC [n;mR`

	// use SGRxxx instead CSISGR
	CSISGR CSIsuffix = 'H' // ESC n m   - SGR - Select Grapthic Rendition. Sets colors and style of the characters following this code.
)

// CSICode will be expanded to `CSI n byte` form.
// For example, Cursor Up ('A') expands as `CSI n A`.
type CSICode struct {
	N      int
	Suffix CSIsuffix
}

type csiCode interface {
	prologue(out CWriter)
	core(out CWriter)
	epilogue(out CWriter)
	And(cs CSIsuffix, n ...int) (codes CSICodes)
}

func (c CSICode) prologue(out CWriter) {
	_, _ = out.WriteString(csi)
}

func (c CSICode) epilogue(out CWriter) {
	_ = out.WriteByte(byte(c.Suffix))
}

func (c CSICode) core(out CWriter) {
	if c.N > 1 {
		_, _ = out.WriteInt(c.N)
	}
}

func (c CSICode) Color() string {
	var sb = NewFmtBuf()
	c.prologue(sb)
	c.core(sb)
	c.epilogue(sb)
	return sb.PutBack()
}

func (c CSICode) String() string { return c.Color() }

func (c CSICode) ColorTo(out io.Writer) {
	wrString(out, c.Color())
}

func (c CSICode) Int() (color int) {
	return
}

func (c CSICode) And(cs CSIsuffix, n ...int) (codes CSICodes) {
	codes.Items = append(codes.Items, c)
	switch len(n) {
	case 0:
		codes.Items = append(codes.Items, cs.Code())
	case 1:
		var nn = 1
		for _, ni := range n {
			nn = ni
		}
		codes.Items = append(codes.Items, cs.Code1(nn))
	case 2:
		nn, nm := n[0], n[1]
		codes.Items = append(codes.Items, cs.Code2(nn, nm))
	}
	return
}

type CSICode2 struct {
	M int
	CSICode
}

func (c CSICode2) core(out CWriter) {
	if c.N > 1 {
		_, _ = out.WriteInt(c.N)
	}
	_ = out.WriteByte(',')
	if c.M > 1 {
		_, _ = out.WriteInt(c.M)
	}
}

type CSIsuffix byte

func (c CSIsuffix) Code() CSICode {
	return CSICode{1, c}
}

func (c CSIsuffix) Code1(n int) CSICode {
	return CSICode{n, c}
}

func (c CSIsuffix) Code2(n, m int) CSICode2 {
	return CSICode2{m, CSICode{n, c}}
}

func (c CSIsuffix) Color() string {
	var sb = NewFmtBuf()
	_, _ = sb.WriteString(csi)
	_ = sb.WriteByte(byte(c))
	return sb.PutBack()
}

func (c CSIsuffix) String() string { return c.Color() }

func (c CSIsuffix) ColorTo(out io.Writer) {
	wrString(out, c.Color())
}

func (c CSIsuffix) Int() (color int) {
	color = int(byte(c))
	return
}

type CSIsgr byte

func (c CSIsgr) Color() string {
	var sb = NewFmtBuf()
	_, _ = sb.WriteString(csi)
	_, _ = sb.WriteInt(int(byte(c)))
	_ = sb.WriteByte('m')
	return sb.PutBack()
}

func (c CSIsgr) String() string { return c.Color() }

func (c CSIsgr) ColorTo(out io.Writer) {
	wrString(out, c.Color())
}

func (c CSIsgr) Int() (color int) {
	color = int(byte(c))
	return
}

const (
	SGRreset                        CSIsgr = iota // reset or normal
	SGRbold                                       // bold or increased intensity
	SGRdim                                        // faint, decreased intensity, or dim
	SGRitalic                                     // italic
	SGRunderline                                  // underline
	SGRslowblink                                  // blink
	SGRrapidblink                                 // fast blink
	SGRinverse                                    // reverse video or invert
	SGRhide                                       // conceal or hide
	SGRstrike                                     // crossed-out or strike
	SGRprimaryfont                                //
	SGRalternativefont1                           //
	SGRalternativefont2                           //
	SGRalternativefont3                           //
	SGRalternativefont4                           //
	SGRalternativefont5                           //
	SGRalternativefont6                           //
	SGRalternativefont7                           //
	SGRalternativefont8                           //
	SGRalternativefont9                           //
	SGRgothic                                     // Fraktur (Gothic), rarely supported
	SGRdoublyUnderline                            // doubly underlined, or not bold
	SGRresetBoldAndDim                            // neither bold nor faint
	SGRresetItalic                                // reset italic
	SGRresetUnderline                             // reset singly or doubly underlined
	SGRresetSlowBlink                             // turn blink off
	SGRresetRapidBlink                            // proportional spacing
	SGRresetInverse                               // reset inverse
	SGRresetHide                                  // not concealed
	SGRresetStrike                                // not cross-out
	SGRfgBlack                                    //
	SGRfgRed                                      //
	SGRfgGreen                                    //
	SGRfgYellow                                   //
	SGRfgBlue                                     //
	SGRfgMagenta                                  //
	SGRfgCyan                                     //
	SGRfgLightGray                                //
	SGRsetFg                                      // use [NewColor256] or [NewColor16m]. 8-bit color; next arguments are `5;n` or `2;r;g;b`
	SGRdefaultFg                                  // reset fg set by [SGRsetFg]
	SGRbgBlack                                    //
	SGRbgRed                                      //
	SGRbgGreen                                    //
	SGRbgYellow                                   //
	SGRbgBlue                                     //
	SGRbgMagenta                                  //
	SGRbgCyan                                     //
	SGRbgLightGray                                //
	SGRsetBg                                      // use [NewColor256] or [NewColor16m]. 8-bit color; next arguments are `5;n` or `2;r;g;b`.
	SGRdefaultBg                                  // reset bg set by [SGRsetBg]
	SGRdisableProportionalSpacing                 //
	SGRframed                                     //
	SGRencircled                                  //
	SGRoverlined                                  // not supported in Terminal.app
	SGRneitherFramedNorEncircled                  //
	SGRnotoverlined                               //
	SGRreserved56                                 //
	SGRreserved57                                 //
	SGRsetUnderlineColor                          // not om standard; implemented in Kitty, VTE, mintty, and iTerm2. Next arguments are `5;n` or `2;r;g;b`.
	SGRdefaultUnderlineColor                      //
	SGRideogramUnderline                          // Ideogram underline or right side line
	SGRideogramDoubleUnderline                    //
	SGRideogramOverline                           //
	SGRideogramDoubleOverline                     //
	SGRideogramStressMarking                      //
	SGRresetIdeogram                              //
	SGRreserved66                                 //
	SGRreserved67                                 //
	SGRreserved68                                 //
	SGRreserved69                                 //
	SGRreserved70                                 //
	SGRreserved71                                 //
	SGRreserved72                                 //
	SGRsuperscript                                //
	SGRsubscript                                  //
	SGRresetSuperscriptAndSubscript               //
	SGRreserved76                                 //
	SGRreserved77                                 //
	SGRreserved78                                 //
	SGRreserved79                                 //
	SGRreserved80                                 //
	SGRreserved81                                 //
	SGRreserved82                                 //
	SGRreserved83                                 //
	SGRreserved84                                 //
	SGRreserved85                                 //
	SGRreserved86                                 //
	SGRreserved87                                 //
	SGRreserved88                                 //
	SGRreserved89                                 //
	SGRfgDarkGray                                 // high density colors, fg, light black
	SGRfgLightRed                                 //
	SGRfgLightGreen                               //
	SGRfgLightYellow                              //
	SGRfgLightBlue                                //
	SGRfgLightMagenta                             //
	SGRfgLightCyan                                //
	SGRfgWhite                                    // Light LightGray
	SGRreserved98                                 //
	SGRreserved99                                 //
	SGRbgDarkGray                                 // high density colors, bg, light black
	SGRbgLightRed                                 //
	SGRbgLightGreen                               //
	SGRbgLightYellow                              //
	SGRbgLightBlue                                //
	SGRbgLightMagenta                             //
	SGRbgLightCyan                                //
	SGRbgWhite                                    // Light LightGray
)

//
// ------ -------
//

func wrPrintf(out io.Writer, format string, args ...any) {
	var data = make([]byte, 0, 32)
	data = fmt.Appendf(data, format, args...)
	_, _ = out.Write(data)
	// buf := NewFmtBuf()
	// buf.Printf(format, args...)
	// _, _ = out.Write([]byte(buf.PutBack()))
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
