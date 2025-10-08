test:
	go mod tidy
	go test -v

test_coverage:
	go mod tidy
	go test -v -coverprofile=coverage.out

install:
	go install

lint:
	golangci-lint run cmd/gendiff

build:
	go build -o bin/gendiff ./cmd/gendiff
