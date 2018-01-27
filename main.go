package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type safeFile struct {
	file    *os.File
	isReady chan bool
}

func newSafeFile(fileName string) *safeFile {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	isReady := make(chan bool, 1)
	isReady <- true

	return &safeFile{
		file,
		isReady,
	}
}

func closeSafeFile(sf *safeFile) {
	sf.file.Close()
}

func containsStr(arr []string, c string) bool {
	for _, b := range arr {
		if b == c {
			return true
		}
	}
	return false
}

func trimQuotesFromString(jsString string) string {
	return jsString[1 : len(jsString)-1]
}

func getExtension(resolvedImportPath string) string {
	parts := strings.Split(resolvedImportPath, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "js"
}

func createVarNameFromPath(resolvedImportPath string) string {
	newName := strings.Replace(resolvedImportPath, "/", "_", -1)
	newName = strings.Replace(newName, ".", "_", -1)
	return newName
}

func isKeyword(t token) bool {
	_, ok := keywords[t.lexeme]
	return ok && t.tType != tNAME
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
	bundleSf *safeFile,
	finishedImportsCh chan string,
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

		addFilesToBundle(unbundledFiles, allImportedPaths, bundleSf)
		writeTokensToFile(result, bundleSf)

	default:
		dstFileName := filepath.Dir(bundleSf.file.Name()) + "/" + createVarNameFromPath(resolvedPath) + "." + ext
		copyFile(dstFileName, resolvedPath)
	}
	finishedImportsCh <- resolvedPath
}

func addFilesToBundle(
	files []string,
	allImportedPaths *[]string,
	bundleSf *safeFile,
) {
	filesImportedCh := make(chan string, len(files))

	for _, unbundledFile := range files {
		go addFileToBundle(unbundledFile, allImportedPaths, bundleSf, filesImportedCh)
	}

	for counter := 0; counter < len(files); counter++ {
		//fmt.Printf("Finished bundling %s\n", <-filesImportedCh)
		<-filesImportedCh
	}
}

func writeTokensToFile(tokens []token, sf *safeFile) {
	<-sf.isReady
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

		sf.file.Write(toWrite)
	}
	sf.isReady <- true
}

func createBundle(entryFileName, bundleFileName string) []string {
	buildStartTime := time.Now()

	allImportedPaths := make([]string, 0)

	os.Remove(bundleFileName)
	sf := newSafeFile(bundleFileName)
	defer closeSafeFile(sf)

	addFilesToBundle([]string{entryFileName}, &allImportedPaths, sf)

	fmt.Printf("Build finished in %s\n", time.Since(buildStartTime))
	return allImportedPaths
}

func main() {
	entryName := "test/index.js"
	bundleName := "test/build/bundle.js"
	bundleHTMLTemplate("test/template.html", bundleName)
	allImportedPaths := createBundle(entryName, bundleName)

	// dev server
	watchBundledFiles(allImportedPaths, entryName, bundleName)

	fmt.Println("Dev server listening at port 8080")
	server := http.FileServer(http.Dir("test/build"))
	err := http.ListenAndServe(":8080", server)
	log.Fatal(err)
}
