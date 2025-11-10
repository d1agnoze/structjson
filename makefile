.PHONY: all test clean

build-cli: 
	go build -o bin/cli ./cmd/cli
