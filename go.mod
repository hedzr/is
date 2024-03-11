module github.com/hedzr/is

go 1.21

//replace github.com/hedzr/env => ../libs.env

require (
	golang.org/x/crypto v0.21.0
	golang.org/x/net v0.22.0
	golang.org/x/term v0.18.0
)

require golang.org/x/sys v0.18.0 // indirect
