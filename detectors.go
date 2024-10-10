package is

import (
	"os"
	"runtime"
	"strings"

	"github.com/hedzr/env/exec"
)

// Root returns true if current user is 'root' or user is in sudo mode.
//
// For windows it's always false.
func Root() bool {
	return Unix() && os.Getuid() == 0
}

// Windows returns true for Microsoft Windows Platform.
func Windows() bool {
	return runtime.GOOS == "windows"
}

// WindowsWSL return true if running under Windows WSL env.
func WindowsWSL() bool {
	_, txt, err := exec.RunWithOutput("uname", "-r")
	return err == nil && strings.Contains(txt, "windows_standard")
}

// Unix returns true for Linux, Darwin, and Others Unix-like Platforms.
func Unix() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin" || runtime.GOOS == "unix"
}

// Linux return true for General Linux Distros.
func Linux() bool {
	return runtime.GOOS == "linux"
}

// Darwin returns true if running under macOS platform, including both Intel and Silicon.
func Darwin() bool {
	return runtime.GOOS == "darwin"
}

// DarwinSilicon returns true if running under Apple Silicon.
func DarwinSilicon() bool {
	return runtime.GOOS == "darwin" && runtime.GOARCH == "arm64"
}

// DarwinIntel returns true if running under Apple Intel Machines.
func DarwinIntel() bool {
	return runtime.GOOS == "darwin" && runtime.GOARCH == "amd64"
}

// Bash returns true if application is running under a Bash shell.
func Bash() bool {
	return os.Getenv("BASH_VERSION") != "" || os.Getenv("BASH") != ""
}

// Zsh returns true if application is running under a Zsh shell.
func Zsh() bool {
	return strings.Contains(ShellName(), "/bin/zsh") && os.Getenv("ZSH_NAME") != ""
}

// Fish returns true if application is running under a Fish shell.
func Fish() bool {
	return os.Getenv("FISH_VERSION") != "" && strings.Contains(ShellName(), "/bin/fish")
}

// Powershell returns true if application is running under a Windows Powershell shell.
//
// Not testing yet.
func Powershell() bool {
	return os.Getenv("PS1") != ""
}

// ShellName returns current SHELL's name.
//
// For Windows, it could be "cmd.exe" or "powershell.exe".
// For Linux or Unix or Darwin, it returns environment variable $SHELL.
// Else it's empty "".
func ShellName() string {
	switch runtime.GOOS {
	case "windows":
		if Powershell() {
			return "powershell.exe"
		}
		return "cmd.exe"
	case "linux", "darwin":
		return os.Getenv("SHELL")
	default:
		return ""
	}
}
