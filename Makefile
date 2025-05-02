PKGROOT=github.com/sstark/gjfy
VERSION_VAR=cmd.Version
VERSION=$(shell git describe --match="v*")
BIN_NAME=gjfy

all: build

build:
	go build -ldflags "-s -w -X $(PKGROOT)/$(VERSION_VAR)=$(VERSION)" -o $(BIN_NAME)

build-debug:
	go build -ldflags "-X $(PKGROOT)/$(VERSION_VAR)=$(VERSION)-debug" -o $(BIN_NAME)

clean:
	go clean

test:
	go test ./...

test-verbose:
	go test -v ./...
