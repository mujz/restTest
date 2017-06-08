BUILD_OUTPUT=../bin/restTest
BUILD_DIR=cmd
GOARCH=amd64

all: run

lint:
	golint

fmt:
	gofmt -s -w .

install:
	cd ${BUILD_DIR}; \
		go build -o ${GOBIN}/restTest

build: lint fmt
	cd ${BUILD_DIR}; \
		go build -o ${BUILD_OUTPUT}

linux: lint fmt
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build -o ${BUILD_OUTPUT}-linux-${GOARCH}

darwin: lint fmt
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BUILD_OUTPUT}-darwin-${GOARCH}

windows: lint fmt
	cd ${BUILD_DIR}; \
	GOOS=windows GOARCH=${GOARCH} go build -o ${BUILD_OUTPUT}-windows-${GOARCH}.exe

run: build
	./bin/restTest

test: fmt
	go test -v -race -coverprofile coverage.out
	go tool cover -func=./coverage.out

clean:
	rm -f ./coverage.out

.PHONY: lint fmt install build linux darwin windows run test clean
