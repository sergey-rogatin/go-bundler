package htmlLoader

import "fmt"

type tokenType int32

const (
	t_UNDEFINED tokenType = iota
	t_ANGLE_BRACKET_LEFT
	t_ANGLE_BRACKET_RIGHT
	t_SLASH
	t_WHITESPACE
	t_ESCAPE_SEQUENCE
	t_STRING
	t_ASSIGN
	t_WORD
	t_END_OF_INPUT
)

var tokenToString = []string{
	"t_UNDEFINED",
	"t_ANGLE_BRACKET_LEFT",
	"t_ANGLE_BRACKET_RIGHT",
	"t_SLASH",
	"t_WHITESPACE",
	"t_ESCAPE_SEQUENCE",
	"t_STRING",
	"t_ASSIGN",
	"t_WORD",
	"t_END_OF_INPUT",
}

func (t tokenType) String() string {
	return tokenToString[int32(t)]
}

type token struct {
	tType  tokenType
	lexeme string
	line   int32
	column int32
}

func (t token) String() string {
	return fmt.Sprintf("{%v, \"%v\", %v:%v}", t.tType, t.lexeme, t.line, t.column)
}

func lex(src []byte) []token {
	if len(src) == 0 {
		return nil
	}

	tokens := []token{}

	i := 0
	c := src[i]
	line := int32(1)
	column := int32(1)

	lastToken := token{}

	skip := func() {
		i++
		column++

		if i < len(src) {
			c = src[i]
		}
	}

	eat := func() {
		lastToken.lexeme += string(c)
		skip()
	}

	end := func(tType tokenType) {
		lastToken.tType = tType
		lastToken.line = line
		lastToken.column = column - int32(len(lastToken.lexeme))
		tokens = append(tokens, lastToken)

		lastToken = token{}
	}

	for i < len(src) {
		switch {
		case c == '<':
			eat()
			end(t_ANGLE_BRACKET_LEFT)

		case c == '>':
			eat()
			end(t_ANGLE_BRACKET_RIGHT)

		case c == '/':
			eat()
			end(t_SLASH)

		case c == '=':
			eat()
			end(t_ASSIGN)

		case c == ' ' || c == '\n':
			eat()
			end(t_WHITESPACE)
			for c == ' ' || c == '\n' {
				skip()
			}

		case c == '"' || c == '\'':
			firstQuote := c
			eat()
			for c != firstQuote {
				eat()
			}
			eat()
			end(t_STRING)

		default:
			for c != '<' &&
				c != '>' &&
				c != '/' &&
				c != '=' &&
				c != ' ' &&
				c != '\n' &&
				c != '"' &&
				c != '\'' {
				eat()
			}
			end(t_WORD)

		}
	}
	end(t_END_OF_INPUT)

	return tokens
}
