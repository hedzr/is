module github.com/hedzr/is

go 1.22.7

//replace github.com/hedzr/env => ../libs.env

require (
	golang.org/x/crypto v0.33.0
	golang.org/x/net v0.35.0
	golang.org/x/term v0.29.0
)

require golang.org/x/sys v0.30.0 // indirect
