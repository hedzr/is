package color

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"golang.org/x/net/html"

	"github.com/hedzr/is/states"
	"github.com/hedzr/is/term"
)

// Index is an indexer for retrieving the entries in this subpackage.
type Index struct{}

// GetColorTranslater returns a translator about ANSI Escaped Code.
// It may or may not translating colored text depending
// on cabin.GetNoColorMode.
func (Index) GetColorTranslater() Translator                    { return GetCPT() }
func (Index) GetColorTranslaterAlways() Translator              { return GetCPTC() }
func (Index) GetColorTranslaterNever() Translator               { return GetCPTNC() }
func (Index) GetDummyTranslater() Translator                    { return GetDummyTranslator() }
func (Index) ToColorString(clr Color) string                    { return cpt.toColorString(clr) }
func (Index) ToColorInt(s string) Color                         { return cpt.ToColorInt(s) }
func (Index) ColoredFast(out io.Writer, clr Color, text string) { WrapColorTo(out, clr, text) }
func (Index) DimFast(out io.Writer, text string)                { WrapDimTo(out, text) }
func (Index) HighlightFast(out io.Writer, text string)          { WrapHighlightTo(out, text) }
func (Index) WrapDimToLite(out io.Writer, text string)          { WrapDimToLite(out, text) }
func (Index) WrapColorAndBgTo(out io.Writer, clr, bg Color, text string) {
	WrapColorAndBgTo(out, clr, bg, text)
}

// GetCPT returns a translator about ANSI Escaped Code.
// It may or may not translating colored text depending
// on cabin.GetNoColorMode.
func GetCPT() Translator {
	if states.Env().IsNoColorMode() {
		return &cptNC
	}
	return &cpt
}

func GetCPTC() Translator  { return &cpt }
func GetCPTNC() Translator { return &cptNC }

func ToColorString(clr Color) string { return cpt.toColorString(clr) }
func ToColorInt(s string) Color      { return cpt.ToColorInt(s) }

// Translator _
type Translator interface {
	Translate(s string, initialFg Color) string

	// StripLeftTabsAndColorize(s string) string

	ColoredFast(out io.Writer, clr Color, text string)
	DimFast(out io.Writer, text string)
	HighlightFast(out io.Writer, text string)

	WriteColor(out io.Writer, clr Color)   // echo ansi color bytes for foreground color
	WriteBgColor(out io.Writer, clr Color) // echo ansi color bytes for background color
	Reset(out io.Writer)                   // echo ansi color bytes for resetting color

	Bold(out io.Writer, cb func(out io.Writer))
	Italic(out io.Writer, cb func(out io.Writer))
	Underline(out io.Writer, cb func(out io.Writer))
	Inverse(out io.Writer, cb func(out io.Writer))
	Dim(out io.Writer, cb func(out io.Writer))
	Blink(out io.Writer, cb func(out io.Writer))
	Bg(out io.Writer, bgColor Color, cb func(out io.Writer))

	// TranslateTo(s string, initialState Color) string
	TranslateTo(s string, initialState Color) string // translate a string with html tags to a colored string

	StripLeftTabsAndColorize(s string) string // strip left tabs and colorize the string
	StripLeftTabs(s string) string            // strip left tabs and colorize the string
	StripLeftTabsOnly(s string) string        // strip left tabs only

	stripHTMLTags(s string) string // aggressively strip HTML tags from a string
	ToColorString(clr Color) string
	ToColorInt(s string) Color
}

func GetDummyTranslator() Translator { return dummy } // return Translator for displaying plain text without color

var dummy dummyS

type dummyS struct{}

// stripHTMLTags implements Translator.
func (c dummyS) stripHTMLTags(s string) string {
	panic("unimplemented")
}

func (dummyS) Translate(s string, initialFg Color) string        { return s } //nolint:revive
func (dummyS) TranslateTo(s string, initialState Color) string   { return s }
func (dummyS) ColoredFast(out io.Writer, clr Color, text string) { _, _ = out.Write([]byte(text)) } //nolint:revive
func (dummyS) DimFast(out io.Writer, text string)                { _, _ = out.Write([]byte(text)) }
func (dummyS) HighlightFast(out io.Writer, text string)          { _, _ = out.Write([]byte(text)) }
func (dummyS) WriteColor(out io.Writer, clr Color)               {} //nolint:revive
func (dummyS) WriteBgColor(out io.Writer, clr Color)             {} //nolint:revive
func (dummyS) Reset(out io.Writer)                               {} //nolint:revive
func (dummyS) ToColorInt(s string) Color                         { return cptNC.toColorInt(s) }
func (dummyS) ToColorString(clr Color) string                    { return cptNC.toColorString(clr) }

