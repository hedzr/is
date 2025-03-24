// Copyright © 2023 Hedzr Yeh.

package strings

import (
	"regexp"
	"strconv"
	"unicode/utf8"
)

// UnescapeUnicode decodes \uxxxx as unicode string.
func UnescapeUnicode(b []byte) string {
	b = reFind.ReplaceAllFunc(b, expandUnicode)
	return string(b)
}

// UnescapeUnicodeInYamlDoc decodes \uxxxx as unicode string.
//
// It assumes the input doc is well-formatted yaml.
func UnescapeUnicodeInYamlDoc(b []byte) string {
	b = reFindInYaml.ReplaceAllFunc(b, expandUnicode)
	return string(b)
}

// var reFind = regexp.MustCompile(`^\s*[^\s\:]+\:\s*["']?.*\\u.*["']?\s*$`)
var reFind = regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)
var reFindInYaml = regexp.MustCompile(`[^\s\:]+\:\s*["']?.*\\u.*["']?`)

var reFindU = regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)

func expandUnicode(line []byte) []byte {
	return reFindU.ReplaceAllFunc(line, expandUnicodeRune)
}

func expandUnicodeRune(esc []byte) []byte {
	ri, _ := strconv.ParseInt(string(esc[2:]), 16, 32)
	r := rune(ri)
	repr := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(repr, r)
	return repr
}
