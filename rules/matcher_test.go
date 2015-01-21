package rules

import (
	"testing"
)

func TestEllipsisBlock(t *testing.T) {
	tokens := []string{
		"{", " ", "$", "variable", " ", "=", " ", "10", "}",
	}

	rule := makeRule(makeRuleName("test"), makeComponent("{", ctypeConst), makeComponent("...", ctypeEllipsis), makeComponent("}", ctypeConst))
	if !matchRule(rule, tokens) {
		t.Error("Expected matched rule")
	}

	rule = makeRule(makeRuleName("test"), makeComponent("{", ctypeConst), makeComponent("...", ctypeEllipsis), makeComponent("]", ctypeConst))
	if matchRule(rule, tokens) {
		t.Error("Unexpected matched rule")
	}
}