func (c dummyS) StripLeftTabsAndColorize(s string) string { return cptNC.StripLeftTabsAndColorize(s) }
func (c dummyS) StripLeftTabs(s string) string            { return StripLeftTabs(s) }
func (c dummyS) StripLeftTabsOnly(s string) string        { return StripLeftTabsOnly(s) }
func (c dummyS) StripHTMLTags(s string) string            { return StripHTMLTags(s) }

func (c dummyS) Bold(out io.Writer, cb func(out io.Writer))      { c.bg(out, BgBoldOrBright, cb) }
func (c dummyS) Italic(out io.Writer, cb func(out io.Writer))    { c.bg(out, BgItalic, cb) }
func (c dummyS) Underline(out io.Writer, cb func(out io.Writer)) { c.bg(out, BgUnderline, cb) }
func (c dummyS) Inverse(out io.Writer, cb func(out io.Writer))   { c.bg(out, BgInverse, cb) }
func (c dummyS) Dim(out io.Writer, cb func(out io.Writer))       { c.bg(out, BgDim, cb) }
func (c dummyS) Blink(out io.Writer, cb func(out io.Writer))     { c.bg(out, BgBlink, cb) }
func (c dummyS) Bg(out io.Writer, bgColor Color, cb func(out io.Writer)) {
	c.bg(out, bgColor, cb)
}
func (c dummyS) bg(out io.Writer, bg Color, cb func(out io.Writer)) {
	cb(out)
	return
}

var _ Translator = (*dummyS)(nil)       // ensure cpTranslator implements Translator
var _ Translator = (*cpTranslator)(nil) // ensure cpTranslator implements Translator

type cpTranslator struct {
	noColorMode bool // strip color code simply
}

func (c *cpTranslator) Translate(s string, initialFg Color) string {
	return c.TranslateTo(s, initialFg)
}

func (c *cpTranslator) resetColors(sb *strings.Builder, states []Color) func() { //nolint:revive,gocritic
	return func() {
		var st string
		st = "\x1b[0m"
		_, _ = (*sb).WriteString(st)
		if len(states) > 0 {
			st = fmt.Sprintf("\x1b[%dm", states[len(states)-1])
			_, _ = (*sb).WriteString(st)
		}
	}
}

func (c *cpTranslator) colorize(sb *strings.Builder, states []Color, walker *func(node *html.Node, level int)) func(node *html.Node, clr Color, representation string, level int) { //nolint:revive,gocritic,lll
	return func(node *html.Node, clr Color, representation string, level int) {
		if representation != "" {
			_, _ = (*sb).WriteString(fmt.Sprintf("\x1b[%sm", representation))
		} else {
			_, _ = (*sb).WriteString(fmt.Sprintf("\x1b[%dm", clr))
		}
		states = append(states, clr) //nolint:revive
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			(*walker)(child, level+1)
		}
		states = states[0 : len(states)-1] //nolint:revive
		c.resetColors(sb, states)()
	}
}

func (c *cpTranslator) TranslateTo(s string, initialState Color) string {
	if c.noColorMode {
		return c._ss(s)
	}

	node, err := html.Parse(bufio.NewReader(strings.NewReader(s)))
	if err != nil {
		return c._sz(s)
	}

	return c.translateTo(node, s, initialState)
}

