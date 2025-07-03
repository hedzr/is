//go:build !nodetectdevmode
// +build !nodetectdevmode

package states

const DetectDevModeFileEnabled = true

func DetectDevModeFileOnce() (isCmdrV2, devModeFilePresent, devmode bool) {
	// println("(DetectDevModeFileOnce)")
	// onceDev.Do(func() {
	isCmdrV2, devModeFilePresent, devmode = DetectDevModeFile()
	// println("(DetectDevModeFileOnce)", "END.", devModeFilePresent, devmode)
	// })
	// if e, ok := env.(*minimalEnv); ok {
	// 	devmode = e.devMode
	// 	devModeFilePresent = e.devModeFilePresent

	// 	data, err := os.ReadFile("go.mod")
	// 	if err != nil {
	// 		return
	// 	}
	// 	content := string(data)

	// 	// I am tiny-app in cmdr/v2, I will be launched in dev-mode always
	// 	if strings.Contains(content, "module github.com/hedzr/cmdr/v2") {
	// 		isCmdrV2 = true
	// 	}
	// }
	return
}
