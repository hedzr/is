module github.com/hedzr/is

go 1.23.7

//replace github.com/hedzr/env => ../libs.env

require (
	golang.org/x/crypto v0.36.0
	golang.org/x/net v0.37.0
	golang.org/x/term v0.30.0
)

require golang.org/x/sys v0.31.0 // indirect
