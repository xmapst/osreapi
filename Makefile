SHELL=/bin/bash
BINARY_NAME=osreapi
VERSION=`git describe --tags --abbrev=0`
GO_VERSION=`go version|awk '{print $$3" "$$4}'`
GIT_URL=`git remote -v|grep push|awk '{print $$2}'`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
GIT_COMMIT=`git rev-parse HEAD`
GIT_LATEST_TAG=`git describe --tags --abbrev=0`
BUILD_TIME=`date +"%Y-%m-%d %H:%M:%S %Z"`

LDFLAGS="-X 'github.com/xmapst/osreapi.Version=${VERSION}' -X 'github.com/xmapst/osreapi.GoVersion=${GO_VERSION}' -X 'github.com/xmapst/osreapi.GitUrl=${GIT_URL}' -X 'github.com/xmapst/osreapi.GitBranch=${GIT_BRANCH}' -X 'github.com/xmapst/osreapi.GitCommit=${GIT_COMMIT}' -X 'github.com/xmapst/osreapi.GitLatestTag=${GIT_LATEST_TAG}' -X 'github.com/xmapst/osreapi.BuildTime=${BUILD_TIME}'"

all: linux windows

linux:
	git pull
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${BINARY_NAME}_linux cmd/linux/main.go
	strip --strip-unneeded ${BINARY_NAME}_linux
	upx --lzma ${BINARY_NAME}_linux

windows:
	git pull
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${BINARY_NAME}_windows.exe cmd/windows/main.go
	strip --strip-unneeded ${BINARY_NAME}_windows.exe
	upx --lzma ${BINARY_NAME}_windows.exe