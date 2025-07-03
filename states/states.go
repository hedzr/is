package states

import (
	"reflect"
	"runtime"

	"github.com/hedzr/is/states/buildtags"
	"github.com/hedzr/is/states/isdelve"
	"github.com/hedzr/is/states/trace"
)

// CmdrMinimal provides the accessors to debug/trace flags
type CmdrMinimal interface {
	InDevMode() bool            // .devmode file existed? or devMode enabled?
	SetDevMode(devMode bool)    // set devMode manually
	IsDevModeFilePresent() bool // .devmode file existed and detected?

	InDebugging() bool      // is debug build
	GetDebugMode() bool     // is debug build or the debug-mode flag is true, settable by `--debug`
	SetDebugMode(b bool)    //
	GetDebugLevel() int     // return debug level as a integer, 0..n, it represents count of `--debug` or set by caller explicitly
	SetDebugLevel(hits int) //

	GetTraceMode() bool     //  the trace-mode flag, settable by `--trace`
	SetTraceMode(b bool)    //
	GetTraceLevel() int     // return trace level as a integer, 0..n, it represents count of `--trace` or set by caller explicitly
	SetTraceLevel(hits int) //

	IsNoColorMode() bool      // settable by `--no-color`
	SetNoColorMode(b bool)    //
	CountOfNoColor() int      //
	SetNoColorCount(hits int) //

	IsVerboseMode() bool      // settable by `--verbose` or `-v`
	IsVerboseModePure() bool  //
	SetVerboseMode(b bool)    //
	CountOfVerbose() int      //
	SetVerboseCount(hits int) //

	IsQuietMode() bool      // settable by `--quiet` or `-q`
	SetQuietMode(b bool)    //
	CountOfQuiet() int      //
	SetQuietCount(hits int) //

	SetOnDevModeChanged(funcs ...OnChanged)
	SetOnDebugChanged(funcs ...OnChanged)
	SetOnTraceChanged(funcs ...OnChanged)
	SetOnNoColorChanged(funcs ...OnChanged)
	SetOnVerboseChanged(funcs ...OnChanged)
	SetOnQuietChanged(funcs ...OnChanged)

	RemovcOnDevModeChanged(funcs ...OnChanged)
	RemovcOnDebugChanged(funcs ...OnChanged)
	RemoveOnTraceChanged(funcs ...OnChanged)
	RemoveOnNoColorChanged(funcs ...OnChanged)
	RemoveOnVerboseChanged(funcs ...OnChanged)
	RemoveOnQuietChanged(funcs ...OnChanged)
}

func Env() CmdrMinimal                  { return env }    // return minimal app env
func UpdateEnvWith(environ CmdrMinimal) { env = environ } // If u will, use ur own env

var env CmdrMinimal = &minimalEnv{}

// minimalEnv structure holds the debug/trace flags and provides CmdrMinimal accessors
type minimalEnv struct {
	debugMode          bool
	debugLevel         int
	traceMode          bool
	traceLevel         int
	noColorMode        bool
	noColorCount       int
	verboseMode        bool
	verboseCount       int
	quietMode          bool
	quietCount         int
	devMode            bool
	devModeFilePresent bool
	traceChanged       []OnChanged
	debugChanged       []OnChanged
	verboseChanged     []OnChanged
	quietChanged       []OnChanged
	noColorChanged     []OnChanged
	devModeChanged     []OnChanged
}

type OnChanged func(mod bool, level int)

func (e *minimalEnv) triggerDevModeChanged() {
	for _, cb := range e.devModeChanged {
		if cb != nil {
			cb(e.devMode, 1)
		}
	}
}

func (e *minimalEnv) triggerDebugChanged() {
	for _, cb := range e.debugChanged {
		if cb != nil {
			cb(e.debugMode, e.debugLevel)
		}
	}
}

func (e *minimalEnv) triggerTraceChanged() {
	for _, cb := range e.traceChanged {
		if cb != nil {
			cb(e.traceMode, e.traceLevel)
		}
	}
}

func (e *minimalEnv) triggerNoColorChanged() {
	for _, cb := range e.noColorChanged {
		if cb != nil {
			cb(e.noColorMode, e.noColorCount)
		}
	}
}

func (e *minimalEnv) triggerVerboseChanged() {
	for _, cb := range e.verboseChanged {
		if cb != nil {
			cb(e.verboseMode, e.verboseCount)
		}
	}
}

func (e *minimalEnv) triggerQuietChanged() {
	for _, cb := range e.quietChanged {
		if cb != nil {
			cb(e.quietMode, e.quietCount)
		}
	}
}

func (e *minimalEnv) SetOnDevModeChanged(funcs ...OnChanged) {
	e.debugChanged = append(e.devModeChanged, funcs...)
}

func (e *minimalEnv) SetOnDebugChanged(funcs ...OnChanged) {
	e.debugChanged = append(e.debugChanged, funcs...)
}

func (e *minimalEnv) SetOnTraceChanged(funcs ...OnChanged) {
	e.traceChanged = append(e.traceChanged, funcs...)
}

func (e *minimalEnv) SetOnNoColorChanged(funcs ...OnChanged) {
	e.noColorChanged = append(e.noColorChanged, funcs...)
}

func (e *minimalEnv) SetOnVerboseChanged(funcs ...OnChanged) {
	e.verboseChanged = append(e.verboseChanged, funcs...)
}

