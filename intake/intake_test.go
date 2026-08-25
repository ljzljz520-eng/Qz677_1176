package intake

import (
	"strings"
	"testing"
)

func TestParseAndValidate(t *testing.T) {
	rows, parseIssues := NewImporter().ParseCSV(strings.NewReader("id,citizen,name,amount,region\nr1,c1,Alice,400,bj\nr2,c2,Bob,nope,sh"))
	if len(rows) != 1 || len(parseIssues) != 1 {
		t.Fatalf("unexpected parse: %#v %#v", rows, parseIssues)
	}
	if len(ValidateBatch(rows)) != 0 {
		t.Fatal("valid row rejected")
	}
}
