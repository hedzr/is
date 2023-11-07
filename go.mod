module github.com/hedzr/is

go 1.21

//replace github.com/hedzr/go-errors/v2 => ../libs.errors

//replace github.com/hedzr/env => ../libs.env

// replace github.com/hedzr/go-utils/v2 => ./

require (
	github.com/konsorten/go-windows-terminal-sequences v1.0.3
	golang.org/x/crypto v0.14.0
	golang.org/x/net v0.17.0
	golang.org/x/sys v0.14.0
)

require golang.org/x/term v0.13.0 // indirect
