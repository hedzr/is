package exec

import (
	"testing"
)

func TestSplitCommandString2(t *testing.T) {
	in := `bash -c 'echo hello world!'`
	out := SplitCommandString(in, '"', '\'')
	t.Log(out)
	if out[0] != "bash" || out[1] != "-c" || out[2] != "echo hello world!" {
		t.Fail()
	}
}
