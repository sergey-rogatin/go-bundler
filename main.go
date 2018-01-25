package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func matchByte(expr string, c byte) bool {
	hasMatch, _ := regexp.Match(expr, []byte{c})
	return hasMatch
}

func matchBytes(expr string, b []byte) bool {
	hasMatch, _ := regexp.Match(expr, b)
	return hasMatch
}

func matchString(expr, str string) bool {
	hasMatch, _ := regexp.MatchString(expr, str)
	return hasMatch
}

func contains(arr []byte, c byte) bool {
	for _, b := range arr {
		if b == c {
			return true
		}
	}
	return false
}

func containsStr(arr []string, c string) bool {
	for _, b := range arr {
		if b == c {
			return true
		}
	}
	return false
}

func containsInt(arr []int, c int) bool {
	for _, b := range arr {
		if b == c {
			return true
		}
	}
	return false
}

func errorInvalidToken(c byte) {
	panic(fmt.Errorf("Error: invalid token %c", c))
}

func trimJsString(jsString string) string {
	return jsString[1 : len(jsString)-1]
}

func resolveImportPath(path, currentFileName string) string {
	path = trimJsString(path)
	pathParts := strings.Split(path, "/")

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

func getExtension(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "js"
}

func formatImportPath(fileName string) string {
	newName := strings.Replace(fileName, "/", "_", -1)
	newName = strings.Replace(newName, ".", "_", -1)
	return newName
}

func getLexemes(src []byte) [][]byte {
	special := "\\{|\\}|\\[|\\]|\\(|\\)|\\,|\\?"
	operators := "\\+|\\-|\\=|\\.|\\>|\\<|\\/|\\*|\\%|\\||\\&|\\:|\\!"
	letters := "[A-Za-z_@$]"
	stringParens := "\\'|\"|\\`"

	state := "sNone"

	lexemes := make([][]byte, 0)
	lastLexeme := make([]byte, 0)

	i := 0
	var c byte

	eat := func() {
		lastLexeme = append(lastLexeme, c)
		i++
	}

	skip := func() {
		i++
	}

	change := func(st string) {
		state = st
		if len(lastLexeme) > 0 {
			lexemes = append(lexemes, lastLexeme)
			lastLexeme = make([]byte, 0)
		}
	}

	reset := func() {
		change("sNone")
	}

	for i < len(src) {
		c = src[i]

		switch state {
		case "sNone":

			switch {
			case c == '/' && src[i+1] == '*':
				change("sBlockComment")
				skip()
				skip()
			case c == '/' && src[i+1] == '/':
				change("sLineComment")
				skip()
				skip()
			case matchByte("[0-9]", c):
				change("sNumber")
				eat()
			case matchByte(";", c):
				change("sSemi")
				eat()
				reset()
			case matchByte(letters, c):
				change("sName")
				eat()
			case matchByte(operators, c):
				change("sOperator")
				eat()
			case matchByte(special, c):
				change("sSpecial")
				eat()
				reset()
			case matchByte(stringParens, c):
				change("sString")
				eat()
			default:
				skip()
			}

		case "sBlockComment":
			if c == '*' && src[i+1] == '/' {
				reset()
			} else {
				skip()
			}

		case "sLineComment":
			if c == '\n' {
				reset()
			} else {
				skip()
			}

		case "sNumber":
			if matchByte("[0-9]", c) {
				eat()
			} else if matchByte("\\.", c) {
				if contains(lastLexeme, c) {
					change("sOperator")
					eat()
				} else {
					eat()
				}
			} else if matchByte(letters, c) {
				errorInvalidToken(c)
			} else {
				reset()
			}

		case "sOperator":
			if matchByte(operators, c) {
				eat()
			} else {
				reset()
			}

		case "sName":
			if matchByte("[A-Za-z0-9_$@]", c) {
				eat()
			} else {
				reset()
			}

		case "sString":
			eat()
			if matchByte(stringParens, c) {
				reset()
			}
		}
	}

	if state != "sNone" {
		reset()
	}

	return lexemes
}

type token struct {
	tType  int32
	lexeme string
}

const (
	tAS = iota
	tRETURN
	tFUNCTION
	tFROM
	tDEFAULT
	tEXPORT
	tIMPORT
	tVAR
	tCURLY_LEFT
	tCURLY_RIGHT
	tPAREN_LEFT
	tPAREN_RIGHT
	tBRACKET_LEFT
	tBRACKET_RIGHT
	tCOMMA
	tCOLON
	tASSIGN
	tSEMI
	tNUMBER
	tSTRING
	tNAME
	tCONST
	tLET
	tOF
	tELSE
	tIF
	tNEW
	tTRY
	tCATCH
	tWHILE
	tDO
)

func getTokens(lexemes [][]byte) []token {
	tokens := make([]token, 0)

	i := 0
	var lexeme string

	add := func(tType int32) {
		tokens = append(tokens, token{tType, lexeme})
		i++
	}

	for i := range lexemes {
		lexeme = string(lexemes[i])

		switch {
		case lexeme == "const":
			add(tCONST)
		case lexeme == "let":
			add(tLET)
		case lexeme == "do":
			add(tDO)
		case lexeme == "while":
			add(tWHILE)
		case lexeme == "new":
			add(tNEW)
		case lexeme == "catch":
			add(tCATCH)
		case lexeme == "try":
			add(tTRY)
		case lexeme == "else":
			add(tELSE)
		case lexeme == "if":
			add(tIF)
		case lexeme == "of":
			add(tOF)
		case lexeme == "as":
			add(tAS)
		case lexeme == "return":
			add(tRETURN)
		case lexeme == "function":
			add(tFUNCTION)
		case lexeme == "from":
			add(tFROM)
		case lexeme == "default":
			add(tDEFAULT)
		case lexeme == "export":
			add(tEXPORT)
		case lexeme == "import":
			add(tIMPORT)
		case lexeme == "var":
			add(tVAR)
		case lexeme == "{":
			add(tCURLY_LEFT)
		case lexeme == "}":
			add(tCURLY_RIGHT)
		case lexeme == "(":
			add(tPAREN_LEFT)
		case lexeme == ")":
			add(tPAREN_RIGHT)
		case lexeme == "[":
			add(tBRACKET_LEFT)
		case lexeme == "]":
			add(tBRACKET_RIGHT)
		case lexeme == ",":
			add(tCOMMA)
		case lexeme == ":":
			add(tCOLON)
		case lexeme == "=":
			add(tASSIGN)
		case lexeme == ";":
			add(tSEMI)
		case matchString("^[0-9]", lexeme):
			add(tNUMBER)
		case matchString("\\'|\"|\\`", lexeme):
			add(tSTRING)
		default:
			add(tNAME)
		}
	}

	return tokens
}

var keywords = []int{
	tAS,
	tRETURN,
	tFUNCTION,
	tFROM,
	tDEFAULT,
	tEXPORT,
	tIMPORT,
	tVAR,
	tCONST,
	tLET,
	tOF,
	tELSE,
	tIF,
	tNEW,
	tTRY,
	tCATCH,
	tWHILE,
	tDO,
}

func isKeyword(t token) bool {
	return containsInt(keywords, int(t.tType))
}

type importInfo struct {
	file       string
	vars       []string
	defaultVar string
}

type functionInfo struct {
	args          []string
	startBlockLvl int
}

func parse(tokens []token, fileName string) []byte {
	body := make([]byte, 0)

	i := 0

	exportedVars := ""
	exportedDefault := ""

	getToken := func(index int) token {
		if index >= 0 && index < len(tokens) {
			return tokens[index]
		}
		return token{}
	}

	//assumes current token index!!!
	writeToken := func(t token) {
		if isKeyword(t) && (getToken(i-1).tType == tNAME || getToken(i-1).tType == tNUMBER) {
			body = append(body, ' ')
		}
		body = append(body, []byte(t.lexeme)...)
		if isKeyword(t) && (getToken(i+1).tType == tNAME || getToken(i+1).tType == tNUMBER || isKeyword(getToken(i+1))) {
			body = append(body, ' ')
		}
	}

	// writeNext := func() {
	// 	writeToken(tokens[i])
	// 	i++
	// }

	// skip := func() {}

	result := make([]byte, 0)

	formattedPath := formatImportPath(fileName)
	body = append(body, []byte(fmt.Sprintf("var %s=(function(){", formattedPath))...)

	fileImports := make([]*importInfo, 0)

	functionStack := make([]functionInfo, 0)

	blockLvl := 0

	for i < len(tokens) {
		if tokens[i].tType == tIMPORT {
			//imports
			imported := importInfo{}

			if tokens[i+1].tType == tNAME {
				imported.defaultVar = tokens[i+1].lexeme
				i++
			}

			for i++; tokens[i].tType != tSEMI; i++ {
				if tokens[i].tType == tNAME {
					imported.vars = append(imported.vars, tokens[i].lexeme)
				} else if tokens[i].tType == tDEFAULT {
					i += 2
					imported.defaultVar = tokens[i].lexeme
				}
			}

			importText := ""

			importPath := resolveImportPath(tokens[i-1].lexeme, fileName)
			formattedPath := formatImportPath(importPath)
			imported.file = formattedPath

			ext := getExtension(importPath)
			if ext != "js" {
				from, err := os.Open(importPath)
				if err != nil {
					fmt.Println(err)
				}
				defer from.Close()

				fullFileName := formattedPath + "." + ext
				to, err := os.OpenFile("test/build/"+fullFileName, os.O_RDWR|os.O_CREATE, 0666)
				if err != nil {
					fmt.Println(err)
				}
				defer to.Close()

				io.Copy(to, from)

				importText += fmt.Sprintf("const %v='%v';", imported.defaultVar, fullFileName)
			} else {
				src, _ := ioutil.ReadFile(importPath)
				result = append(result, parse(getTokens(getLexemes([]byte(src))), importPath)...)
				fileImports = append(fileImports, &imported)
			}

			body = append(body, []byte(importText)...)

		} else if tokens[i].tType == tEXPORT {
			// exports
			for i++; tokens[i].tType != tSEMI; i++ {
				if tokens[i].tType == tNAME {
					exportedVars += tokens[i].lexeme + ","
					writeToken(tokens[i])
				} else if tokens[i].tType == tDEFAULT {
					for i++; tokens[i].tType != tSEMI; i++ {
						exportedDefault += tokens[i].lexeme
					}
					break
				} else {
					writeToken(tokens[i])
				}
			}
			writeToken(tokens[i])
		} else if tokens[i].tType == tFUNCTION {
			info := functionInfo{}
			info.args = make([]string, 0)
			info.startBlockLvl = blockLvl

			for tokens[i].tType != tPAREN_LEFT {
				writeToken(tokens[i])
				i++
			}
			for ; tokens[i].tType != tPAREN_RIGHT; i++ {
				writeToken(tokens[i])
				if tokens[i].tType == tNAME {
					info.args = append(info.args, tokens[i].lexeme)
				}
			}
			functionStack = append(functionStack, info)
			writeToken(tokens[i])
		} else if tokens[i].tType == tCURLY_LEFT {
			writeToken(tokens[i])
			blockLvl++
		} else if tokens[i].tType == tCURLY_RIGHT {
			writeToken(tokens[i])
			blockLvl--

			if len(functionStack) > 0 {
				currentFunction := functionStack[len(functionStack)-1]
				if currentFunction.startBlockLvl == blockLvl {
					functionStack = functionStack[:len(functionStack)-1]
				}
			}
		} else {
			if tokens[i].tType == tNAME {
				newLexeme := tokens[i].lexeme

				importRedeclared := false
				for _, fInfo := range functionStack {
					if containsStr(fInfo.args, tokens[i].lexeme) {
						importRedeclared = true
						break
					}
				}

				if !importRedeclared {
					for _, imp := range fileImports {
						if imp.defaultVar == tokens[i].lexeme {
							newLexeme = imp.file + ".default"
						} else {
							for _, importVar := range imp.vars {
								if importVar == tokens[i].lexeme {
									newLexeme = imp.file + "." + importVar
								}
							}
						}
					}
				}

				writeToken(token{tNAME, newLexeme})
			} else {
				writeToken(tokens[i])
			}
		}

		i++
	}

	exportText := ""
	if exportedDefault != "" {
		exportText = fmt.Sprintf("return {default:%v,%v};})();", exportedDefault, exportedVars)
	} else {
		exportText = fmt.Sprintf("return {%v};})();", exportedVars)
	}
	body = append(body, []byte(exportText)...)

	return append(result, body...)
}

func main() {
	srcFileName := "test/index.js"
	src, _ := ioutil.ReadFile(srcFileName)

	result := parse(getTokens(getLexemes([]byte(src))), srcFileName)

	ioutil.WriteFile("test/build/bundle.js", result, 777)
}
