package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hedzr/is"
	"github.com/hedzr/is/basics"
	"github.com/hedzr/is/term"
	"github.com/hedzr/is/term/color"
)

var colortableCmd, sgrCmd, cptCmd *flag.FlagSet

func init() {
	colortableCmd = flag.NewFlagSet("colortable", flag.ExitOnError)
	sgrCmd = flag.NewFlagSet("sgr", flag.ExitOnError)
	cptCmd = flag.NewFlagSet("cpt", flag.ExitOnError)
}

func main() {
	flag.Parse()

	if len(os.Args) <= 1 {
		flag.Usage()
		fmt.Printf("\n%s", color.StripLeftTabsC(`

		Subcommands:
		  ct, color-table     print 256-colors table
		  sgr, effect.        print SGR effects (eg, bold, underline, ...)
		  cpt                 Translator demo

		`))
		return
	}

	switch os.Args[1] {
	case "color-table", "ct":
		colortableCmd.Parse(os.Args[2:])
		runColor256table()
	case "sgr", "16-colors", "4-bit", "effect":
		sgrCmd.Parse(os.Args[2:])
		runSGReffects()
	case "cpt", "translator":
		cptCmd.Parse(os.Args[2:])
		runCPT()
	default:
		log.Fatalf("[ERROR] unknown subcommand '%s', see help for more details.", os.Args[1])
	}
}

func runColor256table() {
	color.Color256table(os.Stdout)
}

func runSGReffects() {
	for i, sgrs := range []struct {
		pre, post color.CSIsgr
		desc      string
	}{
		{color.SGRbold, color.SGRresetBoldAndDim, "bold"},
		{color.SGRdim, color.SGRresetBoldAndDim, "dim"},
		{color.SGRitalic, color.SGRresetItalic, "italic"},
		{color.SGRunderline, color.SGRresetUnderline, "underline"},
		{color.SGRslowblink, color.SGRresetSlowBlink, "blink"},
		{color.SGRrapidblink, color.SGRresetRapidBlink, "fast blink"},
		{color.SGRinverse, color.SGRresetInverse, "inverse"},
		{color.SGRhide, color.SGRresetHide, "hide"},
		{color.SGRstrike, color.SGRresetStrike, "strike"},
		{color.SGRframed, color.SGRneitherFramedNorEncircled, "framed"},
		{color.SGRencircled, color.SGRneitherFramedNorEncircled, "encircled"},
		{color.SGRoverlined, color.SGRnotoverlined, "overlined"},
		{color.SGRideogramUnderline, color.SGRresetIdeogram, "ideogram underline"},
		{color.SGRideogramDoubleUnderline, color.SGRresetIdeogram, "ideogram double underline"},
		{color.SGRideogramOverline, color.SGRresetIdeogram, "ideogram overline"},
		{color.SGRideogramDoubleOverline, color.SGRresetIdeogram, "ideogram double overline"},
		{color.SGRideogramStressMarking, color.SGRresetIdeogram, "ideogram stress marking"},
		{color.SGRsuperscript, color.SGRresetSuperscriptAndSubscript, "superscript"},
		{color.SGRsubscript, color.SGRresetSuperscriptAndSubscript, "subscript"},
		// {color.SGRdim, color.SGRresetDim},
	} {
		str := fmt.Sprintf(`%5d. %s%s%s %s`,
			i, sgrs.pre,
			"Hello, World!",
			sgrs.post,
			sgrs.desc,
		)
		fmt.Println(str)
	}

	fmt.Println(color.SGRreset)

	// SGRsetFg
	fmt.Printf("\x1b[%d;5;9m[ 9 TEST string HERE]%s\n",
		color.SGRsetFg,
		color.SGRdefaultFg,
	)

	fmt.Printf("\x1b[%d;5;21m[21 TEST string HERE]%s\n",
		color.SGRsetFg,
		color.SGRdefaultFg,
	)

	fmt.Println(color.SGRreset)

}

func dbStarter(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
	// initializing database connections...
	// ...
	wgDone.Done()
}

func cacheStarter(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
	// initializing redis cache connections...
	// ...
	wgDone.Done()
}

func mqStarter(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
	// initializing message queue connections...
	// ...
	wgDone.Done()
}

type dbMgr struct{}

func (s *dbMgr) CloseByCatcher(closer func()) {
	s.Close()
	closer()
}
func (s *dbMgr) Close() {
	// do shutdown stuffs...
	log.Printf("dbMgr closed.\n")
}

func runCPT() {
	defer basics.Close()

	is.RegisterStateGetter("custom", func() bool { return is.InVscodeTerminal() })

	println("state.InTesting:   ", is.InTesting())
	println("state.in-testing:  ", is.State("in-testing"))
	println("state.custom:      ", is.State("custom")) // detects a state with custom detector
	println("env.GetDebugLevel: ", is.Env().GetDebugLevel())
	if is.InDebugMode() {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))
	}

	fmt.Printf("\n%v", color.GetCPT().Translate(`<code>code</code> | <kbd>CTRL</kbd>
		<b>bold / strong / em</b>
		<i>italic / cite</i>
		<u>underline</u>
		<mark>inverse mark</mark>
		<del>strike / del </del>
		<font color="green">green text</font>
`, color.FgDefault))

	println("term.IsTerminal:               ", term.IsTerminal(int(os.Stdout.Fd())))
	println("term.IsAnsiEscaped:            ", term.IsAnsiEscaped(color.GetCPT().Translate(`<code>code</code>`, color.FgDefault)))
	println("term.IsCharDevice(stdout):     ", term.IsCharDevice(os.Stdout))
	rows, cols, err := term.GetFdSize(os.Stdout.Fd())
	println("term.GetFdSize(stdout):        ", rows, cols, err)
	rows, cols, err = term.GetTtySizeByFd(os.Stdout.Fd())
	println("term.GetTtySizeByFd(stdout):   ", rows, cols, err)
	rows, cols, err = term.GetTtySizeByFile(os.Stdout)
	println("term.GetTtySizeByFile(stdout): ", rows, cols, err)
	println("term.IsStartupByDoubleClick:   ", term.IsStartupByDoubleClick())

	ctx, cancel := context.WithCancel(context.Background())
	var cancelled int32

	catcher := is.Signals().Catch()
	catcher.
		WithPrompt("Press CTRL-C to quit...").
		// WithVerboseFn(func(msg string, args ...any) {
		// 	fmt.Printf(msg, args...)
		// }).
		// deprecated: WithOnLoop(dbStarter, cacheStarter, mqStarter).
		// not better choice: WithOnLoopFunc((&dbMgr{}).CloseByCatcher).
		WithOnSignalCaught(func(ctx context.Context, sig os.Signal, wg *sync.WaitGroup) {
			println()
			slog.Info("signal caught", "sig", sig)
			atomic.CompareAndSwapInt32(&cancelled, 0, 1)
			cancel() // cancel user's loop, see Wait(...)
		}).
		WaitFor(ctx, func(ctx context.Context, closer func()) {
			slog.Debug("entering looper's loop...")
			go func() {
				// to terminate this app after a while automatically:
				time.Sleep(10 * time.Second)
				// stopChan <- os.Interrupt
				closer()
			}()
			<-ctx.Done() // waiting until any os signal caught
			// wgDone.Done() // and complete myself

			if atomic.CompareAndSwapInt32(&cancelled, 0, 1) {
				is.PressAnyKeyToContinue(os.Stdin)
			} else {
				closer()
			}
		})
}