func (e *minimalEnv) SetOnQuietChanged(funcs ...OnChanged) {
	e.quietChanged = append(e.quietChanged, funcs...)
}

func (e *minimalEnv) RemovcOnDevModeChanged(funcs ...OnChanged) {
	e.debugChanged = e.removeHelper(e.devModeChanged, funcs...)
}

func (e *minimalEnv) RemovcOnDebugChanged(funcs ...OnChanged) {
	e.debugChanged = e.removeHelper(e.debugChanged, funcs...)
}

func (e *minimalEnv) RemoveOnTraceChanged(funcs ...OnChanged) {
	e.traceChanged = e.removeHelper(e.traceChanged, funcs...)
}

func (e *minimalEnv) RemoveOnNoColorChanged(funcs ...OnChanged) {
	e.noColorChanged = e.removeHelper(e.noColorChanged, funcs...)
}

func (e *minimalEnv) RemoveOnVerboseChanged(funcs ...OnChanged) {
	e.verboseChanged = e.removeHelper(e.verboseChanged, funcs...)
}

func (e *minimalEnv) RemoveOnQuietChanged(funcs ...OnChanged) {
	e.quietChanged = e.removeHelper(e.quietChanged, funcs...)
}

func (e *minimalEnv) removeHelper(slice []OnChanged, removingItems ...OnChanged) (result []OnChanged) {
	result = slice
	for _, h := range removingItems {
		if h != nil {
			hn := fnUniName(h)
			for i, cb := range slice {
				if fnUniName(cb) == hn {
					// result = slices.Delete(e.noColorChanged, i, 1)
					result = append(slice[0:i], slice[i+1:]...)
					break
				}
			}
		}
	}
	return
}

func fnUniName(f OnChanged) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func (e *minimalEnv) InDevMode() bool { return e.devMode }
func (e *minimalEnv) SetDevMode(devMode bool) {
	if save := e.devMode; save != devMode {
		e.devMode = devMode
		e.triggerDevModeChanged()
	}
}
func (e *minimalEnv) IsDevModeFilePresent() bool { return e.devModeFilePresent }

// InDebugging check if the delve debugger presents.
func (e *minimalEnv) InDebugging() bool { return isdelve.Enabled }

// GetDebugMode return the debug boolean flag generally.
func (e *minimalEnv) GetDebugMode() bool { return e.debugMode || isdelve.Enabled }

// SetDebugMode set the debug boolean flag generally.
func (e *minimalEnv) SetDebugMode(b bool) {
	if save := e.debugMode; save != b {
		e.debugMode = b
		e.triggerDebugChanged()
	}
}
func (e *minimalEnv) GetDebugLevel() int { return e.debugLevel }
func (e *minimalEnv) SetDebugLevel(hits int) {
	if save := e.debugLevel; save != hits {
		e.debugLevel = hits
		e.triggerDebugChanged()
	}
}

// GetTraceMode return the trace boolean flag generally.
func (e *minimalEnv) GetTraceMode() bool { return e.traceMode || trace.IsEnabled() }

// SetTraceMode set the trace boolean flag generally.
func (e *minimalEnv) SetTraceMode(b bool) {
	if save := e.traceMode; save != b {
		e.traceMode = b
		e.triggerTraceChanged()
	}
}
func (e *minimalEnv) GetTraceLevel() int { return e.traceLevel }
func (e *minimalEnv) SetTraceLevel(hits int) {
	if save := e.traceLevel; save != hits {
		e.traceLevel = hits
		e.triggerTraceChanged()
	}
}

func (e *minimalEnv) IsNoColorMode() bool { return e.noColorMode }
func (e *minimalEnv) SetNoColorMode(b bool) {
	if save := e.noColorMode; save != b {
		e.noColorMode = b
		e.triggerNoColorChanged()
	}
}

func (e *minimalEnv) CountOfNoColor() int { return e.noColorCount }
func (e *minimalEnv) SetNoColorCount(hits int) {
	if save := e.noColorCount; save != hits {
		e.noColorCount = hits
		e.triggerNoColorChanged()
	}
}

func (e *minimalEnv) IsVerboseMode() bool     { return buildtags.VerboseEnabled || e.verboseMode }
func (e *minimalEnv) IsVerboseModePure() bool { return e.verboseMode }
func (e *minimalEnv) SetVerboseMode(b bool) {
	if save := e.verboseMode; save != b {
		e.verboseMode = b
		e.triggerVerboseChanged()
	}
}

func (e *minimalEnv) CountOfVerbose() int { return e.verboseCount }
func (e *minimalEnv) SetVerboseCount(hits int) {
	if save := e.verboseCount; save != hits {
		e.verboseCount = hits
		e.triggerVerboseChanged()
	}
}

func (e *minimalEnv) IsQuietMode() bool { return e.quietMode }
func (e *minimalEnv) SetQuietMode(b bool) {
	if save := e.quietMode; save != b {
		e.quietMode = b
		e.triggerQuietChanged()
	}
}

func (e *minimalEnv) CountOfQuiet() int { return e.quietCount }
func (e *minimalEnv) SetQuietCount(hits int) {
	if save := e.quietCount; save != hits {
		e.quietCount = hits
		e.triggerQuietChanged()
	}
}
