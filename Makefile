NAME     := u2s3

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

.PHONY: cross-build
cross-build:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$(NAME)_darwin_amd64
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$(NAME)_darwin_arm64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$(NAME)_linux_amd64
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$(NAME)_linux_arm64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$(NAME)_windows_amd64.exe
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$(NAME)_windows_arm64.exe

.PHONY: test
test:
	env CUSTOM_HOST=localhost:9000 CUSTOM_REGION=us-east-1 go test -cover -v ./...

