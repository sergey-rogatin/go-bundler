package jsLoader

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

func trimQuotesFromString(s string) string {
	return s[1 : len(s)-1]
}

func isKeyword(t token) bool {
	_, ok := keywords[t.lexeme]
	return ok && t.tType != tNAME
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

func makeToken(text string) token {
	res := lex([]byte(text))
	return res[0]
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

func GetFileExtension(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "js"
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

	expObj := CreateVarNameFromPath(fileName)

	bs := blockStatement{}
	bs.statements = []ast{}

	returnExport := objectLiteral{}
	returnExport.properties = []objectProperty{}

	for _, item := range programAst.statements {
		switch st := item.(type) {

		case exportStatement:
			for _, expVar := range st.exportedVars {
				prop := objectProperty{}
				prop.key = expVar.name
				if st.isDeclaration {
					bs.statements = append(bs.statements, expVar.value)
				} else {
					prop.value = expVar.value
				}
				returnExport.properties = append(returnExport.properties, prop)
			}

		case importStatement:
			resolvedPath := resolveES6ImportPath(st.path, fileName)
			exportObjName := CreateVarNameFromPath(resolvedPath)

			fileImports = append(fileImports, resolvedPath)
			ext := GetFileExtension(resolvedPath)

			if len(st.vars) > 0 {
				vs := varStatement{}
				vs.keyword = "var"
				vs.decls = []declaration{}
				for _, impVar := range st.vars {
					decl := declaration{}
					decl.name = impVar.pseudonym
					if ext == "js" {
						decl.value = memberExpression{
							object:   atom{makeToken(exportObjName)},
							property: impVar.name,
						}
					} else {
						decl.value = atom{makeToken("\"" + exportObjName + "." + ext + "\"")}
					}
					vs.decls = append(vs.decls, decl)
				}
				bs.statements = append(bs.statements, vs)
			}

		default:
			bs.statements = append(bs.statements, st)
		}
	}

	rs := returnStatement{}
	rs.value = returnExport

	bs.statements = append(bs.statements, rs)

	fe := functionExpression{}
	fe.body = bs

	fc := functionCall{}
	fc.name = fe

	decl := declaration{}
	decl.name = atom{makeToken(expObj)}
	decl.value = fc

	vs := varStatement{}
	vs.decls = []declaration{decl}
	vs.keyword = "var"

	result.statements = append(result.statements, vs)

	return result, fileImports
}

func LoadFile(src []byte, filePath string) ([]byte, []string) {
	tokens := lex(src)
	srcAst := parse(tokens)
	programAst, fileImports := transformIntoModule(srcAst, filePath)
	resultBytes := []byte(programAst.String())

	return resultBytes, fileImports
}

func CreateVarNameFromPath(path string) string {
	newName := strings.Replace(path, "/", "_", -1)
	newName = strings.Replace(newName, ".", "_", -1)
	return newName
}
