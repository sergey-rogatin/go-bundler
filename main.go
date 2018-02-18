package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lvl5hm/go-bundler/jsLoader"
	"github.com/lvl5hm/go-bundler/urlLoader"
	"github.com/lvl5hm/go-bundler/util"
)

/* TODO:
better error handling and printing
multiple entry points?
multiple bundles per file type?

*/

type fileCache struct {
	Data        []byte
	LastModTime time.Time
	Imports     []string
	IsReachable bool
}

type bundleCache struct {
	Files   map[string]fileCache
	DirName string
	Lock    sync.RWMutex
}

func (c *bundleCache) read(fileName string) (fileCache, bool) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	file, ok := c.Files[fileName]
	return file, ok
}

func (c *bundleCache) write(fileName string, data fileCache) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Files[fileName] = data
}

func (c *bundleCache) saveFile() {
	if c.DirName == "" {
		return
	}

	err := os.MkdirAll(c.DirName, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	saveFile, err := os.Create(c.DirName + "/cache")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer saveFile.Close()

	enc := gob.NewEncoder(saveFile)
	err = enc.Encode(c.Files)
	if err != nil {
		fmt.Println("Error: cannot save cache to file!")
	}
}

func (c *bundleCache) loadFile() {
	saveFile, err := os.Open(c.DirName + "/cache")
	if err != nil {
		return
	}
	defer saveFile.Close()

	dec := gob.NewDecoder(saveFile)

	var files map[string]fileCache
	err = dec.Decode(&files)
	if err != nil {
		fmt.Println("Error: cache file is corrupted!")
		return
	}

	c.Files = files
}

type configJSON struct {
	Entry        string
	TemplateHTML string
	BundleDir    string
	WatchFiles   bool
	Env          map[string]string
	DevServer    struct {
		Enable bool
		Port   int
	}
	PermanentCache struct {
		Enable  bool
		DirName string
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
		util.Cprintf(util.C_YELLOW, "Warning: Unable to load config file!\n")
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
	if config.PermanentCache.DirName == "" {
		config.PermanentCache.DirName = ".go-bundler-cache"
	}

	// creating bundle
	bundleName := filepath.Join(config.BundleDir, "bundle.js")

	cache := &bundleCache{}
	if config.PermanentCache.Enable {
		cache.DirName = config.PermanentCache.DirName
	}

	cache.loadFile()
	if cache.Files == nil {
		cache.Files = map[string]fileCache{}
	}

	createBundle(config.Entry, bundleName, cache, &config)

	if config.TemplateHTML != "" {
		bundleHTMLTemplate(config.TemplateHTML, bundleName)
	}

	// dev server and watching files
	if config.DevServer.Enable {
		if config.WatchFiles {
			go watchBundledFiles(cache, config.Entry, bundleName, &config)
		}
		util.Cprintf(util.C_GREEN, "Dev server listening at port %v\n", config.DevServer.Port)
		server := http.FileServer(http.Dir(config.BundleDir))
		err := http.ListenAndServe(fmt.Sprintf(":%v", config.DevServer.Port), server)
		log.Fatal(err)
	} else if config.WatchFiles {
		watchBundledFiles(cache, config.Entry, bundleName, &config)
	}
}

func createBundle(entryFileName, bundleFileName string, cache *bundleCache, config *configJSON) {
	startTime := time.Now()
	os.MkdirAll(filepath.Dir(bundleFileName), 0666)
	os.Remove(bundleFileName)
	sf := util.NewSafeFile(bundleFileName)
	defer sf.Close()

	// mark all files as unreachable at the start of the build
	// so the autorebuilder does not try to rebuild when they change
	for fileName, file := range cache.Files {
		file.IsReachable = false
		cache.Files[fileName] = file
	}

	sf.Write(jsLoader.GetJsBundleFileHead())
	err := addFilesToBundle([]string{entryFileName}, sf, cache, config)

	importsMap := map[string][]string{}
	for path, file := range cache.Files {
		importsMap[path] = file.Imports
	}
	tail, warnings := jsLoader.GetJsBundleFileTail(entryFileName, importsMap)
	sf.Write(tail)
	buildTime := time.Since(startTime)

	if config.WatchFiles {
		util.ClearScreen()
	}

	if err == nil {
		if config.PermanentCache.Enable {
			cacheSaveStart := time.Now()
			cache.saveFile()
			cacheSaveTime := time.Since(cacheSaveStart)
			util.Cprintf(util.C_GREEN, "Cache saved to %s in %s\n", config.PermanentCache.DirName, cacheSaveTime)
		}

		if len(warnings) > 0 {
			for _, w := range warnings {
				util.Cprintf(util.C_YELLOW, "%s", w)
			}
		}
		util.Cprintf(util.C_GREEN, "Build finished in %s\n", buildTime)
	} else {
		util.Cprintf(util.C_RED, "%s\n", err)
		return
	}
}

func addFilesToBundle(
	files []string,
	bundleSf *util.SafeFile,
	cache *bundleCache,
	config *configJSON,
) error {
	errorCh := make(chan error, len(files))

	for _, unbundledFile := range files {
		go addFileToBundle(unbundledFile, bundleSf, errorCh, cache, config)
	}

	errs := []error{}
	for counter := 0; counter < len(files); counter++ {
		err := <-errorCh
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return multiError{errs}
}

type loaderError struct {
	fileName string
	err      error
}

func (l loaderError) Error() string {
	return fmt.Sprintf("Error loading file %s:\n%s", l.fileName, l.err)
}

func addFileToBundle(
	fileName string,
	bundleSf *util.SafeFile,
	errorCh chan error,
	cache *bundleCache,
	config *configJSON,
) {
	fileStats, err := os.Stat(fileName)
	if err != nil {
		errorCh <- loaderError{fileName, fmt.Errorf("Cannot find file")}
		return
	}

	var data []byte
	var fileImports []string
	var lastModTime = fileStats.ModTime()

	saveCache := func() {
		cache.write(fileName, fileCache{
			Data:        data,
			Imports:     fileImports,
			LastModTime: lastModTime,
			IsReachable: true,
		})
	}

	cache.Lock.Lock()
	cachedFile, ok := cache.Files[fileName]
	if ok && cachedFile.IsReachable {
		cache.Lock.Unlock()
		errorCh <- nil
		return
	}
	cache.Files[fileName] = fileCache{
		IsReachable: true,
	}
	cache.Lock.Unlock()

	if ok && cachedFile.LastModTime == fileStats.ModTime() && cachedFile.Data != nil {
		data = cachedFile.Data
		fileImports = cachedFile.Imports
	} else {
		ext := filepath.Ext(fileName)

		switch ext {
		case ".js":
			data, fileImports, err = jsLoader.LoadFile(fileName, config.Env)

		default:
			data, fileImports, err = urlLoader.LoadFile(fileName, config.BundleDir)
		}

		if err != nil {
			saveCache()
			errorCh <- loaderError{fileName, err}
			return
		}
	}

	bundleSf.Write(data)
	saveCache()
	multiErr := addFilesToBundle(fileImports, bundleSf, cache, config)
	errorCh <- multiErr
}

type multiError struct {
	errs []error
}

func (me multiError) Error() string {
	res := ""
	for i, e := range me.errs {
		res += e.Error()
		if i < len(me.errs)-1 {
			res += "\n"
		}
	}
	return res
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

func watchBundledFiles(
	cache *bundleCache,
	entryName,
	bundleName string,
	config *configJSON,
) func() {
	fmt.Println("Watching for file changes")

	running := true

	checkFiles := func() {
		for running {
			for path, file := range cache.Files {
				if !file.IsReachable {
					continue
				}

				stats, err := os.Stat(path)
				if err != nil {
					continue
				}
				if file.LastModTime != stats.ModTime() {
					createBundle(entryName, bundleName, cache, config)
					break
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	checkFiles()

	return func() {
		fmt.Println("Stopped watching files")
		running = false
	}
}
