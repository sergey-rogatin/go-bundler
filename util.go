package main

import (
	"fmt"
)

const (
	c_BLACK  = "\033[0;30m"
	c_RED    = "\033[0;31m"
	c_GREEN  = "\033[0;32m"
	c_YELLOW = "\033[0;33m"
	c_BLUE   = "\033[0;34m"
	c_PURPLE = "\033[0;35m"
	c_CYAN   = "\033[0;36m"
	c_WHITE  = "\033[0;37m"
)

func cprintf(color string, format string, a ...interface{}) {
	format = color + format + c_WHITE
	fmt.Printf(format, a...)
}
