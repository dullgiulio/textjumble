package rules

import (
	"io"
	"bytes"
	"bufio"
)

type tokenizer struct {
	reader	*bufio.Reader
	Tokens	chan string
}

const utf8BOM rune = '\uFEFF'
const stoppers string = " \t\v\n-_=+~!@#$%^&*()[]{}:;\"'\\|?/>.<,"

func newTokenizer(reader io.Reader) *tokenizer {
	return &tokenizer{
		reader: bufio.NewReader(reader),
		Tokens: make(chan string),
	}
}

func (t *tokenizer) next() string {
	return ""
}

func isInString(a rune, b string) bool {
	for _, r := range b {
		if a == r {
			return true
		}
	}

	return false
}

func (t *tokenizer) split() error {
	var buf bytes.Buffer

	defer close(t.Tokens)

	for {
		r, _, err := t.reader.ReadRune()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		if isInString(r, stoppers) {
			if buf.Len() > 0 {
				t.Tokens <- buf.String()
				buf.Reset()
			}

			t.Tokens <- string([]rune{r})
		} else {
			buf.WriteRune(r)
		}
	}

	if buf.Len() > 0 {
		t.Tokens <- buf.String()
	}

	return nil
}
