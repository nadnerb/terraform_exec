package command

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/fatih/color"
)

func Execute(cmdName string, cmdArgs []string) {
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		Error("Error creating StdoutPipe for Cmd", err)
	}

	cmdErrorReader, err := cmd.StderrPipe()
	if err != nil {
		Error("Error creating StderrPipe for Cmd", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			Write(scanner.Text())
		}
	}()

	errorScanner := bufio.NewScanner(cmdErrorReader)
	go func() {
		for errorScanner.Scan() {
			Write(errorScanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		Error("Error starting Cmd", err)
	}

	err = cmd.Wait()
	if err != nil {
		Warn("Error waiting for Cmd", err)
	}
}

var ok = regexp.MustCompile(`^ok:.*`)
var skipping = regexp.MustCompile(`^skipping:.*`)
var changed = regexp.MustCompile(`^changed:.*`)
var fatal = regexp.MustCompile(`^fatal:.*`)

func Write(text string) {
	if ok.MatchString(text) {
		color.Green(text)
	} else if skipping.MatchString(text) {
		color.Cyan(text)
	} else if changed.MatchString(text) {
		color.Yellow(text)
	} else if fatal.MatchString(text) {
		color.Red(text)
	} else {
		fmt.Printf("%s\n", text)
	}
}
