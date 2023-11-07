# is


[![Go](https://github.com/hedzr/is/actions/workflows/go.yml/badge.svg)](https://github.com/hedzr/is/actions/workflows/go.yml)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/is.svg?label=release)](https://github.com/hedzr/is/releases)
[![go.dev](https://img.shields.io/badge/go-dev-green)](https://pkg.go.dev/github.com/hedzr/is)


`is` provides a set of detectors for the environment checking.

<kbd>pre-release</kbd>

## Features

- `Env()` holds a global struct for CLI app basic states, such as: verbose/quiet/debug/trace....
- `InDebugging()`, `InTesting()`, and `InTracing()`, ....
- `DebugBuild()`
- `K8sBuild()`, `DockerBuild()`, ....
- `IsColoredTty()`, ....
- Terminal colorizers
- stringtool: `RandomStringPure`

To using environ detecting utilities better and smoother, some terminal (and stringtool) tools are bundle together.

## Usages

```go
package main

import (
	"fmt"

	"github.com/hedzr/is"
	"github.com/hedzr/is/term/color"
)

func main() {
	println(is.InTesting())
	println(is.Env().GetDebugLevel())
	fmt.Printf("%v", color.GetCPT().Translate(`<code>code</code> | <kbd>CTRL</kbd>
	<b>bold / strong / em</b>
	<i>italic / cite</i>
	<u>underline</u>
	<mark>inverse mark</mark>
	<del>strike / del </del>
	<font color="green">green text</font>
	`, color.FgDefault))
}
```

Result got:

![image-20231107101843524](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20231107101843524.png)

### Lists

The partials:

- InDebugging / InDebugMode
- DebuggerAttached (relyes on delve tag)
- InTracing / InTestingT
- InTesting / InTestingT
- InDevelopingTime
- InVscodeTerminal
- InK8s
- InIstio
- InDocker / InDockerEnvSimple
- Build
  - K8sBuild
  - IstioBuild
  - DockerBuild
  - VerboseBuild
  - DebugBuild

- States / Env
  - VerboseModeEnabled
  - GetVerboseLevel / SetVerboseMode / SetVerboseLevel
  - QuietModeEnabled
  - GetQuietLevel / SetQuietMode / SetQuietLevel
  - NoColorMode
  - GetNoColorLevel / SetNoColorMode / SetNoColorLevel
  - DebugMode
  - GetDebugLevel / SetDebugMode / SetDebugLevel
  - Tracing
  - TraceMode
  - GetTraceLevel / SetTraceMode / SetTraceLevel

- Terminal / Tty
  - IsTty(w)
  - IsColoredTty(w)
  - IsTtyEscaped(s)
  - StripEscapes(s)
  - ReadPassword
  - GetTtySize


### Buildtags

Some functions want special buildtags presented. These are including:

- verbose: See VerboseBuild, ...
- delve: See DebugBuild, ...
- k8s: See K8sBuild
- istio: See IstioBuild
- docker: See DockerBuild
- ...



### Colorizers

The test codes:

```go
func TestGetCPT(t *testing.T) {
	t.Logf("%v", GetCPT().Translate(`<code>code</code> | <kbd>CTRL</kbd>
	<b>bold / strong / em</b>
	<i>italic / cite</i>
	<u>underline</u>
	<mark>inverse mark</mark>
	<del>strike / del </del>
	<font color="green">green text</font>
	`, FgDefault))
}
```

Result:

![image-20231107100150520](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20231107100150520.png)

And more:

```go
func TestStripLeftTabs(t *testing.T) {
	t.Logf("%v", StripLeftTabs(`
	
		<code>code</code>
	NC Cool
	 But it's tight.
	  Hold On!
	Hurry Up.
	`))
}

func TestStripHTMLTags(t *testing.T) {
	t.Logf("%v", StripHTMLTags(`
	
		<code>code</code>
	NC Cool
	 But it's tight.
	  Hold On!
	Hurry Up.
	`))
}

```







## Contributions

Kindly welcome, please issue me first for keep this repo smaller.

## License

Apache 2.0