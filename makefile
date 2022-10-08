.PHONY: httpserver grpcserver build lint generate-proto lint-proto

httpserver:
	go build -o httpserver cmd/httpserver/main.go

grpcserver:
	go build -o grpcserver cmd/grpcserver/main.go

build: httpserver grpcserver

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.50.0

lint: bin/golangci-lint
	go fmt ./...
	go vet ./...
	bin/golangci-lint -c .golangci.yaml run ./...
	go mod tidy

generate-proto:
	buf generate

lint-proto:
	buf lint
