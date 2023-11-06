package color

import (
	"testing"
)

func TestGetCPT(t *testing.T) {
	t.Logf("%v", GetCPT().Translate(`<code>code</code> | <kbd>CTRL</kbd>
	<b>bold / strong / em</b>
	<i>italic / cite</i>
	<u>underline</u>
	<mark>inverse mark</mark>
	<del>strike / del </del>
	<font color="green">green text</font>
	`, FgDefault))
}

func TestGetCPTC(t *testing.T) {
	t.Logf("%v", GetCPTC().Translate(`<code>code</code>`, FgDefault))
}

func TestGetCPTNC(t *testing.T) {
	t.Logf("%v", GetCPTNC().Translate(`<code>code</code>`, FgDefault))
}

func TestStripLeftTabsC(t *testing.T) {
	t.Logf("%v", StripLeftTabsC(`
	
		<code>code</code>
	NC Cool
	 But it's tight.
	  Hold On!
	Hurry Up.
	`))
}

func TestStripLeftTabs(t *testing.T) {
	t.Logf("%v", StripLeftTabs(`
	
		<code>code</code>
	NC Cool
	 But it's tight.
	  Hold On!
	Hurry Up.
	`))
}

func TestStripLeftTabsOnly(t *testing.T) {
	t.Logf("%v", StripLeftTabsOnly(`
	
		<code>code</code>
	NC Cool
	 But it's tight.
	  Hold On!
	Hurry Up.
	`))
}

func TestStripHTMLTags(t *testing.T) {
	t.Logf("%v", StripHTMLTags(`
	
		<code>code</code>
	NC Cool
	 But it's tight.
	  Hold On!
	Hurry Up.
	`))
}
