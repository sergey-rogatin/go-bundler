package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

func watchBundledFiles(allImportedPaths []string, entryName, bundleName string) func() {
	resetModTime := func() []time.Time {
		modTime := make([]time.Time, 0, len(allImportedPaths))
		for _, file := range allImportedPaths {
			stats, _ := os.Stat(file)
			modTime = append(modTime, stats.ModTime())
		}
		return modTime
	}

	modTime := resetModTime()
	running := true

	checkFiles := func() {
		for running {
			for i, file := range allImportedPaths {
				stats, _ := os.Stat(file)
				if modTime[i] != stats.ModTime() {
					allImportedPaths = createBundle(entryName, bundleName)
					modTime = resetModTime()
					break
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	go checkFiles()

	return func() {
		fmt.Println("Stopped watching files")
		running = false
	}
}
