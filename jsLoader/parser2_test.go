package jsLoader

import "testing"

var ps parser

func setParser(text string) {
	toks := lex([]byte(text))
	ps = parser{
		tokens: toks,
		i:      0,
		t:      toks[0],
	}
}

func TestFunctionDeclaration(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"foo();",
			"",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		t.Errorf("%v", le)

		res := "" // generateJsCode(le)
		if res != c.exp {
			// t.Errorf("Expected %s, got %s", c.exp, generateJsCode(le))
		}
	}
}