func (c *cpTranslator) translateTo(root *html.Node, source string, initialState Color) string { //nolint:revive,unparam
	states, _ := []Color{initialState}, source //nolint:revive,gocritic
	var sb strings.Builder

	var walker func(node *html.Node, level int)
	colorize := c.colorize(&sb, states, &walker)
	// nilfn := func(node *html.Node, level int) {}
	colorizeIt := func(clr Color) func(node *html.Node, level int) {
		return func(node *html.Node, level int) {
			colorize(node, clr, "", level)
		}
	}
	m := map[string]func(node *html.Node, level int){
		"html": nil, "head": nil, "body": nil,
		"b": colorizeIt(BgBoldOrBright), "strong": colorizeIt(BgBoldOrBright), "em": colorizeIt(BgBoldOrBright),
		"i": colorizeIt(BgItalic), "cite": colorizeIt(BgItalic),
		"u":    colorizeIt(BgUnderline),
		"mark": colorizeIt(BgInverse),
		"del":  colorizeIt(BgStrikeout),
		"dim":  colorizeIt(BgDim),
	}

	walker = func(node *html.Node, level int) {
		switch node.Type {
		case html.DocumentNode, html.DoctypeNode, html.CommentNode:
		case html.ErrorNode:
		case html.ElementNode:
			if fn, ok := m[node.Data]; ok {
				if fn != nil {
					fn(node, level)
					return
				}
			}

			switch node.Data {
			case "font":
				for _, a := range node.Attr {
					if a.Key == "color" {
						clr := c.toColorInt(a.Val)
						colorize(node, clr, "", level)
						return
					}
				}
			case "kbd", "code":
				colorize(node, 51, "51;1", level)
				return
			default:
				// Logger.Debugf("%v, %v, lvl #%d\n", node.Type, node.Data, level)
				// sb.WriteString(node.Data)
				slog.Debug("default node", "data", node.Data)
			}
		case html.TextNode:
			// Logger.Debugf("%v, %v, lvl #%d\n", node.Type, node.Data, level)
			_, _ = sb.WriteString(node.Data)
			return
		default:
			// sb.WriteString(node.Data)
			slog.Debug("default node", "data", node.Data)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walker(child, level+1)
		}
	}

	walker(root, 0)
	return sb.String()
}

func (c *cpTranslator) _sz(s string) string {
	return s
}

func (c *cpTranslator) _ss(s string) string {
	if term.IsAnsiEscaped(s) {
		clean := term.StripEscapes(s)
		return c.stripHTMLTags(clean)
	}
	return c.stripHTMLTags(s)
}

func (c *cpTranslator) StripLeftTabsAndColorize(s string) string {
	return c.stripLeftTabs(s)
}

func (c *cpTranslator) StripLeftTabs(s string) string {
	r := c.stripLeftTabsOnly(s)
	return c.Translate(r, 0)
}

func (c *cpTranslator) stripLeftTabs(s string) string {
	r := c.stripLeftTabsOnly(s)
	return c.Translate(r, 0)
}

func (c *cpTranslator) StripLeftTabsOnly(s string) string {
	return c.stripLeftTabsOnly(s)
}

func (c *cpTranslator) stripLeftTabsOnly(s string) string {
	var lines []string
	tabs := 1000
	var emptyLines []int
	var sb strings.Builder
	var line int
	noLastLF := !strings.HasSuffix(s, "\n")

	scanner := bufio.NewScanner(bufio.NewReader(strings.NewReader(s)))
	for scanner.Scan() {
		str := scanner.Text()
		i, n, allTabs := 0, len(str), true
		for ; i < n; i++ {
			if str[i] != '\t' {
				allTabs = false
				if tabs > i && i > 0 {
					tabs = i
					break
				}
			}
		}
		if i == n && allTabs {
			emptyLines = append(emptyLines, line)
		}
		lines = append(lines, str)
		line++
	}

	pad := strings.Repeat("\t", tabs)
	for i, str := range lines {
		switch {
		case strings.HasPrefix(str, pad):
			_, _ = sb.WriteString(str[tabs:])
		case inIntSlice(i, emptyLines):
		default:
			_, _ = sb.WriteString(str)
		}
		if noLastLF && i == len(lines)-1 {
			break
		}
		_, _ = sb.WriteRune('\n')
	}

	return sb.String()
}

func inIntSlice(i int, slice []int) bool {
	for _, n := range slice {
		if n == i {
			return true
		}
	}
	return false
}

// StripLeftTabsC strips the least left side tab chars from lines.
// StripLeftTabsC strips html tags too.
// At the end, StripLeftTabsC try translate color code in string.
func StripLeftTabsC(s string) string { return cpt.stripLeftTabs(s) }

// StripLeftTabs strips the least left side tab chars from lines.
// StripLeftTabs strips html tags too.
func StripLeftTabs(s string) string { return cptNC.stripLeftTabs(s) }

// StripLeftTabsOnly strips the least left side tab chars from lines.
func StripLeftTabsOnly(s string) string { return cptNC.stripLeftTabsOnly(s) }

