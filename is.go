// Package is provides environ detectors, and other basis.
//
// ## Environment Detectors
//
// fo runtime modes (debug, trace, verbose, quiet, no-color)
//
//   - `is.State(which) bool`: the universal detector entry
//
//   - via `RegisterStateGetter(state string, getter func() bool)` to add your own ones. *Since v0.5.11*
//
//   - `is.Env()` holds a global struct for CLI app basic states, such as: verbose/quiet/debug/trace....
//
//   - `DebugMode`/`DebugLevel`, `TraceMode`/`TraceLevel`, `ColorMode`, ...
//
//   - `is.InDebugging() bool`, `is.InTesting() bool`, and `is.InTracing() bool`, ....
//
//   - `is.DebugBuild() bool`.
//
//   - `is.K8sBuild() bool`, `is.DockerBuild() bool`, ....
//
//   - `is.ColoredTty() bool`, ....
//
//   - `is.Color()` to get an indexer for the functions in our term/color subpackage, ...
//
// os, shell,
//
//   - `is.Zsh()`, `is.Bash()`
//
// `exec` subpackage: for launching another program, shell command, ...
//
//   - starts from `exec.New()`
//
// `dir`: dir, file operations, dir/file existance detector
//
// `term`: terminal env detectors, operations.
//
// `term/color`: ansi escaped sequences wrapper, simple html tags
// tranlator (for color tags).
//
//   - starts from `color.New()`
//   - `RowsBlock` by `color.NewRowsBlock()`
package is
