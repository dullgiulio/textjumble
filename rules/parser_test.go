package rules

import (
	"testing"
)

func TestParseRule(t *testing.T) {
	text := `ruleName: some /re[a-z0-9]+gex/ ... "string. with space: \"and quotes\"" rule`
	p := newParsable("<STDIN>", text)

	if err := p.parse(); err != nil {
		t.Error(err)
	}
}

func TestParseInvalid(t *testing.T) {
	text := `ruleBlah: nonsese: \t some" more`
	p := newParsable("<STDIN>", text)

	if err := p.parse(); err == nil {
		t.Error("Expected some error")
	}
}
