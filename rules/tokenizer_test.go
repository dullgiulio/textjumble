package rules

import (
	"fmt"
	"testing"
	"bytes"
)

func TestTokenizerSplit(t *testing.T) {
	reader := bytes.NewReader([]byte(`
<?php

class TestExampleClass {
        public function testMethod($arg1, $arg2) {
                echo "$arg1, $arg2"; // Just print the arguments
        }

        public function anotherTestMethod() {
                /*
                 * Define some array
                 */
                $arr = array(1, 2, 3, "string with many words");

                foreach ($arr as $a) {
                        echo $a;

                        return $a + $a;
                }
        }
}

?>
`))

	tok := newTokenizer(reader)
	done := make(chan struct{})

	go func() {
		for s := range tok.Tokens {
			fmt.Printf("%s\n", s)
		}

		done <- struct{}{}
	}()

	if err := tok.split(); err != nil {
		t.Error(err)
	}

	<-done
}
