module github.com/hedzr/is

go 1.23.0

toolchain go1.23.3

//replace github.com/hedzr/env => ../libs.env

require (
	golang.org/x/net v0.39.0
	golang.org/x/term v0.31.0
)

require golang.org/x/sys v0.32.0 // indirect
