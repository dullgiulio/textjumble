package lexer

import (
	"fmt"
	"testing"
)

func verifyFlat(t *testing.T, rls rules) {
	if rs, err := rls.flatten(); err != nil {
		t.Error(err)
	} else {
		rulesGen := rs.all()
		for r := range rulesGen {
			if !r.isPure() {
				t.Error(fmt.Errorf("Definition not pure: %s", &r))
			}
		}
	}
}

func TestToString(t *testing.T) {
	var rules = map[string]rule{
		"regex: /[a-zA-Z0-9]+/": makeRule("regex", makeComponent("[a-zA-Z0-9]+", ctypeRegex)),
		"ellipsis: ...":         makeRule("ellipsis", makeComponent("", ctypeEllipsis)),
		"reference: referenced": makeRule("reference", makeComponent("referenced", ctypeReference)),
		"const: \"const\"":      makeRule("const", makeComponent("const", ctypeConst)),
	}

	for s, r := range rules {
		if r.String() != s {
			t.Error(fmt.Errorf("Expected '%s', got '%s'", s, r.String()))
		}
	}
}

func TestSimpleRuleResolution(t *testing.T) {
	rules := makeRules()

	/*

		class: "class"
		function: class "function"

	*/

	rules.add("class", makeRule("class", makeComponent("class", ctypeConst)))
	rules.add("function", makeRule("function", makeComponent("class", ctypeReference), makeComponent("function", ctypeConst)))

	verifyFlat(t, rules)
}

func TestMultipleRulesResolution(t *testing.T) {
	rules := makeRules()

	/*

		module: "module"
		class:  module
		class:  "class"
		function: class "function"

	*/

	rules.add("module", makeRule("module", makeComponent("module", ctypeConst)))
	rules.add("class", makeRule("class", makeComponent("class", ctypeConst)))
	rules.add("class", makeRule("class", makeComponent("module", ctypeReference)))
	rules.add("function", makeRule("function", makeComponent("class", ctypeReference), makeComponent("function", ctypeConst)))

	verifyFlat(t, rules)
}

func TestCyclicalReference(t *testing.T) {
	rules := makeRules()

	/*

		module: function
		function: module
		doable: "doable"
		mix: function doable module

	*/

	rules.add("module", makeRule("module", makeComponent("function", ctypeReference)))
	rules.add("function", makeRule("function", makeComponent("module", ctypeReference)))
	rules.add("doable", makeRule("doable", makeComponent("doable", ctypeConst)))
	rules.add("mix", makeRule("mix", makeComponent("function", ctypeReference), makeComponent("doable", ctypeReference), makeComponent("module", ctypeReference)))

	if _, err := rules.flatten(); err == nil {
		t.Error("Expected an error on resolving cyclical references")
	} else {
		if _, ok := err.(*referenceError); !ok {
			t.Error("Expected a referenceError when encountering cyclical references")
		}
	}
}
