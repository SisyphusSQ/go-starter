BINARY_NAME = go_starter

VARS_PKG = go-starter/vars

BUILD_FLAGS  = -X '${VARS_PKG}.AppName=${BINARY_NAME}'
BUILD_FLAGS += -X '${VARS_PKG}.AppVersion=0.1.0'
BUILD_FLAGS += -X '${VARS_PKG}.GoVersion=$(shell go version)'
BUILD_FLAGS += -X '${VARS_PKG}.BuildTime=$(shell date +"%Y-%m-%d %H:%M:%S")'
BUILD_FLAGS += -X '${VARS_PKG}.GitCommit=$(shell git rev-parse HEAD)'
BUILD_FLAGS += -X '${VARS_PKG}.GitRemote=$(shell git config --get remote.origin.url)'

build:
	GOARCH=amd64 GOOS=linux go build -ldflags="${BUILD_FLAGS}" -o bin/${BINARY_NAME} app/main.go

test:
	go build -ldflags="${BUILD_FLAGS}" -o bin/${BINARY_NAME} app/main.go

test_version:
	bin/${BINARY_NAME} version

test_http:
	bin/${BINARY_NAME} http -c ../config/123.yml

clean:
	@go clean
	@rm -f bin/${BINARY_NAME}
	@rm -rf app/cmd/logs/
