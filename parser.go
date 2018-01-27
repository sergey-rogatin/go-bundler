package main

import (
	"fmt"
	"strings"
)

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
					if t.tType == tNAME {
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
