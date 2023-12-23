package basics

import (
	"context"
)

type OldInfra interface {
	Open() error
}

type Infrastructure interface {
	Peripheral

	// Open does initializing stuffs
	Open(ctx context.Context) (err error)
}
