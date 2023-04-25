.PHONY: build run test

build:
	go build -o app

run:
	go run main.go

test:
	go test ./...

docker-build:
	docker build -t ghcr.io/ericmort/naas/naas:latest .

docker-push:
	docker push ghcr.io/ericmort/naas/naas:latest