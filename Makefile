NAME=spotify-playlists
VERSION=v1.0
HASH=$(shell git rev-parse --short HEAD)
PACKAGES:=$(shell go list ./... | grep -v /vendor/)
LDFLAGS:=-ldflags "-w -X main.AppVersion=${VERSION} -X main.GitCommit=${HASH}"
GOOS=linux
PKG:=github.com/...
PREFIX?=$(shell pwd)

fmt:
	@echo "+ $@"
	@gofmt -s -l . | grep -v '.pb.go:' | grep -v vendor | tee /dev/stderr

test:
	# go test $(packages) 
	go test ./...

test-verbose:
	# go test $(packages) 
	go test ./...

clean:
	rm -rf main main.zip

bin:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=$(GOOS) go build -a -tags netgo ${LDFLAGS} -o bin/$(GOOS)/main 

getVersion:
	@echo $(VERSION)-$(HASH)