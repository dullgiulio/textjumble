package rules

import (
	"bytes"
	"fmt"
)

type parsePosition int

const (
	posNewline parsePosition = iota
	posNameBegin
	posComponents
	posCompRef
	posCompConst
	posCompRegex
	posCompEll1
	posCompEll2
	posCompEll3
)

type parsable struct {
	text  string
	index int
	line  int
	pos   parsePosition
	file  string
	token bytes.Buffer
	comps []component
}

func newParsable(file, text string) *parsable {
	return &parsable{
		text:  text,
		file:  file,
		comps: make([]component, 0),
	}
}

func (p *parsable) getRules() ([]rule, error) {
	var cs []component
	var rls []rule

	for _, c := range p.comps {
		if len(cs) > 0 {
			if c.ctype == ctypeNil {
				if len(cs) > 1 {
					rls = append(rls, makeRule(cs[0], cs[1:]...))
					cs = nil // Reset the components slice
				} else {
					return rls, c.error("Unexpected rule definition; expected rule component")
				}
			}
		}

		cs = append(cs, c)
	}

	if len(cs) > 0 {
		rls = append(rls, makeRule(cs[0], cs[1:]...))
	}

	return rls, nil
}

func (p *parsable) emitToken(t componentType, i, j int) {
	c := makeComponent(p.token.String(), t)
	c.line = p.line + 1
	c.file = p.file

	p.comps = append(p.comps, c)

	p.token.Reset()
}

func (p *parsable) error(str string, i, j int) error {
	return fmt.Errorf("%s:%d: %d-%d: `%s': %s", p.file, p.line+1, i+1, j+1, p.token.String(), str)
}

func (p *parsable) parse() error {
	var index int
	var c rune
	var quoted bool

	for index, c = range p.text {
		if quoted {
			quoted = false
			p.token.WriteRune(c)
			continue
		}

		switch c {
		case '\n', '\r':
			if p.pos == posComponents {
				p.pos = posNewline
			}
			p.line++

			if p.pos == posCompConst {
				return p.error("String constants can only be on a single line", p.index, index)
			}
		case '\\':
			switch p.pos {
			case posCompConst, posCompRegex:
				quoted = true
			default:
				return p.error("Unexpected slash", p.index, index)
			}
		case '"':
			switch p.pos {
			case posComponents:
				p.index = index
				p.pos = posCompConst
			case posCompConst:
				p.pos = posComponents
				p.emitToken(ctypeConst, p.index, index-1)
				p.index = index
			case posCompRegex:
				p.token.WriteRune(c)
			default:
				return p.error("Unexected quote", p.index, index)
			}
		case ' ', '\t', '\v':
			switch p.pos {
			case posNewline, posComponents:
				continue
			case posNameBegin:
				p.pos = posComponents
				p.emitToken(ctypeNil, p.index, index)
				p.index = index
				p.token.WriteRune(c)
			case posCompRef:
				p.pos = posComponents
				p.emitToken(ctypeReference, p.index, index)
				p.index = index
			case posCompConst, posCompRegex:
				p.token.WriteRune(c)
			default:
				return p.error("Expected space or tabulation", p.index, index)
			}
		case ':':
			switch p.pos {
			case posNameBegin:
				p.pos = posComponents
				p.emitToken(ctypeNil, p.index, index)
				p.index = index
			case posCompConst, posCompRegex:
				p.token.WriteRune(c)
			default:
				return p.error("Unexpected colon", p.index, index)
			}
		case '.':
			switch p.pos {
			case posComponents:
				p.pos = posCompEll1
				p.index = index
				p.token.WriteRune(c)
			case posCompEll1:
				p.pos = posCompEll2
				p.token.WriteRune(c)
			case posCompEll2:
				p.pos = posComponents
				p.token.WriteRune(c)
				p.emitToken(ctypeEllipsis, p.index, index)
			case posCompConst, posCompRegex:
				p.token.WriteRune(c)
			default:
				return p.error("Unexpected dot", p.index, index)
			}
		case '/':
			switch p.pos {
			case posCompConst:
				p.token.WriteRune(c)
			case posComponents:
				p.pos = posCompRegex
				p.index = index
			case posCompRegex:
				p.pos = posComponents
				p.emitToken(ctypeRegex, p.index, index-1)
				p.index = index
			default:
				return p.error("Unexpected slash", p.index, index)
			}
		default:
			switch p.pos {
			case posNewline:
				p.pos = posNameBegin
				p.index = index
				p.token.WriteRune(c)
			case posComponents:
				p.pos = posCompRef
				p.index = index
				p.token.WriteRune(c)
			case posNameBegin, posCompRef, posCompConst, posCompRegex:
				p.token.WriteRune(c)
			default:
				return p.error("Unexpected character after rule name", p.index, index)
			}
		}
	}

	switch p.pos {
	case posCompRef:
		p.emitToken(ctypeReference, p.index, index)
	case posCompRegex:
		return p.error("Unterminated regex", p.index, index)
	case posCompConst:
		return p.error("Unterminated string constant", p.index, index)
	}

	return nil
}
