all: run

fmt:
	go fmt

# TODO do this better
build: fmt
	cd cmd && go build -o ../bin/restTest

run: build
	./bin/restTest

test: fmt
	go test -v -race -coverprofile coverage.out -tags test
	go tool cover -func=./coverage.out

clean: fmt
	rm -f ./coverage.out

.PHONY: build run test clean
