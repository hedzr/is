package term

import (
	"regexp"
	"strings"
)

// IsTtyEscaped detects a string if it contains ansi color escaped sequences
func IsTtyEscaped(s string) bool { return isTtyEscaped(s) }
func isTtyEscaped(s string) bool { return strings.Contains(s, "\x1b[") || strings.Contains(s, "\x9b[") }

// StripEscapes removes any ansi color escaped sequences from a string
func StripEscapes(str string) (strCleaned string) { return stripEscapes(str) }

// var reStripEscapesOld = regexp.MustCompile(`\x1b\[[0-9,;]+m`)

const ansi = "[\u001b\u009b][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var reStripEscapes = regexp.MustCompile(ansi)

func stripEscapes(str string) (strCleaned string) {
	strCleaned = reStripEscapes.ReplaceAllString(str, "")
	return
}
