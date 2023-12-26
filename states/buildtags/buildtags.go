package buildtags

import (
	"debug/buildinfo"
	"strings"

	"github.com/hedzr/is/stringtool"
)

// IsBuildTagExists detects if a tag specified at building time.
func IsBuildTagExists(tag string) (yes bool) {
	file := stringtool.GetExecutablePath()

	var inf *buildinfo.BuildInfo
	var err error
	if inf, err = buildinfo.ReadFile(file); err != nil {
		return
	}

	for _, d := range inf.Settings {
		// fmt.Printf("    - %q: %v\n", d.Key, d.Value)
		if d.Key == "-tags" {
			for _, r := range strings.Split(d.Value, ",") {
				if yes = r == tag; yes {
					return
				}
			}
		}
	}
	return
}