// StripHTMLTags aggressively strips HTML tags from a string.
// It will only keep anything between `>` and `<`.
func StripHTMLTags(s string) string { return cptNC.stripHTMLTags(s) }

// Aggressively strips HTML tags from a string.
// It will only keep anything between `>` and `<`.
func (c *cpTranslator) stripHTMLTags(s string) string {
	// Setup a string builder and allocate enough memory for the new string.
	var builder strings.Builder
	builder.Grow(len(s) + utf8.UTFMax)

	in := false // True if we are inside an HTML tag.
	start := 0  // The index of the previous start tag character `<`
	end := 0    // The index of the previous end tag character `>`

	for i, c := range s {
		// If this is the last character and we are not in an HTML tag, save it.
		if end >= start && (i+1) == len(s) {
			_, _ = builder.WriteString(s[end:])
		}

		// Keep going if the character is not `<` or `>`
		if c != htmlTagStart && c != htmlTagEnd {
			continue
		}

		if c == htmlTagStart {
			// Only update the start if we are not in a tag.
			// This make sure we strip out `<<br>` not just `<br>`
			if !in {
				start = i
			}
			in = true

			// Write the valid string between the close and start of the two tags.
			_, _ = builder.WriteString(s[end:start])
			continue
		}
		// else c == htmlTagEnd
		in = false
		end = i + 1
	}
	str := builder.String()
	return str
}

func WrapColorTo(out io.Writer, clr Color, text string) {
	if states.Env().IsNoColorMode() {
		_, _ = out.Write([]byte(text))
		return
	}

	echoColor(out, clr)
	_, _ = out.Write([]byte(text))
	echoResetColor(out)
}

func WrapColorAndBgTo(out io.Writer, clr, bg Color, text string) {
	if states.Env().IsNoColorMode() {
		_, _ = out.Write([]byte(text))
		return
	}

	echoColorAndBg(out, clr, bg)
	_, _ = out.Write([]byte(text))
	echoResetColor(out)
}

