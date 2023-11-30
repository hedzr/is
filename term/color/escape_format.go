package color

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"golang.org/x/net/html"

	"github.com/hedzr/is/states"
	"github.com/hedzr/is/term"
)

// GetCPT returns a translator about ANSI Escaped Code.
// It may or may not translating colored text depending
// on cabin.GetNoColorMode.
func GetCPT() Translator {
	if states.Env().IsNoColorMode() {
		return &cptNC
	} else {
		return &cpt
	}
}

func GetCPTC() Translator {
	return &cpt
}

func GetCPTNC() Translator {
	return &cptNC
}

// Translator _
type Translator interface {
	Translate(s string, initialFg Color) string

	// StripLeftTabsAndColorize(s string) string

	ColoredFast(out io.Writer, clr Color, text string)
	DimFast(out io.Writer, text string)
	HighlightFast(out io.Writer, text string)
}

func GetDummyTranslator() Translator { return dummy }

var dummy dummyS

type dummyS struct{}

func (dummyS) Translate(s string, initialFg Color) string        { return s }
func (dummyS) ColoredFast(out io.Writer, clr Color, text string) { _, _ = out.Write([]byte(text)) }
func (dummyS) DimFast(out io.Writer, text string)                { _, _ = out.Write([]byte(text)) }
func (dummyS) HighlightFast(out io.Writer, text string)          { _, _ = out.Write([]byte(text)) }

type cpTranslator struct {
	noColorMode bool // strip color code simply
}

func (c *cpTranslator) Translate(s string, initialFg Color) string {
	return c.TranslateTo(s, initialFg)
}

func (c *cpTranslator) resetColors(sb *strings.Builder, states []Color) func() {
	return func() {
		var st string
		st = "\x1b[0m"
		(*sb).WriteString(st)
		if len(states) > 0 {
			st = fmt.Sprintf("\x1b[%dm", states[len(states)-1])
			(*sb).WriteString(st)
		}
	}
}

func (c *cpTranslator) colorize(sb *strings.Builder, states []Color, walker *func(node *html.Node, level int)) func(node *html.Node, clr Color, representation string, level int) {
	return func(node *html.Node, clr Color, representation string, level int) {
		if representation != "" {
			(*sb).WriteString(fmt.Sprintf("\x1b[%sm", representation))
		} else {
			(*sb).WriteString(fmt.Sprintf("\x1b[%dm", clr))
		}
		states = append(states, clr)
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			(*walker)(child, level+1)
		}
		states = states[0 : len(states)-1]
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

func (c *cpTranslator) translateTo(root *html.Node, source string, initialState Color) string {
	states := []Color{initialState}
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
			}
		case html.TextNode:
			// Logger.Debugf("%v, %v, lvl #%d\n", node.Type, node.Data, level)
			sb.WriteString(node.Data)
			return
		default:
			// sb.WriteString(node.Data)
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
	if term.IsTtyEscaped(s) {
		clean := term.StripEscapes(s)
		return c.stripHTMLTags(clean)
	}
	return c.stripHTMLTags(s)
}

func (c *cpTranslator) StripLeftTabsAndColorize(s string) string {
	return c.stripLeftTabs(s)
}

func (c *cpTranslator) stripLeftTabs(s string) string {
	r := c.stripLeftTabsOnly(s)
	return c.Translate(r, 0)
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
			sb.WriteString(str[tabs:])
		case inIntSlice(i, emptyLines):
		default:
			sb.WriteString(str)
		}
		if noLastLF && i == len(lines)-1 {
			break
		}
		sb.WriteRune('\n')
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
		if (i+1) == len(s) && end >= start {
			builder.WriteString(s[end:])
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
			builder.WriteString(s[end:start])
			continue
		}
		// else c == htmlTagEnd
		in = false
		end = i + 1
	}
	s = builder.String()
	return s
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

// ColoredFast outputs formatted message to stdout while logger level
// less than slog.Warn.
// For slog.SetLevel(slog.Error), the text will be discarded.
func (c *cpTranslator) ColoredFast(out io.Writer, clr Color, text string) {
	WrapColorTo(out, clr, text)
}

func (c *cpTranslator) DimFast(out io.Writer, text string) {
	WrapDimTo(out, text)
}

func (c *cpTranslator) HighlightFast(out io.Writer, text string) {
	WrapHighlightTo(out, text)
}

func (c *cpTranslator) color(out io.Writer, clr Color) { echoColor(out, clr) }
func (c *cpTranslator) resetColor(out io.Writer)       { echoResetColor(out) }

func echoColor(out io.Writer, clr Color) {
	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
	_, _ = out.Write([]byte("\x1b["))
	_, _ = out.Write([]byte(strconv.Itoa(int(clr))))
	_, _ = out.Write([]byte{'m'})
}

func echoColorAndBg(out io.Writer, clr, bg Color) {
	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
	_, _ = out.Write([]byte("\x1b["))
	_, _ = out.Write([]byte(strconv.Itoa(int(clr))))
	_, _ = out.Write([]byte{'m'})
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

func (c *cpTranslator) toColorInt(s string) Color {
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
	})
	if i, ok := cptCM[strings.ToLower(s)]; ok {
		return Color(i)
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

	// cptHTM    map[string]func(node *html.Node, level int)

	// onceColorPrintTranslator sync.Once

	cpt   cpTranslator
	cptNC = cpTranslator{noColorMode: true}
)
