package main

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

func TestRuleComponents(t *testing.T) {
	text := `class: "class" name block`
	p := newParsable("<STDIN>", text)

	if err := p.parse(); err != nil {
		t.Error(err)
	}

	if rls, err := p.getRules(); err != nil {
		t.Error(err)
	} else {
		if len(rls) != 1 {
			t.Error("Unexpected number or rules parsed")
		}

		if rls[0].String() != text {
			t.Error("Unexpected parsing of rule")
		}
	}
}

func TestMultipleRules(t *testing.T) {
	text := `
class: "class" name block
name:	/[a-zA-Z0-9_][a-zA-Z0-9_]*/
block: "{" ... "}"

`

	p := newParsable("<STDIN>", text)

	if err := p.parse(); err != nil {
		t.Error(err)
	}

	if rls, err := p.getRules(); err != nil {
		t.Error(err)
	} else {
		if len(rls) != 3 {
			t.Error("Unexpected number or rules parsed")
		}
	}
}
