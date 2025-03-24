/*
 * Copyright © 2021 Hedzr Yeh.
 */

package timing

import (
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
}

// Opt is a type for implementing functional options pattern
type Opt func(profiling *timeProfiling)

// Writer is a formatter and printer such as log.Printf, t.Logf, ...
type Writer func(msg string, args ...interface{})

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
