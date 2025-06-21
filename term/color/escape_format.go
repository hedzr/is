package color

import (
	"io"
	"strconv"

	"github.com/hedzr/is/states"
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

func echoColor(out io.Writer, clr Color) {
	_, _ = out.Write([]byte(clr.Color()))
}

func echoColor16(out io.Writer, clr Color16) {
	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
	_, _ = out.Write([]byte("\x1b["))
	_, _ = out.Write([]byte(strconv.Itoa(int(clr))))
	_, _ = out.Write([]byte{'m'})
}

func echoColorAndBg(out io.Writer, clr, bg Color) {
	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
	if clr != NoColor {
		echoColor(out, clr)
	}
	echoBg(out, bg)
}

func echoColorAndBg16(out io.Writer, clr, bg Color16) {
	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
	if clr != NoColor {
		_, _ = out.Write([]byte("\x1b["))
		_, _ = out.Write([]byte(strconv.Itoa(int(clr))))
		_, _ = out.Write([]byte{'m'})
	}
	echoBg(out, bg)
}

func echoBg(out io.Writer, bg Color) {
	_, _ = out.Write([]byte(bg.Color()))
}

func echoBg16(out io.Writer, bg Color16) {
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
