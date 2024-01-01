# term package

`term.go` and its sub-sequences (term_unix...) are copied from golang.org/x/term, and:

- disable err returned by readPassword when unsupported (and plan9, maybe nacl)
