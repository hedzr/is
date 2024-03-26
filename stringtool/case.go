package stringtool

import (
	"unicode"
)

func ToCamelCase(s string) string {
	if s != "" {
		return makeSnakeCase(s, 0, true, false)
	}
	return s
}

func ToSmallCamelCase(s string) string {
	if s != "" {
		return makeSnakeCase(s, 0, true, true)
	}
	return s
}

func ToKebabCase(s string) string {
	if s != "" {
		return makeSnakeCase(s, '-', false, false)
	}
	return s
}

func ToSnakeCase(s string) string {
	if s != "" {
		return makeSnakeCase(s, '_', false, false)
	}
	return s
}

func makeSnakeCase(s string, delimiter rune, capitalize, lower1stWord bool) string {
	if s != "" {
		a := wordSplitter(s)
		if capitalize {
			for i, word := range a {
				if lower1stWord && i == 0 {
					a[i] = makeLowerCase1st(word)
				} else {
					a[i] = makeCapitalize1st(word)
				}
			}
		} else {
			for i, word := range a {
				a[i] = makeLowerCase1st(word)
			}
		}

		var r []rune
		for i, word := range a {
			if i > 0 && delimiter != 0 {
				r = append(r, delimiter)
			}
			r = append(r, word...)
		}
		return string(r)
	}
	return s
}

// CapitalizeFirstLetter make first letter of a string to upper-case.
func CapitalizeFirstLetter(s string) string {
	return string(makeCapitalize1st([]rune(s)))
}

// ToExportedName converts any name to Golang Exported Name.
//
// Basically, it is a kebab/snake-case to Camel-case transformer.
func ToExportedName(s string) string {
	return toExportedName(s)
}

func toExportedName(s string) string {
	if s != "" {
		a := wordSplitter(s)
		for i, word := range a {
			a[i] = makeCapitalize1st(word)
		}

		var r []rune
		for _, word := range a {
			r = append(r, word...)
		}
		return string(r)
	}
	return s
}

func wordSplitter(s string) (result [][]rune) {
	runes := []rune(s)
	var word []rune
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, word)
			}
			word = nil
		} else if r == '-' || r == '_' || unicode.IsSpace(r) {
			if i > 0 {
				result = append(result, word)
			}
			word = nil
			continue
		}
		word = append(word, r)
	}
	if len(word) > 0 {
		result = append(result, word)
	}
	return
}

func makeCapitalize1st(r []rune) (ret []rune) {
	if len(r) > 0 {
		ret = append(ret, unicode.ToUpper(r[0]))
		ret = append(ret, r[1:]...)
		return
	}
	return r
}

func makeLowerCase1st(r []rune) (ret []rune) {
	if len(r) > 0 {
		ret = append(ret, unicode.ToLower(r[0]))
		ret = append(ret, r[1:]...)
		return
	}
	return r
}
