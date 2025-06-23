package term

import (
	"github.com/hedzr/is/term/chk"
)

// ToBool translate a value (int, bool, string) to boolean
func ToBool(val any, defaultVal ...bool) (ret bool) { return chk.ToBool(val, defaultVal...) }

func StringToBool(val string, defaultVal ...bool) (ret bool) {
	return chk.StringToBool(val, defaultVal...)
}
