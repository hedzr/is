package is

import (
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/hedzr/is/states"
	"github.com/hedzr/is/states/buildtags"
	"github.com/hedzr/is/states/isdelve"
	"github.com/hedzr/is/states/trace"
)

// InDevMode return the devmode state.
//
// For cmdr app, devmode means a '.devmode' file is detected
// at work-dir. You could touch a '.devmode' file in your
// app's src-tree folder to enable this state automatically.
//
// Detecting devmode file and state is invoked at app starting up.
// If you wanna disable it, using build tag '-tags="nodetectdevmode"'.
// If you have disabled it but still want it to be invoked manually,
// try [states.DetectDevModeFile] or do it youself.
func InDevMode() bool { return states.Env().InDevMode() }

// InDebugging and DebuggerAttached returns status if golang debugger 'dlv' is attached.
//
// **Since v0.6.0**, InDebugging checks if the parent process is
// 'dlv' or not. It supports Linux, Darwin and Windows currently.
//
// If you're looking for old behavior, DebugBuild() returns if the
// executable is built with delve tag:
//
//	 When you runs `go build` with `-tags=delve` options. eg:
//
//		go run -tags=delve ./cli
//		go build -tags=delve -o my-app ./cli
//
//	 The executable will hold a 'isdelve' tag. For more details please goto
//	 https://stackoverflow.com/questions/47879070/how-can-i-see-if-the-goland-debugger-is-running-in-the-program
//
// Performance Note:
//
// For security reason InDebugging checks the debuggers lively. It might
// take unexpected times for detecting.
//
// Debug States:
//
//  1. is.InDebugging, loaded by a golang debugger (dlv) at runtime?
//  2. is.DebugBuild, build tags decide it
//  3. is.DebugMode, set by command line flag `--debug`
func InDebugging() bool {
	return states.IsUnderDebugger()
}

// DebuggerAttached returns status if debugger attached or is a debug build.
//
// See also InDebugging.
func DebuggerAttached() bool {
	return states.IsUnderDebugger()
}

// InDebugMode and DebugMode returns true if:
//
//   - a debug build, see the DebugBuild.
//   - SetDebugMode(true) called.
//
// **NOTE** Since v0.6.0, InDebugMode does not check DebuggerAttached.
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
	switch runtime.GOOS {
	case "windows":
		if strings.HasSuffix(args[0], ".test.exe") {
			for _, s := range args {
				if strings.HasPrefix(s, "-test.") {
					return true
				}
			}
		}
	default:
		if strings.HasSuffix(args[0], ".test") ||
			strings.Contains(args[0], "/T/___Test") {
			// [0] = /var/folders/td/2475l44j4n3dcjhqbmf3p5l40000gq/T/go-build328292371/b001/exe/main
			// !strings.Contains(SavedOsArgs[0], "/T/go-build")

			for _, s := range args {
				if strings.HasPrefix(s, "-test.") {
					return true
				}
			}
		}
	}
	return false
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

func InBenchmark() bool { return isInBench(os.Args) }

func isInBench(args []string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-test.bench") || strings.HasPrefix(arg, "-bench") {
			return true
		}
		// if strings.HasPrefix(arg, "-test.bench=") {
		// 	// ignore the benchmark name after an underscore
		// 	bench = strings.SplitN(arg[12:], "_", 2)[0]
		// 	break
		// }
	}
	return false
}

