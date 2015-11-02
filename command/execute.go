package command

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os/exec"
	"regexp"
)

func Execute(cmdName string, cmdArgs []string) {
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		Error("Error creating StdoutPip for Cmd", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			Write(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		Error("Error starting Cmd", err)
	}

	err = cmd.Wait()
	if err != nil {
		Error("Error waiting for Cmd", err)
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
