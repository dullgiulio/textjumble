package main

import (
	"bytes"
	"errors"
	"fmt"
)

var loopErr = errors.New("Loop or cyclical reference detected")

type referenceError struct {
	err  string
	name string
	r    rule
}

func newReferenceError(err, name string, r rule) *referenceError {
	return &referenceError{
		err:  err,
		name: name,
		r:    r,
	}
}

func (rf *referenceError) Error() string {
	return fmt.Sprintf("%s: `%s' in `%s'", rf.err, rf.name, &rf.r)
}

type componentType int

const (
	ctypeNil componentType = iota
	ctypeRegex
	ctypeConst
	ctypeEllipsis
	ctypeReference
)

type component struct {
	ctype componentType
	value string
	line  int
	file  string
}

func (c *component) String() string {
	switch c.ctype {
	case ctypeRegex:
		return fmt.Sprintf("/%s/", c.value) // TODO: Escape stuff
	case ctypeConst:
		return fmt.Sprintf("\"%s\"", c.value) // TODO: Escape "
	case ctypeEllipsis:
		return "..."
	case ctypeReference:
		return c.value
	}

	return "<nil>"
}

func (c *component) error(str string) error {
	return fmt.Errorf("%s:%d: `%s': %s", c.file, c.line, c.value, str)
}

func makeComponent(value string, ctype componentType) component {
	return component{
		ctype: ctype,
		value: value,
	}
}

func makeRuleName(value string) component {
	return component{
		value: value,
	}
}

type rule struct {
	name       component
	components []component
}

func (r *rule) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "%s:", r.name.value)

	for _, c := range r.components {
		fmt.Fprintf(&b, " %s", &c)
	}

	return b.String()
}

func (r rule) isPure() bool {
	for _, c := range r.components {
		if c.ctype == ctypeReference {
			return false
		}
	}

	return true
}

func (r rule) replace(sr rule) rule {
	nr := makeRuleEmpty(r.name)

	for _, c := range r.components {
		if c.ctype == ctypeReference && c.value == sr.name.value {
			nr.components = append(nr.components, sr.components...)
		} else {
			nr.components = append(nr.components, c)
		}
	}

	return nr
}

func makeRule(name component, cs ...component) rule {
	return rule{
		name:       name,
		components: cs,
	}
}

func makeRuleEmpty(name component) rule {
	return rule{
		name:       name,
		components: make([]component, 0),
	}
}

type rules map[string][]rule

func makeRules() rules {
	return make(map[string][]rule)
}

func (rs rules) add(name string, r rule) {
	if _, ok := rs[name]; !ok {
		rs[name] = make([]rule, 0)
	}

	rs[name] = append(rs[name], r)
}

func (rs rules) get(name string) ([]rule, bool) {
	r, ok := rs[name]
	return r, ok
}

func (rs rules) all() <-chan rule {
	ch := make(chan rule)

	go func() {
		for _, rls := range rs {
			for _, r := range rls {
				ch <- r
			}
		}

		close(ch)
	}()

	return ch
}

func (rs rules) merge(nrs rules) {
	ch := nrs.all()

	for r := range ch {
		rs.add(r.name.value, r)
	}
}

func (rs rules) partitionByPurity() (a, b rules) {
	a = makeRules()
	b = makeRules()

	rlsGen := rs.all()

	for r := range rlsGen {
		if r.isPure() {
			a.add(r.name.value, r)
		} else {
			b.add(r.name.value, r)
		}
	}

	return
}

func (rs rules) replaceReferences(reps rules) (rules, *referenceError) {
	res := makeRules()
	rulesGen := rs.all()

	for r := range rulesGen {
		for _, c := range r.components {
			if c.ctype == ctypeReference {
				if replacements, ok := reps.get(c.value); !ok {
					return res, newReferenceError("Invalid reference or cyclical definition", c.value, r)
				} else {
					for _, rep := range replacements {
						res.add(r.name.value, r.replace(rep))
					}
				}

				if replacements, ok := res.get(c.value); ok {
					for _, rep := range replacements {
						res.add(r.name.value, r.replace(rep))
					}
				}
			}
		}
	}

	return res, nil
}

func (rs rules) flatten() (rules, error) {
	var lastNSize int

	pure := makeRules()
	nonPure := rs

	for {
		p, np := nonPure.partitionByPurity()

		nonPure = np
		pure.merge(p)

		nonPureLen := len(nonPure)

		if nonPureLen == 0 {
			break
		}

		// If the last partition didn't find any new pure definition,
		// there is some cycle and we cannot resolve everything.
		if lastNSize != 0 && nonPureLen == lastNSize {
			return nonPure, loopErr
		}

		if np, err := nonPure.replaceReferences(pure); err != nil {
			return nonPure, err
		} else {
			nonPure = np
		}

		// Keep the number of not resolved rules to detect cyclical references.
		lastNSize = nonPureLen
	}

	return pure, nil
}
