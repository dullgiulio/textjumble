package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTokenizerSplit(t *testing.T) {
	reader := bytes.NewReader([]byte(`
	public function testMethod($arg1, $arg2) {
}
`))

	tok := newTokenizer(reader)
	done := make(chan struct{})

	go func() {
		expectedTokens := []string{
			"\n", "\t", "public", " ", "function", " ", "testMethod", "(", "$", "arg1", ",", " ",
			"$", "arg2", ")", " ", "{", "\n", "}", "\n",
		}

		n := 0
		nEx := len(expectedTokens)

		for s := range tok.Tokens {
			if n >= nEx {
				t.Error(fmt.Errorf("%d: Generated unexpected token `%s'", n, s))
				continue
			}

			if s != expectedTokens[n] {
				t.Error(fmt.Errorf("%d: Expected token `%s', got `%s'", n, expectedTokens[n], s))
			}

			n++
		}

		done <- struct{}{}
	}()

	if err := tok.split(); err != nil {
		t.Error(err)
	}

	<-done
}
