package is

import (
	"os"
	"strings"

	"github.com/hedzr/is/states"
	"github.com/hedzr/is/states/buildtags"
	"github.com/hedzr/is/states/isdelve"
	"github.com/hedzr/is/states/trace"
)

// InDebugging returns status if debugger attached or is a debug build.
//
// To enable the debugger attached mode for cmdr, run `go build` with `-tags=delve` options. eg:
//
//	go run -tags=delve ./cli
//	go build -tags=delve -o my-app ./cli
//
// For Goland, you can enable this under 'Run/Debug Configurations', by adding the following into 'Go tool arguments:'
//
//	-tags=delve
//
// InDebugging() is a synonym to DebuggerAttached().
//
// NOTE that `isdelve` algor is from https://stackoverflow.com/questions/47879070/how-can-i-see-if-the-goland-debugger-is-running-in-the-program
//
// To check wildly like cli app debug mode (via --debug), call InDebugMode.
//
// To find parent process is dlv (that is, detecting a real debugger present), another library needed.
func InDebugging() bool {
	return states.Env().InDebugging() // isdelve.Enabled
}

// DebuggerAttached returns status if debugger attached or is a debug build.
//
// See also InDebugging.
//
// # To check wildly like cli app debug mode (via --debug), call InDebugMode
//
// To find parent process is dlv (that is, detecting a real debugger present), another library needed.
func DebuggerAttached() bool {
	return states.Env().InDebugging() // isdelve.Enabled
}

// InDebugMode returns if:
//
//   - debugger attached
//   - a debug build
//   - SetDebugMode(true) called.
//
// To find parent process is dlv (that is, detecting a real debugger present), another library needed.
func InDebugMode() bool {
	return states.Env().GetDebugMode() // isdelve.Enabled
}

// InTracing tests if trace.IsEnabled() or env.traceMode (cli app trace mode via --trace).
//
// See the github.com/is/states/trace package.
func InTracing() bool {
	return states.Env().GetTraceMode()
}

// InTestingT detects whether is running under 'go test' mode
func InTestingT(args []string) bool {
	if !strings.HasSuffix(args[0], ".test") &&
		!strings.Contains(args[0], "/T/___Test") {

		// [0] = /var/folders/td/2475l44j4n3dcjhqbmf3p5l40000gq/T/go-build328292371/b001/exe/main
		// !strings.Contains(SavedOsArgs[0], "/T/go-build")

		for _, s := range args {
			if s == "-test.v" || s == "-test.run" {
				return true
			}
		}
		return false

	}
	return true
}

// InTesting detects whether is running under go test mode
func InTesting() bool {
	return InTestingT(os.Args)
	// if !strings.HasSuffix(tool.SavedOsArgs[0], ".test") &&
	//	!strings.Contains(tool.SavedOsArgs[0], "/T/___Test") {
	//
	//	// [0] = /var/folders/td/2475l44j4n3dcjhqbmf3p5l40000gq/T/go-build328292371/b001/exe/main
	//	// !strings.Contains(SavedOsArgs[0], "/T/go-build")
	//
	//	for _, s := range tool.SavedOsArgs {
	//		if s == "-test.v" || s == "-test.run" {
	//			return true
	//		}
	//	}
	//	return false
	//
	// }
	// return true
}

// InDevelopingTime detects whether is in developing time (debugging or testing).
//
// If the main program has been built as an executable binary, we
// would assume which is not in developing time.
//
// If GetDebugMode() is true, that's in developing time too.
func InDevelopingTime() (status bool) {
	return InDebugging() || InDebugMode() || InTesting()
}

// InDockerEnvSimple detects whether is running within docker
// container environment.
//
// InDockerEnvSimple finds if `/.dockerenv` exists or not.
func InDockerEnvSimple() (status bool) {
	return isRunningInDockerContainer()
}

func isRunningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then the viewer is running
	// from inside a container so return true

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

// InVscodeTerminal tests if running under visual studio code integrated terminal
func InVscodeTerminal() bool {
	return os.Getenv("VSCODE_INJECTION") == "1"
}

// InK8s detects if the service is running under k8s environment.
func InK8s() bool {
	return os.Getenv("KUBERNETES_SERVICE_HOST") != "" || buildtags.IsK8sBuild()
}

