package rules

type matchData struct {
	r rule
	ts []string
	tsLen int
	t int
}

func (m *matchData) token(p, nextP *component) bool {
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
	}

	return false
}

func matchRule(r rule, ts []string) (matched bool) {
	m := &matchData{
		r: r,
		ts: ts,
		tsLen: len(ts),
	}

	compLen := len(r.components)

	for i := 0; i < compLen; i++ {
		if i < compLen - 1 {
			matched = m.token(&r.components[i], &r.components[i+1])
		} else {
			matched = m.token(&r.components[i], nil)
		}

		if m.t >= m.tsLen {
			break
		}
	}

	return
}
