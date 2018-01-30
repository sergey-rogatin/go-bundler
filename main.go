package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type safeFile struct {
	file *os.File
	lock sync.RWMutex
}

func newSafeFile(fileName string) *safeFile {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	return &safeFile{file, sync.RWMutex{}}
}

func (sf *safeFile) close() {
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

type cachedFile struct {
	data        []byte
	lastModTime time.Time
	imports     []string
}

type bundleCache struct {
	files map[string]cachedFile
	lock  sync.RWMutex
}

func (c *bundleCache) read(fileName string) (cachedFile, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	file, ok := c.files[fileName]
	return file, ok
}

func (c *bundleCache) write(fileName string, data cachedFile) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.files[fileName] = data
}

func addFileToBundle(
	resolvedPath string,
	bundleSf *safeFile,
	finishedImportsCh chan string,
	cache *bundleCache,
) {
	ext := getExtension(resolvedPath)
	fileStats, _ := os.Stat(resolvedPath)

	var data []byte
	var fileImports []string

	switch ext {
	case "js":
		cachedData, ok := cache.read(resolvedPath)
		if ok && cachedData.lastModTime == fileStats.ModTime() {
			data = cachedData.data
			fileImports = cachedData.imports
		} else {
			src, err := ioutil.ReadFile(resolvedPath)
			if err != nil {
				panic(err)
			}

			data, fileImports = loadJsFile(src, resolvedPath)
		}

	default:
		bundleDir := filepath.Dir(bundleSf.file.Name())
		dstFileName := bundleDir + "/" + createVarNameFromPath(resolvedPath) + "." + ext
		copyFile(dstFileName, resolvedPath)
	}

	cache.write(resolvedPath, cachedFile{
		data:        data,
		lastModTime: fileStats.ModTime(),
		imports:     fileImports,
	})

	addFilesToBundle(fileImports, bundleSf, cache)
	writeToSafeFile(data, bundleSf)
	finishedImportsCh <- resolvedPath
}

func addFilesToBundle(
	files []string,
	bundleSf *safeFile,
	cache *bundleCache,
) {
	filesImportedCh := make(chan string, len(files))

	for _, unbundledFile := range files {
		go addFileToBundle(unbundledFile, bundleSf, filesImportedCh, cache)
	}

	for counter := 0; counter < len(files); counter++ {
		<-filesImportedCh
	}
}

func writeToSafeFile(data []byte, sf *safeFile) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	sf.file.Write(data)
}

func createBundle(entryFileName, bundleFileName string, cache *bundleCache) {
	buildStartTime := time.Now()

	os.MkdirAll(filepath.Dir(bundleFileName), 0666)
	os.Remove(bundleFileName)
	sf := newSafeFile(bundleFileName)
	defer sf.close()

	addFilesToBundle([]string{entryFileName}, sf, cache)

	fmt.Printf("Build finished in %s\n", time.Since(buildStartTime))
}

type configJSON struct {
	Entry        string
	TemplateHTML string
	BundleDir    string
	WatchFiles   bool
	DevServer    struct {
		Enable bool
		Port   int
	}
}

func main() {
	src, _ := ioutil.ReadFile("test.js")

	tokens := lex(src)
	//fmt.Println(tokens)
	parse(tokens)

	// config := configJSON{}
	// config.TemplateHTML = "test/template.html"
	// config.WatchFiles = false
	// config.DevServer.Enable = false

	// if len(os.Args) > 1 {
	// 	configFileName := os.Args[1]

	// 	configFile, err := ioutil.ReadFile(configFileName)
	// 	if err != nil {
	// 		fmt.Println("Unable to load config file!")
	// 		config = configJSON{}
	// 	}
	// 	json.Unmarshal(configFile, &config)
	// }

	// // config defaults
	// if config.Entry == "" {
	// 	config.Entry = "test/index.js"
	// }
	// if config.BundleDir == "" {
	// 	config.BundleDir = "test/build"
	// }
	// if config.DevServer.Port == 0 {
	// 	config.DevServer.Port = 8080
	// }

	// entryName := config.Entry
	// bundleName := filepath.Join(config.BundleDir, "bundle.js")

	// cache := bundleCache{}
	// cache.files = make(map[string]cachedFile)
	// createBundle(entryName, bundleName, &cache)

	// if config.TemplateHTML != "" {
	// 	bundleHTMLTemplate(config.TemplateHTML, bundleName)
	// }

	// // dev server and watching files
	// if config.DevServer.Enable {
	// 	if config.WatchFiles {
	// 		go watchBundledFiles(&cache, entryName, bundleName)
	// 	}
	// 	fmt.Printf("Dev server listening at port %v\n", config.DevServer.Port)
	// 	server := http.FileServer(http.Dir(config.BundleDir))
	// 	err := http.ListenAndServe(fmt.Sprintf(":%v", config.DevServer.Port), server)
	// 	log.Fatal(err)
	// } else if config.WatchFiles {
	// 	watchBundledFiles(&cache, entryName, bundleName)
	// }
}
