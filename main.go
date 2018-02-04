package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lvl5hm/goBundler/jsLoader"
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

func (sf *safeFile) write(data []byte) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	sf.file.Write(data)
}

func (sf *safeFile) close() {
	sf.file.Close()
}

type fileCache struct {
	data        []byte
	lastModTime time.Time
	imports     []string
	isReachable bool
}

type bundleCache struct {
	files map[string]fileCache
	lock  sync.RWMutex
}

func (c *bundleCache) read(fileName string) (fileCache, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	file, ok := c.files[fileName]
	return file, ok
}

func (c *bundleCache) write(fileName string, data fileCache) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.files[fileName] = data
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
	// setting up config
	config := configJSON{}

	configFileName := "config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}

	configFile, err := ioutil.ReadFile(configFileName)
	if err != nil {
		fmt.Println("Unable to load config file!")
	} else {
		json.Unmarshal(configFile, &config)
	}

	// config defaults
	if config.Entry == "" {
		config.Entry = "index.js"
	}
	if config.BundleDir == "" {
		config.BundleDir = "build"
	}
	if config.DevServer.Port == 0 {
		config.DevServer.Port = 8080
	}

	// creating bundle
	bundleName := filepath.Join(config.BundleDir, "bundle.js")

	cache := bundleCache{}
	cache.files = make(map[string]fileCache)
	createBundle(config.Entry, bundleName, &cache)

	if config.TemplateHTML != "" {
		bundleHTMLTemplate(config.TemplateHTML, bundleName)
	}

	// dev server and watching files
	if config.DevServer.Enable {
		if config.WatchFiles {
			go watchBundledFiles(&cache, config.Entry, bundleName)
		}
		fmt.Printf("Dev server listening at port %v\n", config.DevServer.Port)
		server := http.FileServer(http.Dir(config.BundleDir))
		err := http.ListenAndServe(fmt.Sprintf(":%v", config.DevServer.Port), server)
		log.Fatal(err)
	} else if config.WatchFiles {
		watchBundledFiles(&cache, config.Entry, bundleName)
	}
}

func createBundle(entryFileName, bundleFileName string, cache *bundleCache) {
	buildStartTime := time.Now()

	os.MkdirAll(filepath.Dir(bundleFileName), 0666)
	os.Remove(bundleFileName)
	sf := newSafeFile(bundleFileName)
	defer sf.close()

	// mark all files as unreachable at the start of the build
	// so the autorebuilder does not try to rebuild when they change
	for fileName, file := range cache.files {
		file.isReachable = false
		cache.files[fileName] = file
	}

	err := addFilesToBundle([]string{entryFileName}, sf, cache)
	if err == nil {
		fmt.Printf("Build finished in %s\n", time.Since(buildStartTime))
	} else {
		fmt.Println(err)
	}
}

func addFilesToBundle(
	files []string,
	bundleSf *safeFile,
	cache *bundleCache,
) error {
	errorCh := make(chan error, len(files))

	for _, unbundledFile := range files {
		go addFileToBundle(unbundledFile, bundleSf, errorCh, cache)
	}

	for counter := 0; counter < len(files); counter++ {
		err := <-errorCh
		if err != nil {
			return err
		}
	}

	return nil
}

func addFileToBundle(
	resolvedPath string,
	bundleSf *safeFile,
	errorCh chan error,
	cache *bundleCache,
) {
	ext := filepath.Ext(resolvedPath)
	fileStats, _ := os.Stat(resolvedPath)

	var data []byte
	var fileImports []string

	switch ext {
	case ".js":
		cachedFile, ok := cache.read(resolvedPath)
		if ok && cachedFile.lastModTime == fileStats.ModTime() {
			data = cachedFile.data
			fileImports = cachedFile.imports
		} else {
			src, err := ioutil.ReadFile(resolvedPath)
			if err != nil {
				panic(err)
			}

			data, fileImports, err = jsLoader.LoadFile(src, resolvedPath)
			if err != nil {
				errorCh <- err
			}
		}

	default:
		bundleDir := filepath.Dir(bundleSf.file.Name())
		dstFileName := bundleDir + "/" + jsLoader.CreateVarNameFromPath(resolvedPath) + ext
		copyFile(dstFileName, resolvedPath)
	}

	cache.write(resolvedPath, fileCache{
		data:        data,
		lastModTime: fileStats.ModTime(),
		imports:     fileImports,
		isReachable: true,
	})

	err := addFilesToBundle(fileImports, bundleSf, cache)
	bundleSf.write(data)
	errorCh <- err
}

func bundleHTMLTemplate(templateName, bundleName string) {
	template, err := ioutil.ReadFile(templateName)
	if err != nil {
		log.Fatal("Can't find or open html template")
	}

	templateStr := string(template)
	insertIndex := strings.Index(templateStr, "\n</body")
	if insertIndex < 0 {
		log.Fatal("Can't find end of <body> in html template")
	}

	result := templateStr[:insertIndex] +
		"\n  <script src=\"" + filepath.Base(bundleName) + "\"></script>\n" +
		templateStr[insertIndex+1:]

	bundleDir := filepath.Dir(bundleName)
	ioutil.WriteFile(filepath.Join(bundleDir, "index.html"), []byte(result), 0666)
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
