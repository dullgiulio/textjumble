package main

import (
	"testing"
)

func TestEllipsisBlock(t *testing.T) {
	tokens := []string{
		"{", " ", "$", "variable", " ", "=", " ", "10", "}",
	}

	rule := makeRule(makeRuleName("test"), makeComponent("{", ctypeConst), makeComponent("...", ctypeEllipsis), makeComponent("}", ctypeConst))
	if matchRule(rule, tokens) == nil {
		t.Error("Expected matched rule")
	}

	rule = makeRule(makeRuleName("test"), makeComponent("{", ctypeConst), makeComponent("...", ctypeEllipsis), makeComponent("]", ctypeConst))
	if matchRule(rule, tokens) != nil {
		t.Error("Unexpected matched rule")
	}
}

func TestRegex(t *testing.T) {
	tokens := []string{
		"{", "\n", "\t", "$", "long", "_", "variable", "=", " ", "100", ";", "\n", "}",
	}

	rule := makeRule(makeRuleName("test"), makeComponent("{", ctypeConst), makeComponent("...", ctypeEllipsis),
		makeComponent("$", ctypeConst), makeComponent("[a-zA-Z_][a-zA-Z0-9_]*", ctypeRegex),
		makeComponent("...", ctypeEllipsis), makeComponent("}", ctypeConst))
	if matchRule(rule, tokens) == nil {
		t.Error("Expected matched rule")
	}
}

func TestMultipleMatches(t *testing.T) {
	tokens := []string{
		"{", "\n",
			"\t", "$", "long", "_", "variable", "=", " ", "100", ";", "\n",
			"\t", "$", "another", "_", "variable", " ", "=", " ", "100", ";", "\n",
		"}",
	}

	rule := makeRule(makeRuleName("test"), makeComponent("$", ctypeConst), makeComponent("[a-zA-Z_][a-zA-Z0-9_]*", ctypeRegex), makeComponent("...", ctypeEllipsis), makeComponent(";", ctypeConst))
	//rule := makeRule(makeRuleName("test"), makeComponent("...", ctypeEllipsis), makeComponent(";", ctypeConst))

	if matchRule(rule, tokens) == nil {
		t.Error("Expected matched rule")
	}
}
