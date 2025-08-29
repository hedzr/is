# is

[![Go](https://github.com/hedzr/is/actions/workflows/go.yml/badge.svg)](https://github.com/hedzr/is/actions/workflows/go.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/hedzr/is)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/is.svg?label=release)](https://github.com/hedzr/is/releases)
[![go.dev](https://img.shields.io/badge/go-dev-green)](https://pkg.go.dev/github.com/hedzr/is)
[![deps.dev](https://img.shields.io/badge/deps-dev-green)](https://deps.dev/go/github.com%2Fhedzr%2Fis)

`is` provides numerous detectors for checking the states of environment (build, executive, ...).

## Features

- `is.State(which) bool`: the universal detector entry - via `RegisterStateGetter(state string, getter func() bool)` to add your own ones. *Since v0.5.11*
- `is.Env()` holds a global struct for CLI app basic states, such as: verbose/quiet/debug/trace....
  - `DebugMode`/`DebugLevel`, `TraceMode`/`TraceLevel`, `ColorMode`, ...
- `is.InDebugging() bool`, `is.InTesting() bool`, and `is.InTracing() bool`, ....
- `is.DebugBuild() bool`.
- `is.K8sBuild() bool`, `is.DockerBuild() bool`, ....
- `is.ColoredTty() bool`, ....
- `is.Color()` to get an indexer for the functions in our term/color subpackage, ...
- Terminal Colorizer, Detector, unescape tools.
  - `is/term/color.Color` interface
- stringtool: `RandomStringPure`, case-converters ...
- basics: closable, closer, signals.
  - easier `Press any key to exit...` prompt: `is.Signals().Catch()`
- exec: Run, RunWithOutput, Sudo, ...
- ~~go1.23.7+ required since v0.7.0~~
- ~~go 1.24.0+ required~~
- go1.24.5 required since v0.8.55

See the above badge to get the exact required go toolchain version.

To using environment detecting utilities better and smoother, some terminal (and stringtool, basics) tools are bundled together.

Since v0.6.0, `is.InDebugging()` checks if the running process' parent is `dlv`.
The old `DebugMode` and `DebugBuild` are still work:

- `InDebugging`: checks this process is being debugged by `dlv`.
- `DebugBuild`: `-tags=delve` is set at building.
- `DebugMode`: `--debug` is specified at command line.

Since v0.8.27, `basics.Signals().Catcher().WaitFor()` wants `ctx` param passed in.

## Usages

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "sync"
    "time"

    "github.com/hedzr/is"
    "github.com/hedzr/is/basics"
    "github.com/hedzr/is/term/color"
)

func main() {
    // defer basics.Close() // uncomment if not using Catcher.WaitFor and/or cmdr.v2

    is.RegisterStateGetter("custom", func() bool { return is.InVscodeTerminal() })

    println(is.InTesting())
    println(is.State("in-testing"))
    println(is.State("custom")) // detects a state with custom detector
    println(is.Env().GetDebugLevel())
    if is.InDebugMode() {
        slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))
    }

    // or:
    //    is.Color().GetColorTranslator().Translate("<b>bold</b>")
    fmt.Printf("%v", color.GetCPT().Translate(`<code>code</code> | <kbd>CTRL</kbd>
        <b>bold / strong / em</b>
        <i>italic / cite</i>
        <u>underline</u>
        <mark>inverse mark</mark>
        <del>strike / del </del>
        <font color="green">green text</font>
`, color.FgDefault))

    var cancelled int32
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    catcher := is.Signals().Catch()
    catcher.
        WithPrompt("Press CTRL-C to quit...").
        // WithOnLoopFunc(dbStarter, cacheStarter, mqStarter).
        WithPeripherals(&dbMgr{}).
        WithOnSignalCaught(func(ctx context.Context, sig os.Signal, wg *sync.WaitGroup) {
            println()
            slog.Info("signal caught", "sig", sig)
            cancel() // cancel user's loop, see Wait(...)
        }).
        WaitFor(ctx, func(ctx context.Context, closer func()) {
            slog.Debug("entering looper's loop...")
            defer close() // notify catcher we want to shutdown
            // to terminate this app after a while automatically:
            time.Sleep(10 * time.Second)

            if atomic.CompareAndSwapInt32(&cancelled, 0, 1) {
                is.PressAnyKeyToContinue(os.Stdin)
            }
        })
}

type dbMgr struct{}

func (*dbMgr) Close()                           {}         // before app terminatine
func (*dbMgr) Open(context.Context) (err error) { return } // ran before catcher.WaitFor()
```

Result is similar with:

![image-20240113071930661](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20240113071930661.png)

NOTE that `is.Signals().Catch()` will produce a prompt and enter a infinite loop to wait for user's keystroke pressed.

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
  - `color.Color` interface
  - `color.New()` return a stream-callable color object: `color.Cursor`.

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

### `Cursor`

Since v0.8+, A new `color.Cursor` object can be initialized by `color.New()`, which support format the colorful text with streaming calls, for console/tty.

> See [the online docs](https://docs.hedzr.com/docs/is/) for more usages.

<details title="Expand"><caption>expand</caption><detail>

The examples are:

```go

func ExampleNew() {
 // start a color text builder
 var c = color.New()

 // specially for running on remote ci server
 if states.Env().IsNoColorMode() {
  states.Env().SetNoColorMode(true)
 }

 // paint and get the result (with ansi-color-seq ready)
 var result = c.Println().
  Color16(color.FgRed).
  Printf("hello, %s.", "world").Println().
  SavePos().
  Println("x").
  Color16(color.FgGreen).Printf("hello, %s.\n", "world").
  Color256(160).Printf("[160] hello, %s.\n", "world").
  Color256(161).Printf("[161] hello, %s.\n", "world").
  Color256(162).Printf("[162] hello, %s.\n", "world").
  Color256(163).Printf("[163] hello, %s.\n", "world").
  Color256(164).Printf("[164] hello, %s.\n", "world").
  Color256(165).Printf("[165] hello, %s.\n", "world").
  Up(3).Echo(" ERASED ").
  RGB(211, 211, 33).Printf("[16m] hello, %s.", "world").
  Println().
  RestorePos().
  Println("z").
  Down(8).
  Println("DONE").
  Build()

  // and render the result
 fmt.Println(result)

 // For most of ttys, the output looks like:
 //
 // [31mhello, world.[0m
 // [sx
 // [32mhello, world.
 // [38;5;160m[160] hello, world.
 // [38;5;161m[161] hello, world.
 // [38;5;162m[162] hello, world.
 // [38;5;163m[163] hello, world.
 // [38;5;164m[164] hello, world.
 // [38;5;165m[165] hello, world.
 // [0m[3A ERASED [38;2;211;211;33m[16m] hello, world.
 // [uz
 // [8BDONE
}

func ExampleCursor_Color16() {
 // another colorful builfer
 var c = color.New()
 fmt.Println(c.Color16(color.FgRed).
  Printf("hello, %s.", "world").Println().Build())
 // Output: [31mhello, world.[0m
}

func ExampleCursor_Color() {
 // another colorful builfer
 var c = color.New()
 fmt.Println(c.Color(color.FgRed, "hello, %s.", "world").Build())
 // Output: [31mhello, world.[0m
}

func ExampleCursor_Bg() {
 // another colorful builfer
 var c = color.New()
 fmt.Println(c.Bg(color.BgRed, "hello, %s.", "world").Build())
 // Output: [41mhello, world.[0m
}

func ExampleCursor_Effect() {
 // another colorful builfer
 var c = color.New()
 fmt.Println(c.Effect(color.BgDim, "hello, %s.", "world").Build())
 // Output: [2mhello, world.[0m
}

func ExampleCursor_Color256() {
 // another colorful builfer
 var c = color.New()
 fmt.Print(c.
  Color256(163).Printf("[163] hello, %s.\n", "world").
  Color256(164).Printf("[164] hello, %s.\n", "world").
  Color256(165).Printf("[165] hello, %s.\n", "world").
  Build())
 // Output:
 // [38;5;163m[163] hello, world.
 // [38;5;164m[164] hello, world.
 // [38;5;165m[165] hello, world.
}

func ExampleCursor_RGB() {
 // another colorful builfer
 var c = color.New()
 fmt.Print(c.
  RGB(211, 211, 33).Printf("[16m] hello, %s.\n", "world").
  BgRGB(211, 211, 33).Printf("[16m] hello, %s.\n", "world").
  Build())
 // Output:
 // [38;2;211;211;33m[16m] hello, world.
 // [48;2;211;211;33m[16m] hello, world.
}

func ExampleCursor_EDim() {
 // another colorful builfer
 var c = color.New()
 fmt.Print(c. // Color16(color.FgRed).
   EDim("[DIM] hello, %s.\n", "world").String())
 // Output:
 // [2m[DIM] hello, world.
 // [0m
}

func ExampleCursor_Black() {
 // another colorful builfer
 var c = color.New()
 fmt.Print(c. // Color16(color.FgRed).
   Black("[BLACK] hello, %s.\n", "world").String())
 // Output:
 // [30m[BLACK] hello, world.
 // [0m
}

func ExampleCursor_BgBlack() {
 // another colorful builfer
 var c = color.New()
 fmt.Print(c. // Color16(color.FgRed).
   BgBlack("[BGBLACK] hello, %s.\n", "world").String())
 // Output:
 // [40m[BGBLACK] hello, world.
 // [0m
}

func ExampleCursor_Translate() {
 // another colorful builfer
 var c = color.New()
 fmt.Print(c. // Color16(color.FgRed).
   Translate(`<code>code</code> | <kbd>CTRL</kbd>
  <b>bold / strong / em</b>
  <i>italic / cite</i>
  <u>underline</u>
  <mark>inverse mark</mark>
  <del>strike / del </del>
  <font color="green">green text</font>
  `).String())
 // Output:
 // [51;1mcode[0m[39m | [51;1mCTRL[0m[39m
 //  [1mbold / strong / em[0m[39m
 //  [3mitalic / cite[0m[39m
 //  [4munderline[0m[39m
 //  [7minverse mark[0m[39m
 //  [9mstrike / del [0m[39m
 //  [32mgreen text[0m[39m
}

func ExampleCursor_StripLeftTabsColorful() {
 // another colorful builfer
 var c = color.New()
 fmt.Print(c. // Color16(color.FgRed).
   StripLeftTabsColorful(`
  <code>code</code> | <kbd>CTRL</kbd>
  <b>bold / strong / em</b>
  <i>italic / cite</i>
  <u>underline</u>
  <mark>inverse mark</mark>
  <del>strike / del </del>
  <font color="green">green text</font>
  `).String())
 // Output:
 // [51;1mcode[0m[0m | [51;1mCTRL[0m[0m
 // [1mbold / strong / em[0m[0m
 // [3mitalic / cite[0m[0m
 // [4munderline[0m[0m
 // [7minverse mark[0m[0m
 // [9mstrike / del [0m[0m
 // [32mgreen text[0m[0m
}
```

</detail></details>

### `color` subpackage

Package color provides a wrapped standard output device like printf but with colored enhancements.

The main types are [Cursor](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#Cursor) and [Translator](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#Translator).

[Cursor](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#Cursor) allows formatting colorful text and moving cursor to another coordinate.

[New](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#New) will return a [Cursor](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#Cursor) object.

[RowsBlock](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#RowsBlock) is another cursor controller, which can treat the current line and following lines as a block and updating these lines repeatedly. This feature will help the progressbar writers or the continuous lines updater.

[Translator](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#Translator) is a text and tiny HTML tags translator to convert these markup text into colorful console text sequences. [GetCPT](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#GetCPT) can return a smart translator which translate colorful text or strip the ansi escaped sequence from result text if `states.Env().IsNoColorMode()` is true.

[Color](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#Color) is an interface type to represent a terminal color object, which can be serialized to ansi escaped sequence directly by [Color.Color].

To create a [Color](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#Color) object, there are several ways:

- by [NewColor16](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#NewColor16), or use [Color16](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#Color16) constants directly like [FgBlack](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#FgBlack), [BgGreen](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#BgGreen), ...
- by [NewColor256](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#NewColor256) to make a 8-bit 256-colors object
- by [NewColor16m](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#NewColor16m) to make a true-color object
- by [NewControlCode](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#NewControlCode) or [ControlCode](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#ControlCode) constants
- by [NewFeCode](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#NewFeCode) or [FeCode](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#FeCode) constants
- by [NewSGR](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#NewSGR) or use [CSIsgr](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#CSIsgr) constants directly like [SGRdim](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#SGRdim), [SGRstrike](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#SGRstrike), ...
- by [NewStyle](https://pkg.go.dev/github.com/hedzr/is@v0.8.31/term/color#NewStyle) to make a compounded object
- ...

As to v0.8.53, a `Color` object can wrap itself arroubd the given text:

```go
color.BgBold.Wrap(color.FgRed.Wrap("ERROR!"))
color.BgDim.Wrap(color.FgDarkGray.Wrap("debug message here."))
```

## Integrated with `cmdr`

### `Closers`

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

### `Signals`

`Signals()` could catch OS signals and entering a infinite loop.

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

    // make a catcher so that it can catch ths signals,
    catcher := is.Signals().Catch()
    catcher.
        // WithVerboseFn(func(msg string, args ...any) {
        //     logz.WithSkip(2).Verbose(fmt.Sprintf("[verbose] %s", fmt.Sprintf(msg, args...)))
        // }).
        WithOnSignalCaught(func(ctx context.Context, sig os.Signal, wg *sync.WaitGroup) {
            println()
            logz.Debug("signal caught", "sig", sig)
            if err := server.Shutdown(); err != nil {
                logz.Error("server shutdown error", "err", err)
            }
            cancel()
        }).
        WaitFor(ctx, func(ctx context.Context, closer func()) {
            logz.Debug("entering looper's loop...")

            server.WithOnShutdown(func(err error, ss net.Server) { closer() })
            err := server.ListenAndServe(ctx, nil)
            if err != nil {
                logz.Fatal("server serve failed", "err", err)
            } else {
                closer()
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

## Contributions

Kindly welcome, please issue me first for keeping this repo smaller.

## License

under Apache 2.0