// InK8sYN is yet another DetectInK8s impl
func InK8sYN() bool {
	return fileExists("/var/run/secrets/kubernetes.io") || buildtags.IsK8sBuild()
}

// InIstio detects if the service is running under istio injected.
//
// ### IMPORTANT
//
// To make this detector work properly, you must mount a DownwordAPI
// volume to your container/pod. See also:
//
// https://kubernetes.io/en/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/
func InIstio() bool {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		const kf = "/etc/podinfo/labels"
		if fileExists(kf) {
			if data, err := readFile(kf); err == nil {
				// lines:=strings.Split(string(data),"\n")
				if strings.Contains(string(data), "service.istio.io/canonical-name") {
					return true
				}
			}
		}
	}
	return buildtags.IsIstioBuild()
}

// InDocker detects if the service is running under docker environment.
//
// We tests these two conditions:
//
//  1. find if `/.dockerenv` exists or not.
//  2. `docker` in buildtags
func InDocker() bool {
	if fileExists("/.dockerenv") {
		return true
	}
	return buildtags.IsDockerBuild()
}

func DockerBuild() bool  { return buildtags.IsDockerBuild() } // need build tag 'docker' present
func K8sBuild() bool     { return buildtags.IsK8sBuild() }    // need build tag 'k8s' present
func IstioBuild() bool   { return buildtags.IsIstioBuild() }  // need build tag 'istio' present
func DebugBuild() bool   { return isdelve.Enabled }           // is debug build? need build tag 'delve' present
func VerboseBuild() bool { return buildtags.VerboseEnabled }  // is verbose build? need build tag 'verbose' present

func VerboseModeEnabled() bool { return Env().IsVerboseMode() }  // is verbose build, or is CLI Verbose mode enabled (by `--verbose`)?
func GetVerboseLevel() int     { return Env().CountOfVerbose() } // returns verbose state level
func SetVerboseMode(b bool)    { Env().SetVerboseMode(b) }       // sets verbose state
func SetVerboseLevel(hits int) { Env().SetVerboseCount(hits) }   // sets verbose level

func QuietModeEnabled() bool { return Env().IsQuietMode() }  // is quiet build, or is CLI Quiet mode enabled (by `--verbose`)?
func GetQuietLevel() int     { return Env().CountOfQuiet() } // returns quiet state level
func SetQuietMode(b bool)    { Env().SetQuietMode(b) }       // sets quiet state
func SetQuietLevel(hits int) { Env().SetQuietCount(hits) }   // sets quiet level

func NoColorMode() bool        { return Env().IsNoColorMode() }  // plain mode (non-colorful mode)
func GetNoColorLevel() int     { return Env().CountOfNoColor() } // returns no-color state level
func SetNoColorMode(b bool)    { Env().SetNoColorMode(b) }       // sets no-color state
func SetNoColorLevel(hits int) { Env().SetNoColorCount(hits) }   // setd no-color level

func DebugMode() bool        { return Env().GetDebugMode() }  // is debug build, or is CLI debug mode enabled (by `--debug`)?
func GetDebugLevel() int     { return Env().GetDebugLevel() } // for debug build, return the debug level integer, 0-9
func SetDebugMode(b bool)    { Env().SetDebugMode(b) }        // sets debug state
func SetDebugLevel(hits int) { Env().SetDebugLevel(hits) }    // sets debug level

func Tracing() bool          { return trace.IsEnabled() }     // is tracing-flag true in trace package
func TraceMode() bool        { return Env().GetTraceMode() }  // is CLI trace mode enabled (by `--trace`)? or is tracing-flag true in trace package
func GetTraceLevel() int     { return Env().GetTraceLevel() } // return the trace level integer, 0-9
func SetTraceMode(b bool)    { Env().SetTraceMode(b) }        // sets trace state
func SetTraceLevel(hits int) { Env().SetTraceLevel(hits) }    // sets trace level

// States or Env returns a minimal environment settings for a typical CLI app.
//
// See also [states.CmdrMinimal].
func States() states.CmdrMinimal           { return states.Env() }
func Env() states.CmdrMinimal              { return states.Env() }       // States or Env returns a minimal environment settings for a typical CLI app.
func UpdateEnvWith(env states.CmdrMinimal) { states.UpdateEnvWith(env) } // If u will, use ur own env
