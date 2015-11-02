package command

import (
	"fmt"
	"os"
	"github.com/fatih/color"
)

var cyan = color.New(color.FgCyan).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()

func Error(errorMessage string, error error) {
	fmt.Println()
	fmt.Println("---------------------------------------------")
	fmt.Fprintf(os.Stderr, "ERROR\n")
	fmt.Fprintf(os.Stderr, "%s: %s\n", cyan(errorMessage), red(error))
	fmt.Println("---------------------------------------------")
	fmt.Println()
	os.Exit(1)
}

func Warn(errorMessage string, error error) {
	fmt.Fprintf(os.Stderr, "WARNING\n")
	fmt.Fprintf(os.Stderr, "%s: %s\n", cyan(errorMessage), yellow(error))
	fmt.Println()
}

