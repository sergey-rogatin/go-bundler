package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func containsStr(arr []string, c string) bool {
	for _, b := range arr {
		if b == c {
			return true
		}
	}
	return false
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

func getExtension(resolvedImportPath string) string {
	parts := strings.Split(resolvedImportPath, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "js"
}

func formatImportPath(resolvedImportPath string) string {
	newName := strings.Replace(resolvedImportPath, "/", "_", -1)
	newName = strings.Replace(newName, ".", "_", -1)
	return newName
}

func isKeyword(t token) bool {
	_, ok := keywords[t.lexeme]
	return ok && t.tType != tNAME
}

type jsImportInfo struct {
	exportObjName string
	resolvedPath  string
	def           string
	vars          []string
}

type jsExportInfo struct {
	def  []token
	vars []token
}

func transformIntoModule(tokens []token, resolvedPath string) ([]token, []jsImportInfo) {
	result := make([]token, 0, len(tokens)+128)
	imports := make([]jsImportInfo, 0)

	if len(tokens) == 0 {
		return result, imports
	}

	fileExports := jsExportInfo{}

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

	exportObjName := formatImportPath(resolvedPath)
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
			jsImport := jsImportInfo{}
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
			importPath := resolveImportPath(t.lexeme, resolvedPath)
			skip() // "./foo"
			skip() // ;

			ext := getExtension(importPath)
			formattedName := formatImportPath(importPath)

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

func copyFile(dst, src string) {
	from, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer to.Close()
	io.Copy(to, from)
}

func addFileToBundle(
	resolvedPath string,
	allImportedPaths *[]string,
	bundleFile *os.File,
	bundleFileReadyCh chan bool,
	importsFinishedCh chan string,
) {
	*allImportedPaths = append(*allImportedPaths, resolvedPath)

	ext := getExtension(resolvedPath)

	switch ext {
	case "js":
		src, err := ioutil.ReadFile(resolvedPath)
		if err != nil {
			panic(err)
		}
		tokens := lex(src)

		result, fileImports := transformIntoModule(tokens, resolvedPath)

		unbundledFiles := make([]string, 0)
		for _, imp := range fileImports {
			if !containsStr(*allImportedPaths, imp.resolvedPath) {
				unbundledFiles = append(unbundledFiles, imp.resolvedPath)
			}
		}

		addFilesToBundleAsync(unbundledFiles, allImportedPaths, bundleFile, bundleFileReadyCh)
		writeToFile(result, bundleFile, bundleFileReadyCh)

	default:
		dstFileName := filepath.Dir(bundleFile.Name()) + "/" + formatImportPath(resolvedPath) + "." + ext
		copyFile(dstFileName, resolvedPath)
	}
	importsFinishedCh <- resolvedPath
}

func writeToFile(tokens []token, file *os.File, isFileReady chan bool) {
	<-isFileReady
	for i, t := range tokens {
		tIsKeyword := isKeyword(t)

		toWrite := make([]byte, 0)
		if tIsKeyword && i > 0 && (tokens[i-1].tType == tNAME || tokens[i-1].tType == tNUMBER || isKeyword(tokens[i-1])) {
			toWrite = append(toWrite, ' ')
		}
		toWrite = append(toWrite, []byte(t.lexeme)...)
		if tIsKeyword && i < len(tokens) && (tokens[i+1].tType == tNAME || tokens[i+1].tType == tNUMBER) {
			toWrite = append(toWrite, ' ')
		}

		file.Write(toWrite)
	}
	isFileReady <- true
}

func addFilesToBundleAsync(
	files []string,
	allImportedPaths *[]string,
	bundleFile *os.File,
	bundleFileReadyCh chan bool,
) {
	filesImportedCh := make(chan string, len(files))

	for _, unbundledFile := range files {
		go addFileToBundle(unbundledFile, allImportedPaths, bundleFile, bundleFileReadyCh, filesImportedCh)
	}

	for counter := 0; counter < len(files); counter++ {
		fmt.Printf("Finished bundling %s\n", <-filesImportedCh)
	}
}

func createBundle(entryFileName, bundleFileName string) {
	allImportedPaths := make([]string, 0)

	os.Remove(bundleFileName)
	bundleFile, err := os.Create(bundleFileName)
	if err != nil {
		panic(err)
	}
	defer bundleFile.Close()

	bundleFileReadyCh := make(chan bool, 1)
	bundleFileReadyCh <- true

	addFilesToBundleAsync([]string{entryFileName}, &allImportedPaths, bundleFile, bundleFileReadyCh)
}

func main() {
	createBundle("test/index.js", "test/build/bundle.js")
}
