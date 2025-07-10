/*
 * Copyright © 2021 Hedzr Yeh.
 */

package timing

import (
	"io"
	"log"
	"time"
)

// New returns a timing tool for calculating the elapsed time in time.Duration
func New(opts ...Opt) P {
	x := &timeProfiling{tmStart: time.Now(), w: log.Printf}
	for _, opt := range opts {
		opt(x)
	}
	return x
}

// P is a timing tool interface type
type P interface {
	// WithWriter allows putting a printer for dump the timing information
	WithWriter(writer Writer) P
	// WithoutWriter could clear the internal writer
	WithoutWriter() P
	// Duration returns the timing result of an invocation
	Duration() time.Duration

	// CalcNow calculates [Duration] right now to make
	// log.Printf output for the eplaased timing.
	//
	// In general, you need to pass a log.Printf like
	// [Writer] into [P], for example:
	//
	//    p := timing.New(timing.WithWriter(func(msg string, args ...any) {
	//        fmt.Printf(msg, args...)
	//        fmt.Println()
	//    }))
	//    defer p.CalcNow()
	//
	// If nothing applied, `log.Printf` will be used but you
	// can also disable it by `timing.WithoutWriter()`.
	// Then you must have to process the returning of [P.Duration]
	// youself.
	CalcNow()
}

// Opt is a type for implementing functional options pattern
type Opt func(profiling *timeProfiling)

// Writer is a formatter and printer such as log.Printf, t.Logf, ...
type Writer func(msg string, args ...any)

// WithMsgFormat specify a msg-format string such as "xxx takes %v"
func WithMsgFormat(msgTemplate string) Opt {
	return func(p *timeProfiling) {
		p.msg = msgTemplate
	}
}

// WithWriter specify a msg printer/formmater
func WithWriter(w Writer) Opt {
	return func(p *timeProfiling) {
		p.w = w
	}
}

// WithoutWriter specify an empty msg printer/formmater
func WithoutWriter() Opt {
	return func(p *timeProfiling) {
		p.w = nil
	}
}

type timeProfiling struct {
	tmStart time.Time
	w       Writer
	msg     string
}

func (s *timeProfiling) WithWriter(writer Writer) P {
	s.w = writer
	return s
}

func (s *timeProfiling) WithoutWriter() P {
	s.w = nil
	return s
}

func (s *timeProfiling) Duration() time.Duration {
	elapsed := time.Since(s.tmStart)
	if s.w != nil {
		if len(s.msg) == 0 {
			s.msg = "operation lasted %v"
		}
		s.w(s.msg, elapsed)
	}
	return elapsed
}

func (s *timeProfiling) CalcNow() {
	io.Discard.Write([]byte(s.Duration().String()))
}
