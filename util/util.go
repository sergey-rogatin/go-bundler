package util

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func ClearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func IndexOf(arr []string, str string) int {
	for i, item := range arr {
		if item == str {
			return i
		}
	}
	return -1
}

const (
	C_RESET  = "\033[0;0m"
	C_BLACK  = "\033[0;30m"
	C_RED    = "\033[0;31m"
	C_GREEN  = "\033[0;32m"
	C_YELLOW = "\033[0;33m"
	C_BLUE   = "\033[0;34m"
	C_PURPLE = "\033[0;35m"
	C_CYAN   = "\033[0;36m"
	C_WHITE  = "\033[0;37m"
)

func Cprintf(color, format string, a ...interface{}) {
	format = color + format + C_RESET
	fmt.Printf(format, a...)
}

type SafeFile struct {
	file *os.File
	lock sync.RWMutex
}

func NewSafeFile(fileName string) *SafeFile {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	return &SafeFile{file, sync.RWMutex{}}
}

func (sf *SafeFile) Write(data []byte) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	sf.file.Write(data)
}

func (sf *SafeFile) Close() {
	sf.file.Close()
}
