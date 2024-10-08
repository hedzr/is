module github.com/hedzr/is

go 1.22.7

//replace github.com/hedzr/env => ../libs.env

require (
	golang.org/x/crypto v0.28.0
	golang.org/x/net v0.29.0
	golang.org/x/term v0.25.0
)

require golang.org/x/sys v0.26.0 // indirect
