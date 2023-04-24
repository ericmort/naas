.PHONY: build run test

build:
	go build -o app

run:
	go run main.go

test:
	go test ./...

