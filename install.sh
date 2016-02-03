#!/bin/bash
go install -ldflags "-X main.GitCommit `git rev-parse HEAD`"
