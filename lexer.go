package main

import "fmt"

type token struct {
	tType  tokenType
	lexeme string
}

func (t token) String() string {
	return fmt.Sprintf("{%v, \"%v\"}", t.tType, t.lexeme)
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

	isNumber := func(c byte) bool {
		return c >= '0' && c <= '9'
	}

	isLetter := func(c byte) bool {
		return (c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			(c == '_') || (c == '$') || (c == '@')
	}

	isWhitespace := func(c byte) bool {
		return c == ' ' || c == '\n'
	}

	isStringParen := func(c byte) bool {
		return c == '\'' || c == '"' || c == '`'
	}

	eat := func(tType tokenType) {
		lastToken.tType = tType
		lastToken.lexeme = lastToken.lexeme + string(c)
		i++

		if i < len(src) {
			c = src[i]
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
		if start >= 0 && end < len(src) {
			return string(src[start:end])
		}
		return ""
	}

	skip := func() {
		i++
		if i < len(src) {
			c = src[i]
		}
	}

	for i < len(src) {
		switch {
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
			eat(tSTRING)
			for !isStringParen(c) {
				eat(tSTRING)
			}
			eat(tSTRING)
			end()

		case isWhitespace(c):
			skip()

		case substr(i, i+2) == "//":
			for c != '\n' {
				skip()
			}

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

	eat(tUNDEFINED)

	return tokens
}
