package common

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// ErrExit allows to exit on error with exit code 1 after printing error message
// NB: specific function assigned to a variable for mock in unit tests.
var ErrExit = func(msg interface{}) {
	if msg != nil {
		fmt.Println(msg)
	}
	os.Exit(1)
}

func GetUsableColWidth(colNumber int) int {
	defaultTTYWidth := 100
	width, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("[WARNING] Failed to get terminal width due to err:%s\n", err.Error())
		return defaultTTYWidth
	}
	// remove chars for separators according to the nb of columns (colNumber+1) and add configurable margin
	return width - (colNumber + 5)
}

func GetRatio(value int, ratio float64) int {
	return int(float64(value) * ratio)
}
