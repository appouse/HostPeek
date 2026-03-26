.PHONY: build run clean test vet

VERSION ?= dev

## Build the binary for the current platform
build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o hostpeek ./cmd/hostpeek

## Run the application
run: build
	./hostpeek

## Run go vet
vet:
	go vet ./...

## Run tests
test:
	go test ./...

## Clean build artifacts
clean:
	rm -f hostpeek hostpeek.exe

## Build for all platforms (local cross-compile)
build-all:
	GOOS=linux   GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o dist/hostpeek-linux-amd64       ./cmd/hostpeek
	GOOS=linux   GOARCH=arm64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o dist/hostpeek-linux-arm64       ./cmd/hostpeek
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o dist/hostpeek-windows-amd64.exe ./cmd/hostpeek

## Create a snapshot release (local, no publish)
snapshot:
	goreleaser release --snapshot --clean
