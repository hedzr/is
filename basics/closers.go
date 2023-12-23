package basics

import (
	"sync/atomic"
)

// RegisterPeripheral adds a peripheral/closable into our global closers set.
// a basics.Peripheral object is a closable instance.
func RegisterPeripheral(servers ...Peripheral) { closers.RegisterPeripheral(servers...) }

// RegisterClosable adds a peripheral/closable into our global closers set.
// a basics.Peripheral object is a closable instance.
func RegisterClosable(servers ...Closable) { closers.RegisterClosable(servers...) }

// RegisterCloseFns adds a simple closure into our global closers set
func RegisterCloseFns(fns ...func()) { closers.RegisterCloseFns(fns...) }

// RegisterCloseFn adds a simple closure into our global closers set
func RegisterCloseFn(fn func()) { closers.RegisterCloseFn(fn) }

// RegisterClosers adds a simple closure into our global closers set
func RegisterClosers(cc ...Closer) { closers.RegisterClosers(cc...) }

// Close will cleanup all registered closers.
// You must make a call to Close before your app shutting down. For example:
//
//	func main() {
//	    defer basics.Close()
//	    // ...
//	}
func Close() {
	closers.Close()
}

// Closers returns the closers set as a basics.Peripheral
func Closers() Peripheral { return closers }

// ClosersClosers returns the closers set as a basics.Peripheral array
func ClosersClosers() []Peripheral { return closers.closers }

var closers = new(c)

type c struct {
	closers []Peripheral
	closed  int32
}

func (s *c) RegisterPeripheral(servers ...Peripheral) {
	s.closers = append(s.closers, servers...)
}

func (s *c) RegisterClosable(servers ...Closable) {
	for _, ci := range servers {
		s.closers = append(s.closers, ci)
	}
}

func (s *c) RegisterCloseFns(fns ...func()) {
	s.closers = append(s.closers, &w{fns})
}

func (s *c) RegisterCloseFn(fn func()) {
	s.closers = append(s.closers, &cf{fn})
}

func (s *c) RegisterClosers(cc ...Closer) {
	for _, ci := range cc {
		if ci != nil {
			s.closers = append(s.closers, &cw{ci})
		}
	}
}

type cw struct {
	cc Closer
}

func (s *cw) Close() {
	if s.cc != nil {
		if err := s.cc.Close(); err != nil {
			println("closing Closer failed.", err)
		}
	}
}

type w struct {
	fns []func()
}

func (s *w) Close() {
	for _, c := range s.fns {
		if c != nil {
			c()
		}
	}
}

type cf struct {
	fn func()
}

func (s *cf) Close() {
	if s.fn != nil {
		s.fn()
	}
}

func (s *c) Close() {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		for _, c := range s.closers {
			// log.Debugf("  c.Close(), %v", c)
			if c != nil {
				c.Close()
			}
		}
	}
}
