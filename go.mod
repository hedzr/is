module github.com/hedzr/is

go 1.22.7

//replace github.com/hedzr/env => ../libs.env

require (
	golang.org/x/crypto v0.27.0
	golang.org/x/net v0.29.0
	golang.org/x/term v0.24.0
)

require golang.org/x/sys v0.25.0 // indirect
