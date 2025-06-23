package basics

type OldInfra interface {
	Open() error
}

type Infrastructure interface {
	Peripheral
	Openable
}
