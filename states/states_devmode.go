package states

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hedzr/is/dir"
)

func init() { DetectDevModeFileOnce() }

func DetectDevModeFile() (isCmdrV2, devModeFilePresent, devMode bool) {
	// var err error
	// var fi os.FileInfo
	// if fi, devModeFilePresent, err = dir.Exists(".devmode"); devModeFilePresent && err == nil {
	// 	//
	// } else {
	// 	_, _ = err, fi
	// }

	d := dir.GetCurrentDir()
	cliName := filepath.Dir(d)
	prjName := filepath.Dir(cliName)
	if filepath.Base(cliName) == "cli" {
		d = prjName
	}

	// isCmdrV2 := false
	devModeFile := filepath.Join(d, ".dev-mode")
	if devModeFilePresent = dir.FileExists(devModeFile); devModeFilePresent {
		devMode = true
	} else {
		devModeFile := filepath.Join(d, ".devmode")
		if devModeFilePresent = dir.FileExists(devModeFile); devModeFilePresent {
			devMode = true
		}
	}
	if dir.FileExists("go.mod") {
		data, err := os.ReadFile("go.mod")
		if err != nil {
			return
		}
		content := string(data)

		// if strings.Contains(content, "github.com/hedzr/cmdr/v2/pkg/") {
		// 	devMode = false
		// }

		// I am tiny-app in cmdr/v2, I will be launched in dev-mode always
		if strings.Contains(content, "module github.com/hedzr/cmdr/v2") {
			isCmdrV2, devMode = true, true
		}
	}

	if e, ok := env.(*minimalEnv); ok {
		e.devModeFilePresent = devModeFilePresent
		e.SetDevMode(devMode)
	}
	return
}

var onceDev sync.Once
var devMode bool
var devModeFilePresent bool
