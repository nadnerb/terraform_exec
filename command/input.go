package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// does this work on windows?
func InputAffirmative() bool {
	fmt.Print("Are you sure? (yes)\n")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text) == "yes"
}
