package states

// IsUnderDebugger detects if golang debugger 'dlv' is
// controlling this process.
//
// Only for Linux, Darwin and Windows, need more completely test.
func IsUnderDebugger() bool { return isDebuggerAttached() }

// Process interface.
//
// Part of findProcess, there's no plan to implement it.
type Process interface {
	// Pid is the process ID for this process.
	Pid() int

	// PPid is the parent process ID for this process.
	PPid() int

	// Executable name running this process. This is not a path to the
	// executable.
	Executable() string
}
