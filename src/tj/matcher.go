package main

import (
	"fmt"
	"bytes"
	"regexp"
	"strings"
)

type match struct {
	lineBegin int
	lineEnd	  int
	i, j	  int
}

func (m *match) failed() bool {
	return m.i == 0 && m.j == 0
}

type matchData struct {
	r     rule
	ts    []string
	tsLen int
	t     int
}

func isSpaceOrValue(x, str string) bool {
	switch x {
	case str, " ", "\t", "\v", "\n", "\r":
		return true
	}

	return false
}

func (m *matchData) regexpGetString(nextVal string) string {
	var matchable bytes.Buffer

	// Peek forward in the list of tokens.
	origMt := m.t

	// Regex can be applied until the next constant token or a space.
	for ; m.t < m.tsLen; m.t++ {
		if isSpaceOrValue(m.ts[m.t], nextVal) {
			break
		}

		matchable.WriteString(m.ts[m.t])
	}

	m.t = origMt

	return matchable.String()
}

// Restore the position in the token slice as the last
// unmatched token from the regular expression
func (m *matchData) regexpSkipMatched(matched string) {
	strIndex := 0
	matchedLen := len(matched)

	for ; m.t < m.tsLen; m.t++ {
		tokenLen := len(m.ts[m.t])
		nextIndex := strIndex + tokenLen

		if nextIndex >= matchedLen {
			break
		}

		if m.ts[m.t] != matched[strIndex:nextIndex] {
			break
		}

		strIndex = nextIndex
	}

	m.t++
}

func (m *matchData) matchToken(p, nextP *component) bool {
	switch p.ctype {
	case ctypeConst:
		if m.ts[m.t] == p.value {
			m.t++
			return true
		}
	case ctypeEllipsis:
		if nextP != nil {
			for ; m.t < m.tsLen; m.t++ {
				if m.ts[m.t] == nextP.value {
					return true
				}
			}
		} else {
			m.t = m.tsLen
			return true
		}
	case ctypeRegex:
		var nextVal string

		if nextP != nil {
			nextVal = nextP.value
		}

		matchable := m.regexpGetString(nextVal)

		// TODO: Compiled regex should be cached.
		// TODO: Return error, don't use MustCompile.
		rexp := regexp.MustCompile(p.value)
		matched := rexp.FindString(matchable)

		if matched != "" {
			m.regexpSkipMatched(matched)
			return true
		}
	}

	return false
}

func (m *matchData) getMatchingString(match match) string {
	if match.failed() {
		return ""
	}

	return strings.Join(m.ts[match.i:match.j], "")
}

func matchRule(r rule, ts []string) *match {
	var match match
	var matchStarted bool

	m := &matchData{
		r:     r,
		ts:    ts,
		tsLen: len(ts),
	}

	compLen := len(r.components)
	matchBeginning := m.t

	for i := 0; i < compLen; i++ {
		var matched bool

		if i < compLen-1 {
			matched = m.matchToken(&r.components[i], &r.components[i+1])
		} else {
			matched = m.matchToken(&r.components[i], nil)
		}

		if m.t >= m.tsLen {
			break
		}

		if matched {
			if !matchStarted {
				match.i = matchBeginning
				matchStarted = true
			} else {
				match.j = m.t - 1
			}
		} else {
			if matched {
				match.j = m.t - 1
			}
		}
	}

	if match.j < match.i {
		match.j = m.t
	}

	fmt.Printf("%d %d %s\n", match.i, match.j, m.getMatchingString(match))
	return &match
}

/*
func matchAll(r rule, ts []string) bool {
	for i := 0; i < len(ts); i++ {
		matched := matchRule(r, ts[i:])

		if matched != nil {
			fmt.Printf("Matched \"%s\"\n", strings.Join(ts[i+matched.i:matched.j+i], ""))
			i = matched.j + i
		}
	}

	return false
}
*/