func WrapDimTo(out io.Writer, text string) {
	if states.Env().IsNoColorMode() {
		_, _ = out.Write([]byte(text))
		return
	}

	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[2m\x1b[37m%v\x1b[0m", sb.String())

	_, _ = out.Write([]byte("\x1b[2m\x1b[37m"))
	_, _ = out.Write([]byte(text))
	_, _ = out.Write([]byte("\x1b[0m"))
}

func WrapDimToLite(out io.Writer, text string) {
	if states.Env().IsNoColorMode() {
		_, _ = out.Write([]byte(text))
		return
	}

	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[2m\x1b[37m%v\x1b[0m", sb.String())

	_, _ = out.Write([]byte("\x1b[2m"))
	_, _ = out.Write([]byte(text))
	// _, _ = out.Write([]byte("\x1b[0m"))
}

func WrapHighlightTo(out io.Writer, text string) {
	if states.Env().IsNoColorMode() {
		_, _ = out.Write([]byte(text))
		return
	}

	// str = fmt.Sprintf("\x1b[0;1m%v\x1b[0m", str)

	_, _ = out.Write([]byte("\x1b[0;1m"))
	_, _ = out.Write([]byte(text))
	_, _ = out.Write([]byte("\x1b[0m"))
}

// ColoredFast outputs formatted message to stdout with colored ansi codes.
func (c *cpTranslator) ColoredFast(out io.Writer, clr Color, text string) {
	WrapColorTo(out, clr, text)
}

// DimFast outputs formatted message to stdout with colored ansi codes.
func (c *cpTranslator) DimFast(out io.Writer, text string) {
	WrapDimTo(out, text)
}

// HighlightFast outputs formatted message to stdout with colored ansi codes.
func (c *cpTranslator) HighlightFast(out io.Writer, text string) {
	WrapHighlightTo(out, text)
}

func (c *cpTranslator) Bold(out io.Writer, cb func(out io.Writer))      { c.bg(out, BgBoldOrBright, cb) }
func (c *cpTranslator) Italic(out io.Writer, cb func(out io.Writer))    { c.bg(out, BgItalic, cb) }
func (c *cpTranslator) Underline(out io.Writer, cb func(out io.Writer)) { c.bg(out, BgUnderline, cb) }
func (c *cpTranslator) Inverse(out io.Writer, cb func(out io.Writer))   { c.bg(out, BgInverse, cb) }
func (c *cpTranslator) Dim(out io.Writer, cb func(out io.Writer))       { c.bg(out, BgDim, cb) }
func (c *cpTranslator) Blink(out io.Writer, cb func(out io.Writer))     { c.bg(out, BgBlink, cb) }
func (c *cpTranslator) Bg(out io.Writer, bgColor Color, cb func(out io.Writer)) {
	c.bg(out, bgColor, cb)
}

func (c *cpTranslator) bg(out io.Writer, bg Color, cb func(out io.Writer)) {
	if states.Env().IsNoColorMode() {
		cb(out)
		return
	}

	echoBg(out, bg)
	cb(out)
	echoResetColor(out)
}

func (c *cpTranslator) WriteBgColor(out io.Writer, clr Color) { echoBg(out, clr) }
func (c *cpTranslator) WriteColor(out io.Writer, clr Color)   { echoColor(out, clr) }
func (c *cpTranslator) Reset(out io.Writer)                   { echoResetColor(out) }
func (c *cpTranslator) color(out io.Writer, clr Color)        { echoColor(out, clr) } //nolint:unused
func (c *cpTranslator) resetColor(out io.Writer)              { echoResetColor(out) } //nolint:unused

func echoColor(out io.Writer, clr Color) {
	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
	_, _ = out.Write([]byte("\x1b["))
	_, _ = out.Write([]byte(strconv.Itoa(int(clr))))
	_, _ = out.Write([]byte{'m'})
}

func echoColorAndBg(out io.Writer, clr, bg Color) {
	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
	if clr != NoColor {
		_, _ = out.Write([]byte("\x1b["))
		_, _ = out.Write([]byte(strconv.Itoa(int(clr))))
		_, _ = out.Write([]byte{'m'})
	}
	echoBg(out, bg)
}

func echoBg(out io.Writer, bg Color) {
	if bg != NoColor {
		_, _ = out.Write([]byte("\x1b["))
		_, _ = out.Write([]byte(strconv.Itoa(int(bg))))
		_, _ = out.Write([]byte{'m'})
	}
}

func echoResetColor(out io.Writer) { //nolint:unused //no
	// _, _ = fmt.Fprint(os.Stdout, "\x1b[0m")
	_, _ = out.Write([]byte("\x1b[0m"))
}

func (c *cpTranslator) onceInit() {
	onceoptCM.Do(func() {
		cptCM = map[string]Color{
			"black":     FgBlack,
			"red":       FgRed,
			"green":     FgGreen,
			"yellow":    FgYellow,
			"blue":      FgBlue,
			"magenta":   FgMagenta,
			"cyan":      FgCyan,
			"lightgray": FgLightGray, "light-gray": FgLightGray,
			"darkgray": FgDarkGray, "dark-gray": FgDarkGray,
			"lightred": FgLightRed, "light-red": FgLightRed,
			"lightgreen": FgLightGreen, "light-green": FgLightGreen,
			"lightyellow": FgLightYellow, "light-yellow": FgLightYellow,
			"lightblue": FgLightBlue, "light-blue": FgLightBlue,
			"lightmagenta": FgLightMagenta, "light-magenta": FgLightMagenta,
			"lightcyan": FgLightCyan, "light-cyan": FgLightCyan,
			"white": FgWhite,
		}
		cptNM = make(map[Color]string)
		for k, v := range cptCM {
			cptNM[v] = k
		}
	})
}

func (c *cpTranslator) ToColorString(clr Color) string { return c.toColorString(clr) }

func (c *cpTranslator) toColorString(clr Color) string { //nolint:revive
	c.onceInit()
	if ss, ok := cptNM[clr]; ok {
		return ss
	}
	return "gray"
}

func (c *cpTranslator) ToColorInt(s string) Color { return c.toColorInt(s) }

func (c *cpTranslator) toColorInt(s string) Color { //nolint:revive
	c.onceInit()
	if i, ok := cptCM[strings.ToLower(s)]; ok {
		return i
	}
	return Color(0)
}

const (
	htmlTagStart = 60 // Unicode `<`
	htmlTagEnd   = 62 // Unicode `>`
)

var (
	onceoptCM sync.Once
	cptCM     map[string]Color
	cptNM     map[Color]string

	// cptHTM    map[string]func(node *html.Node, level int)

	// onceColorPrintTranslator sync.Once

	cpt   cpTranslator
	cptNC = cpTranslator{noColorMode: true}
)
