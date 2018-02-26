package jsLoader

import (
	"fmt"
)

type token struct {
	tType  tokenType
	lexeme string
	line   int32
	column int32
}

func (t token) String() string {
	return fmt.Sprintf("{%v, \"%v\", %v:%v}", t.tType, t.lexeme, t.line, t.column)
}

func isNumber(c byte) bool {
	return (c >= '0' && c <= '9')
}

func isHexadecimal(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isLetter(c byte) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c == '_') || (c == '$')
}

func lex(src []byte) []token {
	tokens := make([]token, 0)
	if len(src) == 0 {
		return append(tokens, token{t_END_OF_INPUT, "", 0, 0})
	}

	lastToken := token{}
	lastToken.lexeme = ""
	lastToken.tType = t_UNDEFINED

	i := 0
	c := src[i]

	var line int32 = 1
	var column int32 = 0

	eat := func(tType tokenType) {
		lastToken.tType = tType
		lastToken.lexeme = lastToken.lexeme + string(c)
		lastToken.line = line
		lastToken.column = column - int32(len(lastToken.lexeme))

		i++
		column++

		if i < len(src) {
			c = src[i]
		} else {
			c = ' '
		}
	}

	end := func() {
		if len(lastToken.lexeme) > 0 {
			tokens = append(tokens, lastToken)
			lastToken = token{}
		}
		lastToken.lexeme = ""
		lastToken.tType = t_UNDEFINED
	}

	substr := func(start, end int) string {
		if start >= 0 && end <= len(src) {
			return string(src[start:end])
		}
		return ""
	}

	// skip := func() {
	// 	i++
	// 	column++
	// 	if i < len(src) {
	// 		c = src[i]
	// 	}
	// }

	for i < len(src) {
		switch {

		case c == '0' && (src[i+1] == 'o' || src[i+1] == 'O'):
			eat(t_NUMBER)
			eat(t_NUMBER)
			for isNumber(c) {
				eat(t_NUMBER)
			}
			end()

		case c == '0' && (src[i+1] == 'b' || src[i+1] == 'B'):
			eat(t_NUMBER)
			eat(t_NUMBER)
			for isNumber(c) {
				eat(t_NUMBER)
			}
			end()

		case c == '0' && (src[i+1] == 'x' || src[i+1] == 'X'):
			eat(t_NUMBER)
			eat(t_NUMBER)
			for isHexadecimal(c) {
				eat(t_NUMBER)
			}
			end()

		case c == '.' && isNumber(src[i+1]):
			eat(t_NUMBER)
			for isNumber(c) {
				eat(t_NUMBER)
			}
			if c == 'e' || c == 'E' {
				eat(t_NUMBER)
				for isNumber(c) {
					eat(t_NUMBER)
				}
			}
			end()

		case isNumber(c):
			for isNumber(c) {
				eat(t_NUMBER)
			}
			if c == '.' {
				eat(t_NUMBER)
				for isNumber(c) {
					eat(t_NUMBER)
				}
			}
			if c == 'e' || c == 'E' {
				eat(t_NUMBER)
				for isNumber(c) {
					eat(t_NUMBER)
				}
			}
			end()

		case isLetter(c):
			for isLetter(c) || isNumber(c) {
				eat(t_NAME)
			}
			if keyword, ok := keywords[lastToken.lexeme]; ok {
				lastToken.tType = keyword
			}
			end()

		case c == '\'' || c == '"':
			eat(t_STRING_QUOTE)
			end()

		case c == '`':
			eat(t_TEMPLATE_LITERAL_QUOTE)
			end()

		case c == '\n' || c == '\v' || c == '\f':
			eat(t_NEWLINE)
			end()
			line++
			column = 0

		case c == ' ' || c == '\t':
			eat(t_SPACE)
			end()

		default:
			fourOp := substr(i, i+4)
			threeOp := substr(i, i+3)
			twoOp := substr(i, i+2)
			oneOp := substr(i, i+1)

			if op, ok := operators[fourOp]; ok {
				eat(op)
				eat(op)
				eat(op)
				eat(op)
				end()
			} else if op, ok := operators[threeOp]; ok {
				eat(op)
				eat(op)
				eat(op)
				end()
			} else if op, ok := operators[twoOp]; ok {
				eat(op)
				eat(op)
				end()
			} else if op, ok := operators[oneOp]; ok {
				eat(op)
				end()
			} else {
				eat(t_ANY)
				end()
			}
		}
	}

	eat(t_END_OF_INPUT)
	end()

	return tokens
}
