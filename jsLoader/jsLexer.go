package jsLoader

import (
	"fmt"
)

type token struct {
	tType  tokenType
	lexeme string
	line   int
	column int
}

func (t token) String() string {
	return fmt.Sprintf("{%v, \"%v\", %v:%v}", t.tType, t.lexeme, t.line, t.column)
}

func trimQuotesFromString(s string) string {
	return s[1 : len(s)-1]
}

func isKeyword(t token) bool {
	_, ok := keywords[t.lexeme]
	return ok && t.tType != tNAME
}

func isNumber(c byte) bool {
	return (c >= '0' && c <= '9')
}

func isHexadecimal(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')
}

func isLetter(c byte) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c == '_') || (c == '$')
}

func lex(src []byte) []token {
	tokens := make([]token, 0)
	if len(src) == 0 {
		return tokens
	}

	lastToken := token{}
	lastToken.lexeme = ""
	lastToken.tType = tUNDEFINED

	i := 0
	c := src[i]

	line := 1
	column := 0

	isStringParen := func(c byte) bool {
		return c == '\'' || c == '"' || c == '`'
	}

	eat := func(tType tokenType) {
		lastToken.tType = tType
		lastToken.lexeme = lastToken.lexeme + string(c)
		lastToken.line = line
		lastToken.column = column - len(lastToken.lexeme)

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
		lastToken.tType = tUNDEFINED
	}

	substr := func(start, end int) string {
		if start >= 0 && end <= len(src) {
			return string(src[start:end])
		}
		return ""
	}

	skip := func() {
		i++
		column++
		if i < len(src) {
			c = src[i]
		}
	}

	for i < len(src) {
		switch {

		case substr(i, i+2) == "0x":
			eat(tHEX)
			eat(tHEX)
			for isHexadecimal(c) {
				eat(tHEX)
			}
			end()

		case isNumber(c):
			for isNumber(c) {
				eat(tNUMBER)
			}
			if c == '.' {
				eat(tNUMBER)
				for isNumber(c) {
					eat(tNUMBER)
				}
			}
			end()

		case isLetter(c):
			for isLetter(c) || isNumber(c) {
				eat(tNAME)
			}
			if keyword, ok := keywords[lastToken.lexeme]; ok {
				lastToken.tType = keyword
			}
			end()

		case isStringParen(c):
			startSymbol := c
			eat(tSTRING)
			for !(c == startSymbol && src[i-1] != '\\') {
				eat(tSTRING)
			}
			eat(tSTRING)
			end()

		case c == '\n':
			//skip()
			eat(tNEWLINE)
			end()
			line++
			column = 0

		case c == ' ':
			eat(tSPACE)
			end()

		case substr(i, i+2) == "//":
			for c != '\n' {
				skip()
			}

		case substr(i, i+2) == "/*":
			for substr(i, i+2) != "*/" {
				if c == '\n' {
					line++
				}
				skip()
			}
			skip()
			skip()

		default:
			threeOp := substr(i, i+3)
			twoOp := substr(i, i+2)
			oneOp := substr(i, i+1)

			if op, ok := operators[threeOp]; ok {
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
				skip()
			}
		}
	}

	eat(tEND_OF_INPUT)
	end()

	return tokens
}
