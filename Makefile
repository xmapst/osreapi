SHELL=/bin/bash
GIT_URL := "https://github.com/xmapst/osreapi.git"
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD || echo "Unknown")
GIT_COMMIT := $(shell git rev-parse HEAD || echo "Unknown")
VERSION := $(shell git describe --tags --abbrev=0 || echo "Unknown")
USER_NAME := $(shell git config user.name || echo "Unknown")
USER_EMAIL := $(shell git config user.email || echo "Unknown")
BUILD_TIME := $(shell date +"%Y-%m-%d %H:%M:%S %Z" || echo "Unknown")
CGO_ENABLED := 0
LDFLAGS := "-w -s \
-X 'github.com/xmapst/osreapi/internal/info.Version=$(VERSION)' \
-X 'github.com/xmapst/osreapi/internal/info.GitUrl=$(GIT_URL)' \
-X 'github.com/xmapst/osreapi/internal/info.GitBranch=$(GIT_BRANCH)' \
-X 'github.com/xmapst/osreapi/internal/info.GitCommit=$(GIT_COMMIT)' \
-X 'github.com/xmapst/osreapi/internal/info.BuildTime=$(BUILD_TIME)' \
-X 'github.com/xmapst/osreapi/internal/info.UserName=$(USER_NAME)' \
-X 'github.com/xmapst/osreapi/internal/info.UserEmail=$(USER_EMAIL)' \
"

all: vet fmt windows linux darwin sha256sum

sha256sum:
	sha256sum bin/*remote_executor* > bin/latest.sha256sum

fmt:
	go fmt ./...

vet:
	go vet ./...

swag:
	swag init -d internal/router -g router.go -o internal/docs

dev:
	go mod tidy
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o bin/linux-remote_executor-amd64 cmd/osreapi.go
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o bin/windows-remote_executor-amd64.exe cmd/osreapi.go

windows:
	go mod tidy
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=386 go build -trimpath -ldflags $(LDFLAGS) -o bin/windows-remote_executor-386.exe cmd/osreapi.go
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o bin/windows-remote_executor-amd64.exe cmd/osreapi.go

linux:
	go mod tidy
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=386 go build -trimpath -ldflags $(LDFLAGS) -o bin/linux-remote_executor-386 cmd/osreapi.go
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o bin/linux-remote_executor-amd64 cmd/osreapi.go
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=arm go build -trimpath -ldflags $(LDFLAGS) -o bin/linux-remote_executor-arm cmd/osreapi.go
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=arm64 go build -trimpath -ldflags $(LDFLAGS) -o bin/linux-remote_executor-arm64 cmd/osreapi.go

darwin:
	go mod tidy
	CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o bin/darwin-remote_executor-amd64 cmd/osreapi.go
	CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags $(LDFLAGS) -o bin/darwin-remote_executor-arm64 cmd/osreapi.go