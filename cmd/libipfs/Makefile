all: install
commit = $(shell git rev-parse --short HEAD 2> /dev/null || echo unknown)
ldflags = "-X "github.com/ipfs/go-ipfs/repo/config".CurrentCommit=$(commit) -w -s"

build:
	go build -ldflags=$(ldflags) -buildmode=c-archive -x

install: build
	go install
