test:
	go mod tidy
	go test -v

install:
	go install

lint:
	golangci-lint run cmd/gendiff

build:
	go build -o bin/gendiff ./cmd/gendiff