// InDevelopingTime detects whether is in developing time (debugging or testing).
//
// If the main program has been built as an executable binary, we
// would assume which is not in developing time.
//
// If GetDebugMode() is true, that's in developing time too.
func InDevelopingTime() (status bool) {
	return InDebugMode() || InTesting() || InBenchmark() || InDebugging()
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

// DevMode return the devmode state.
//
// For cmdr app, devmode means a '.devmode' file is detected
// at work-dir. You could touch a '.devmode' file in your
// app's src-tree folder to enable this state automatically.
//
// Detecting devmode file and state is invoked at app starting up.
// If you wanna disable it, using build tag '-tags="nodetectdevmode"'.
// If you have disabled it but still want it to be invoked manually,
// try [states.DetectDevModeFile] or do it youself.
func DevMode() bool     { return Env().InDevMode() }
func SetDevMode(b bool) { Env().SetDevMode(b) } // set devMode state

// DevModeFilePresent returns a state to identify ".devmode" (or ".dev-mode")
// file is detected. This state relyes on once [states.DetectDevModeFile]
// was invoked.
func DevModeFilePresent() bool { return Env().IsDevModeFilePresent() }

func DebugMode() bool        { return Env().GetDebugMode() }  // is debug build, or is CLI debug mode enabled (by `--debug`)?
func GetDebugLevel() int     { return Env().GetDebugLevel() } // for debug build, return the debug level integer, 0-9
func SetDebugMode(b bool)    { Env().SetDebugMode(b) }        // sets debug state
func SetDebugLevel(hits int) { Env().SetDebugLevel(hits) }    // sets debug level

func Tracing() bool          { return trace.IsEnabled() }     // is tracing-flag true in trace package
func TraceMode() bool        { return Env().GetTraceMode() }  // is CLI trace mode enabled (by `--trace`)? or is tracing-flag true in trace package
func GetTraceLevel() int     { return Env().GetTraceLevel() } // return the trace level integer, 0-9
func SetTraceMode(b bool)    { Env().SetTraceMode(b) }        // sets trace state
func SetTraceLevel(hits int) { Env().SetTraceLevel(hits) }    // sets trace level

func SetOnDeeModeChanged(funcs ...states.OnChanged) { Env().SetOnDevModeChanged(funcs...) } // sets ondebugchanged callbacks
func SetOnDebugChanged(funcs ...states.OnChanged)   { Env().SetOnDebugChanged(funcs...) }   // sets ondebugchanged callbacks
func SetOnTraceChanged(funcs ...states.OnChanged)   { Env().SetOnTraceChanged(funcs...) }   // sets ontracechanged callbacks
func SetOnVerboseChanged(funcs ...states.OnChanged) { Env().SetOnVerboseChanged(funcs...) } // sets onverbosechanged callbacks
func SetOnQuietChanged(funcs ...states.OnChanged)   { Env().SetOnQuietChanged(funcs...) }   // sets onquietchanged callbacks
func SetOnNoColorChanged(funcs ...states.OnChanged) { Env().SetOnNoColorChanged(funcs...) } // sets onnocolorchanged callbacks

// States or Env returns a minimal environment settings for a typical CLI app.
//
// See also [states.CmdrMinimal].
func States() states.CmdrMinimal           { return states.Env() }
func Env() states.CmdrMinimal              { return states.Env() }       // States or Env returns a minimal environment settings for a typical CLI app.
func UpdateEnvWith(env states.CmdrMinimal) { states.UpdateEnvWith(env) } // If u will, use ur own env

// Detected returns the detected state associated with the given state name.
//
// = State, see also it.
func Detected(stateName string) (state bool) { return State(stateName) }

// State returns a state with the given state name.
func State(stateName string) (state bool) {
	initmstates()
	if fn, ok := mstates[stateName]; ok {
		return fn()
	}
	return
}

// RegisterStateGetter allows integrating your own detector into State(name) bool.
//
// For example:
//
//	func customState() bool { reutrn ... }
//	is.RegisterStateGetter("custom", customState)
//	println(is.State("custom"))
func RegisterStateGetter(state string, getter func() bool) {
	initmstates()
	mstates[state] = getter
}

func initmstates() {
	oncemstates.Do(func() {
		mstates = make(map[string]func() bool)

		mstates["debug"] = states.Env().GetDebugMode
		mstates["trace"] = states.Env().GetTraceMode
		mstates["verbose"] = states.Env().IsVerboseMode
		mstates["quiet"] = states.Env().IsQuietMode
		mstates["no-color"] = states.Env().IsNoColorMode

		mstates["docker-build"] = DockerBuild
		mstates["k8s-build"] = K8sBuild
		mstates["istio-build"] = IstioBuild
		mstates["debug-build"] = DebugBuild
		mstates["verbose-build"] = VerboseBuild

		mstates["in-docker"] = InDocker
		mstates["in-k8s"] = InK8s
		mstates["in-istio"] = InIstio

		mstates["in-vscode-terminal"] = InVscodeTerminal

		mstates["in-testing"] = InTesting
		mstates["in-developing-time"] = InDevelopingTime

		mstates["in-tracing"] = InTracing
		mstates["has-debugger"] = InDebugMode
		mstates["in-debugging"] = InDebugMode // note that we cannot detect if a debugger attached really
	})
}

var (
	mstates     map[string]func() bool
	oncemstates sync.Once
)

// FuncPtrSame compares two functors if they are same, with an unofficial way.
func FuncPtrSame[T any](fn1, fn2 T) bool {
	sf1 := reflect.ValueOf(fn1)
	sf2 := reflect.ValueOf(fn2)
	if sf1.Kind() != reflect.Func || sf2.Kind() != reflect.Func {
		return false
	}
	return sf1.Pointer() == sf2.Pointer()
}
