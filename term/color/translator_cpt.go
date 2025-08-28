package color

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/hedzr/is/states"
	"github.com/hedzr/is/term"
	"golang.org/x/net/html"
)

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

var _ Translator = (*cpTranslator)(nil) // ensure cpTranslator implements Translator

type cpTranslator struct {
	noColorMode     bool // strip color code simply
	noLeadingSpaces bool
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

	if c.noLeadingSpaces {
		node, err := html.Parse(bufio.NewReader(strings.NewReader(string(s))))
		if err != nil {
			return c._sz(s)
		}

		return c.translateTo(node, "", s, initialState)
	}

	runes, pos := []rune(s), 0
	for unicode.IsSpace(runes[pos]) {
		pos++
	}
	cs := string(runes[pos:])
	node, err := html.Parse(bufio.NewReader(strings.NewReader(cs)))
	if err != nil {
		return c._sz(s)
	}

	return c.translateTo(node, string(runes[:pos]), cs, initialState)
}

func (c *cpTranslator) translateTo(root *html.Node, leading, source string, initialState Color) string {
	states, _ := []Color{initialState}, source //nolint:revive,gocritic
	var sb strings.Builder

	_, _ = sb.WriteString(leading)

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
		"hide": colorizeIt(BgHidden),
		"del":  colorizeIt(BgStrikeout),
		"dim":  colorizeIt(BgDim),
		"dbl":  colorizeIt(Color16(21)), // doubly underlined
		"over": colorizeIt(Color16(53)), // overlined
		// "box":  colorizeIt(Color16(51)), // framed (not work)
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
				// printf "\033[%sm%s\033[0m\n" "51;1" "text here"
				// draw a frame arround the character(s), rarely supported.
				colorize(node, Reset, "51;1", level)
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
	return c.Translate(r, Reset)
}

func (c *cpTranslator) stripLeftTabs(s string) string {
	r := c.stripLeftTabsOnly(s)
	return c.Translate(r, Reset)
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
			if in && c == htmlTagEnd {
				// the end-tag is in recoganizing, nothing to done
			} else {
				_, _ = builder.WriteString(s[end:])
			}
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
		if in {
			in = false
			end = i + 1
		}
	}
	str := builder.String()
	return str
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
	return Reset
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
