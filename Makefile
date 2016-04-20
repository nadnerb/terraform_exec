install:
	go get ./...
	go get github.com/stretchr/testify/assert
	go install

build:
	go get github.com/mitchellh/gox
	gox -osarch "linux/amd64 darwin/amd64 windows/amd64"
