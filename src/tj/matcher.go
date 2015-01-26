package main

import (
	"bytes"
	"regexp"
)

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

func matchRule(r rule, ts []string) (matched bool) {
	m := &matchData{
		r:     r,
		ts:    ts,
		tsLen: len(ts),
	}

	compLen := len(r.components)

	for i := 0; i < compLen; i++ {
		if i < compLen-1 {
			matched = m.matchToken(&r.components[i], &r.components[i+1])
		} else {
			matched = m.matchToken(&r.components[i], nil)
		}

		if m.t >= m.tsLen {
			break
		}
	}

	return
}
