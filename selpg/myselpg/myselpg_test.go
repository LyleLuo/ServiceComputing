package myselpg

import (
	"bytes"
	"testing"
)

func TestSelpg(t *testing.T) {
	page := new(Page)
	page.StartPage = 4
	page.EndPage = 5
	page.InFileName = "../test/testfile"
	page.PageLen = 2
	buf := new(bytes.Buffer)
	Selpg(page, buf)
	output := buf.String()
	exceptOutput := "\t\"io\"\n\t\"os\"\n\t\"os/exec\"\n\n"
	if output != exceptOutput {
		t.Errorf("Excepted output is:\n%s, but you output is:\n%s", exceptOutput, output)
	}
}
