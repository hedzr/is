# is


[![Go](https://github.com/hedzr/is/actions/workflows/go.yml/badge.svg)](https://github.com/hedzr/is/actions/workflows/go.yml)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/is.svg?label=release)](https://github.com/hedzr/is/releases)
[![go.dev](https://img.shields.io/badge/go-dev-green)](https://pkg.go.dev/github.com/hedzr/is)


`is` provides a set of detectors for the environment checking.

## Features

- `is.Env()` holds a global struct for CLI app basic states, such as: verbose/quiet/debug/trace....
- `is.InDebugging() bool`, `is.InTesting() bool`, and `is.InTracing() bool`, ....
- `is.DebugBuild() bool`
- `is.K8sBuild() bool`, `is.DockerBuild() bool`, ....
- `is.ColoredTty() bool`, ....
- Terminal Colorizer, Detector, unescape tool
- stringtool: `RandomStringPure`, ...
- basics: closable, closer, signals

To using environ detecting utilities better and smoother, some terminal (and stringtool, basics) tools are bundled together.

## Usages

```go
package main

import (
	"fmt"

	"github.com/hedzr/is"
	"github.com/hedzr/is/basics"
	"github.com/hedzr/is/term/color"
)

func main() {
	defer basics.Close()

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

```go
is.Terminal(os.Stdout)
```
The partials:

- is.InDebugging / InDebugMode
- is.DebuggerAttached (relyes on delve tag)
- is.InTracing / InTestingT
- is.InTesting / InTestingT
- is.InDevelopingTime
- is.InVscodeTerminal
- is.InK8s
- is.InIstio
- is.InDocker / InDockerEnvSimple
- Build
  - is.K8sBuild
  - is.IstioBuild
  - is.DockerBuild
  - is.VerboseBuild
  - is.DebugBuild
  - buildtags.IsBuildTagExists

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
  - is.Terminal(file)
  - is.TerminalFd(fd)
  - is.Tty(wr)
  - is.ColoredTty(wr)
  - is.AnsiEscaped(s) (~~IsTtyEscaped(s)~~)
  - StripEscapes(s)
  - ReadPassword
  - GetTtySize
  - is.GetTtySizeByName(filename) (cols,rows,err)
  - is.GetTtySizeByFile(file) (cols,rows,err)
  - is.GetTtySizeByFd(fd) (cols,rows,err)
  - StartupByDoubleClick() bool

- [Special] Terminal / Color
  - escaping tools: GetCPT()/GetCPTC()/GetCPTNC()
  - Highlight, Dimf, Text, Dim, ToDim, ToHighlight, ToColor, ...

- Basics
  - closers 
    - Peripheral, Closable, Closer
    - RegisterClosable
    - RegisterClosers
    - RegisterCloseFns
  - `is.Signals().Catcher()`
  - is.FileExists(filepath)
  - is.ToBool, StringToBool

### Build Tags

Some functions want special buildtags presented. These are including:

- `verbose`: See VerboseBuild, ...
- `delve`: See DebugBuild, ...
- `k8s`: See K8sBuild
- `istio`: See IstioBuild
- `docker`: See DockerBuild
- ...
- `buildtags.IsBuildTagExists(tag) bool`

### Colorizes

The test codes:

```go
import "github.com/hedzr/is/term/color"

func TestGetCPT(t *testing.T) {
t.Logf("%v", color.GetCPT().Translate(`<code>code</code> | <kbd>CTRL</kbd>
	<b>bold / strong / em</b>
	<i>italic / cite</i>
	<u>underline</u>
	<mark>inverse mark</mark>
	<del>strike / del </del>
	<font color="green">green text</font>
    `, color.FgDefault))
}
```

Result:

![image-20231107100150520](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20231107100150520.png)

And more:

```go
func TestStripLeftTabs(t *testing.T) {
t.Logf("%v", color.StripLeftTabs(`
	
		<code>code</code>
	NC Cool
	 But it's tight.
	  Hold On!
	Hurry Up.
    `))
}

func TestStripHTMLTags(t *testing.T) {
t.Logf("%v", color.StripHTMLTags(`
	
		<code>code</code>
	NC Cool
	 But it's tight.
	  Hold On!
	Hurry Up.
    `))
}
```

### Basics (Closers)

The `Closers()` collects all closable objects and allow shutting down them at once.

```go
package main

import (
	"os"

	"github.com/hedzr/is/basics"
)

type redisHub struct{}

func (s *redisHub) Close() {
    // close the connections to redis servers
    println("redis connections closed")
}

func main() {
    defer basics.Close()

    tmpFile, _ := os.CreateTemp(os.TempDir(), "1*.log")
    basics.RegisterClosers(tmpFile)

    basics.RegisterCloseFn(func() {
        // do some shutdown operations here
        println("close single functor")
    })

    basics.RegisterPeripheral(&redisHub{})
}
```

### Basics (Signals)

`Signals()` could catch os signals and entering a infinite loop.

For example, a tcp server could be:

```go
package main

import (
	"context"
	"os"
	"sync"

	"github.com/hedzr/go-socketlib/net"
	"github.com/hedzr/is"
	logz "github.com/hedzr/logg/slog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logz.New("new-dns")
	server := net.NewServer(":7099",
		net.WithServerOnListening(func(ss net.Server, l stdnet.Listener) {
			go runClient(ctx, ss, l)
		}),
		net.WithServerLogger(logger.WithSkip(1)),
	)
	defer server.Close()

	catcher := is.Signals().Catch()
	catcher.
		// WithVerboseFn(func(msg string, args ...any) {
		// 	logz.WithSkip(2).Verbose(fmt.Sprintf("[verbose] %s", fmt.Sprintf(msg, args...)))
		// }).
		WithOnSignalCaught(func(sig os.Signal, wg *sync.WaitGroup) {
			println()
			logz.Debug("signal caught", "sig", sig)
			if err := server.Shutdown(); err != nil {
				logz.Error("server shutdown error", "err", err)
			}
			cancel()
		}).
		Wait(func(stopChan chan<- os.Signal, wgShutdown *sync.WaitGroup) {
			logz.Debug("entering looper's loop...")

			server.WithOnShutdown(func(err error, ss net.Server) { wgShutdown.Done() })
			err := server.ListenAndServe(ctx, nil)
			if err != nil {
				logz.Fatal("server serve failed", "err", err)
			}
		})
}

func runClient(ctx context.Context, ss net.Server, l stdnet.Listener) {
	c := net.NewClient()

	if err := c.Dial("tcp", ":7099"); err != nil {
		logz.Fatal("connecting to server failed", "err", err, "server-endpoint", ":7099")
	}
	logz.Info("[client] connected", "server.addr", c.RemoteAddr())
	c.RunDemo(ctx)
}
```

> some packages has stayed in progress so the above codes is just a skeleton (from go-socketlib/_examples/new-loop/main.go/v1).




## Contributions

Kindly welcome, please issue me first for keeping this repo smaller.

## License

Apache 2.0