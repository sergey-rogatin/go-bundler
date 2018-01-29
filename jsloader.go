package main

import (
	"fmt"
	"strings"
)

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
		return c == ' '
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
			eat(tNEWLINE)
			end()

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

func transformIntoModule(tokens []token, resolvedPath string) ([]token, []importFileInfo) {
	result := make([]token, 0, len(tokens)+128)
	imports := make([]importFileInfo, 0)

	if len(tokens) == 0 {
		return result, imports
	}

	fileExports := exportInfo{}

	i := 0
	t := tokens[i]

	write := func(tType tokenType, lexeme string) {
		if tType == tNAME {
			// check if imported variable
			for _, imp := range imports {
				ext := getExtension(imp.resolvedPath)
				if ext != "js" {
					continue
				}
				if imp.def == lexeme {
					result = append(result, token{tNAME, imp.exportObjName})
					result = append(result, token{tDOT, "."})
					result = append(result, token{tNAME, "default"})
					return
				}
				for _, importVar := range imp.vars {
					if importVar == lexeme {
						result = append(result, token{tNAME, imp.exportObjName})
						result = append(result, token{tDOT, "."})
						result = append(result, token{tNAME, lexeme})
						return
					}
				}
			}
		}
		result = append(result, token{tType, lexeme})
	}

	eat := func() {
		write(t.tType, t.lexeme)
		i++
		if i < len(tokens) {
			t = tokens[i]
		}
	}

	skip := func() {
		i++
		if i < len(tokens) {
			t = tokens[i]
		}
	}

	back := func() {
		i--
		if i > 0 {
			t = tokens[i]
		}
	}

	exportObjName := createVarNameFromPath(resolvedPath)
	// add module pattern
	write(tVAR, "var")
	write(tNAME, exportObjName)
	write(tEQUALS, "=")
	write(tPAREN_LEFT, "(")
	write(tFUNCTION, "function")
	write(tPAREN_LEFT, "(")
	write(tPAREN_RIGHT, ")")
	write(tCURLY_LEFT, "{")

	for i < len(tokens) {
		switch t.tType {
		case tNEWLINE:
			skip()

		case tIMPORT:
			jsImport := importFileInfo{}
			jsImport.vars = make([]string, 0)

			for t.tType != tSEMI {
				// no curly brace encountered
				if t.tType == tNAME {
					jsImport.def = t.lexeme
				}
				// destructuring import
				if t.tType == tCURLY_LEFT {
					for t.tType != tCURLY_RIGHT {
						if t.tType == tDEFAULT {
							skip()                  // default
							skip()                  // as
							jsImport.def = t.lexeme // foo
						} else if t.tType == tNAME {
							jsImport.vars = append(jsImport.vars, t.lexeme)
						}
						skip()
					}
				}
				skip()
			}
			// end of import statement found
			back()
			importPath := resolveES6ImportPath(t.lexeme, resolvedPath)
			skip() // "./foo"
			skip() // ;

			ext := getExtension(importPath)
			formattedName := createVarNameFromPath(importPath)

			jsImport.resolvedPath = importPath
			jsImport.exportObjName = formattedName
			imports = append(imports, jsImport)

			if ext != "js" {
				fullFileName := formattedName + "." + ext
				write(tVAR, "var")
				write(tNAME, jsImport.def)
				write(tASSIGN, "=")
				write(tSTRING, fmt.Sprintf("\"%s\"", fullFileName))
				write(tSEMI, ";")
			}
		case tEXPORT:
			skip() // export
			if t.tType == tDEFAULT {
				skip() // default
				for t.tType != tSEMI {
					fileExports.def = append(fileExports.def, t)
					skip()
				}
			} else {
				for t.tType != tSEMI {
					if t.tType == tFUNCTION {
						eat()
						fileExports.vars = append(fileExports.vars, t)
						break
					} else if t.tType == tNAME {
						fileExports.vars = append(fileExports.vars, t)
					}
					eat()
				}
				eat()
			}

		default:
			eat()
		}
	}

	// append exports object return
	write(tRETURN, "return")
	write(tCURLY_LEFT, "{")

	if len(fileExports.def) > 0 {
		write(tNAME, "default")
		write(tCOLON, ":")

		for _, defToken := range fileExports.def {
			write(defToken.tType, defToken.lexeme)
		}

		write(tCOMMA, ",")
	}

	for _, varToken := range fileExports.vars {
		write(varToken.tType, varToken.lexeme)
		write(tCOMMA, ",")
	}

	write(tCURLY_RIGHT, "}")
	write(tSEMI, ";")

	// finish module pattern
	write(tCURLY_RIGHT, "}")
	write(tPAREN_RIGHT, ")")
	write(tPAREN_LEFT, "(")
	write(tPAREN_RIGHT, ")")
	write(tSEMI, ";")

	return result, imports
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
	resultTokens, fileImports := transformIntoModule(tokens, filePath)
	resultBytes := writeTokens(resultTokens)

	fileImportPaths := make([]string, 0)
	for _, imp := range fileImports {
		fileImportPaths = append(fileImportPaths, imp.resolvedPath)
	}
	return resultBytes, fileImportPaths
}
