package main

import (
	"fmt"
	"os"
	"time"
)

func watchBundledFiles(cache *bundleCache, entryName, bundleName string) func() {
	fmt.Println("Watching for file changes")

	running := true

	checkFiles := func() {
		for running {
			for path, file := range cache.files {
				stats, _ := os.Stat(path)
				if file.lastModTime != stats.ModTime() {
					createBundle(entryName, bundleName, cache)
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
