package is

import (
	"os"
	"runtime"
	"strings"

	"os/exec"

	"github.com/hedzr/is/dirs"
)

// Root returns true if current user is 'root' or user is in sudo mode.
//
// For windows it's always false.
func Root() bool {
	return dirs.Unix() && os.Getuid() == 0
}

// AMD64 returns true if CPU arch is amd64.
func AMD64() bool {
	return runtime.GOARCH == "amd64"
}

// AMD64 returns true if CPU arch is amd64.
func I386() bool {
	return runtime.GOARCH == "386"
}

// AMD32 returns true if CPU arch is arm32.
func ARM32() bool {
	return runtime.GOARCH == "arm"
}

// AMD32BE returns true if CPU arch is arm32be.
func ARM32BE() bool {
	return runtime.GOARCH == "armbe"
}

// ARM64 returns true if CPU arch is arm64.
func ARM64() bool {
	return runtime.GOARCH == "arm64"
}

// ARM64BE returns true if CPU arch is arm64be.
func ARM64BE() bool {
	return runtime.GOARCH == "arm64be"
}

// Windows returns true for Microsoft Windows Platform.
func Windows() bool {
	return runtime.GOOS == "windows"
}

// WindowsWSL return true if running under Windows WSL env.
func WindowsWSL() bool {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "windows_standard")
}

// Unix returns true for Linux, Darwin, and Others Unix-like Platforms.
//
// [NOTE] for unix platforms, more assetions and tests needed.
func Unix() bool {
	return dirs.Unix()
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

// BSD returns true if running under Any BSD platform, including FreeBSD, NetBSD and OpenBSD.
func BSD() bool {
	return runtime.GOOS == "freebsd" || runtime.GOOS == "netbsd" || runtime.GOOS == "openbsd"
}

// Aix returns true if running under Aix platform.
func Aix() bool {
	return runtime.GOOS == "aix"
}

// Android returns true if running under Android platform.
func Android() bool {
	return runtime.GOOS == "android"
}

// Dragonfly returns true if running under Dragonfly platform.
func Dragonfly() bool {
	return runtime.GOOS == "dragonfly"
}

// FreeBSD returns true if running under FreeBSD platform.
func FreeBSD() bool {
	return runtime.GOOS == "freebsd"
}

// Hurd returns true if running under Hurd platform.
func Hurd() bool {
	return runtime.GOOS == "hurd"
}

// Illumos returns true if running under Illumos platform.
func Illumos() bool {
	return runtime.GOOS == "illumos"
}

// IOS returns true if running under iOS platform.
func IOS() bool {
	return runtime.GOOS == "ios"
}

// JS returns true if running under JS/WASM platform.
func JS() bool {
	return runtime.GOOS == "js"
}

// Nacl returns true if running under Nacl platform.
func Nacl() bool {
	return runtime.GOOS == "nacl"
}

// NetBSD returns true if running under NetBSD platform.
func NetBSD() bool {
	return runtime.GOOS == "netbsd"
}

// OpenBSD returns true if running under OpenBSD platform.
func OpenBSD() bool {
	return runtime.GOOS == "openbsd"
}

// Plan9 returns true if running under Plan9 platform.
func Plan9() bool {
	return runtime.GOOS == "plan9"
}

// Solaris returns true if running under Solaris platform.
func Solaris() bool {
	return runtime.GOOS == `solaris`
}

// Wasip1 returns true if running under Wasip1 platform.
func Wasip1() bool {
	return runtime.GOOS == `wasip1`
}

// Zos returns true if running under Zos platform.
func Zos() bool {
	return runtime.GOOS == `zos`
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
