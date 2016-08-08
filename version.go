package main

import (
	"fmt"
)

const ProjectName = "Terraform Exec"
const Usage = `Execute terraform commands across environments maintaining state in s3

Default project layout

  ./config/<environment>`

// The git commit that at compile time. Use -ldflags "-X main.GitCommit `git rev-parse HEAD`"
var GitCommit string
func gitCommit() string {
	if GitCommit != "" {
		return GitCommit
	}
	return "unknown"
}

func CommitMessage() {
	fmt.Printf("\nGit commit: %s\n\n", gitCommit())
}

// The version number.
const Version = "0.0.6"

