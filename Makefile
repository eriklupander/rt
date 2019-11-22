all: fmt build test

build:
	go build -o bin/rt main.go

run: build
	./bin/rt

test:
	go test ./... -count=1 -v

fmt:
	go fmt ./...