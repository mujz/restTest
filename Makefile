all: get run

get:
	go get ./...

fmt:
	go fmt ./...

build: fmt
	go build -o ./bin/restTest

run: build
	./bin/restTest

test: fmt
	go test -v -race -coverprofile coverage.out *.go
	go tool cover -func=./coverage.out

clean: fmt
	rm -f ./coverage.out

.PHONY: get build run test clean
