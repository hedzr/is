package basics

import "github.com/hedzr/env/ctx"

type OldInfra interface {
	Open() error
}

type Infrastructure interface {
	Peripheral

	// Open does initializing stuffs
	Open(ctx ctx.Context) (err error)
}
