BUILD_OUTPUT=../bin/restTest
BUILD_DIR=cmd
GOARCH=amd64

all: run

fmt:
	go fmt

install:
	cd ${BUILD_DIR}; \
		go build -o ${GOBIN}/restTest

build: fmt
	cd ${BUILD_DIR}; \
		go build -o ${BUILD_OUTPUT}

linux:
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build -o ${BUILD_OUTPUT}-linux-${GOARCH}

darwin:
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BUILD_OUTPUT}-darwin-${GOARCH}

windows:
	cd ${BUILD_DIR}; \
	GOOS=windows GOARCH=${GOARCH} go build -o ${BUILD_OUTPUT}-windows-${GOARCH}.exe

run: build
	./bin/restTest

test: fmt
	go test -v -race -coverprofile coverage.out
	go tool cover -func=./coverage.out

clean: fmt
	rm -f ./coverage.out

.PHONY: build run test clean
