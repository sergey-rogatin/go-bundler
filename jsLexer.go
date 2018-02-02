package main

import (
	"fmt"
	"strings"
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

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
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

	line := 0
	column := 0

	isWhitespace := func(c byte) bool {
		return c == ' '
	}

	isStringParen := func(c byte) bool {
		return c == '\'' || c == '"' || c == '`'
	}

	eat := func(tType tokenType) {
		lastToken.tType = tType
		lastToken.lexeme = lastToken.lexeme + string(c)
		lastToken.line = line
		lastToken.column = column
		// fmt.Println(lastToken)
		i++

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

		case c == '\n':
			skip()
			//eat(tNEWLINE)
			end()
			line++
			column = 0

		case isWhitespace(c):
			skip()

		case substr(i, i+2) == "//":
			for c != '\n' {
				skip()
			}

		case substr(i, i+2) == "/*":
			for substr(i, i+2) != "*/" {
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
		column++
	}

	eat(tEND_OF_INPUT)
	end()

	return tokens
}

type importFileInfo struct {
	exportObjName string
	resolvedPath  string
	def           string
	vars          []string
}

type exportInfo struct {
	def  []token
	vars []token
}

func resolveES6ImportPath(importPath, currentFileName string) string {
	importPath = trimQuotesFromString(importPath)
	pathParts := strings.Split(importPath, "/")

	locationParts := strings.Split(currentFileName, "/")
	locationParts = locationParts[:len(locationParts)-1]

	for _, part := range pathParts {
		if part == ".." {
			locationParts = locationParts[:len(locationParts)-1]
			pathParts = pathParts[1:]
		}
		if part == "." {
			pathParts = pathParts[1:]
		}
	}

	fullFileName := strings.Join(append(locationParts, pathParts...), "/")

	ext := ""
	if strings.Index(pathParts[len(pathParts)-1], ".") < 0 {
		ext = ".js"
	}

	result := fullFileName + ext
	return result
}

func transformIntoModule(programAst program, fileName string) (program, []string) {
	result := program{}
	result.statements = []ast{}
	fileImports := []string{}

	for _, item := range programAst.statements {
		switch st := item.(type) {
		case importStatement:
			resolvedPath := resolveES6ImportPath(st.path, fileName)
			exportObjName := createVarNameFromPath(resolvedPath)

			fileImports = append(fileImports, resolvedPath)
			if len(st.vars) > 0 {
				vs := varStatement{}
				vs.keyword = "var"
				vs.decls = []declaration{}
				for _, impVar := range st.vars {
					decl := declaration{}
					decl.name = impVar.pseudonym
					decl.value = memberExpression{
						object:   atom{token{tType: tNAME, lexeme: exportObjName}},
						property: impVar.name,
					}
					vs.decls = append(vs.decls, decl)
				}
				result.statements = append(result.statements, vs)
			}

		default:
			result.statements = append(result.statements, st)
		}
	}

	fmt.Println(result)
	return result, nil
}

func writeTokens(tokens []token) []byte {
	result := make([]byte, 0)

	for i, t := range tokens {
		tIsKeyword := isKeyword(t)

		if tIsKeyword && i > 0 && (tokens[i-1].tType == tNAME || tokens[i-1].tType == tNUMBER || isKeyword(tokens[i-1])) {
			result = append(result, ' ')
		}
		result = append(result, []byte(t.lexeme)...)
		if tIsKeyword && i < len(tokens) && (tokens[i+1].tType == tNAME || tokens[i+1].tType == tNUMBER) {
			result = append(result, ' ')
		}
	}

	return result
}

func loadJsFile(src []byte, filePath string) ([]byte, []string) {
	tokens := lex(src)
	srcAst := parse(tokens)
	programAst, fileImports := transformIntoModule(srcAst, filePath)
	resultBytes := []byte(programAst.String())

	return resultBytes, fileImports
}
