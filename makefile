SHELL := /bin/bash
PLATFORM := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

setup:
	go mod download

build:
	go fmt ./...
	GOOS=$(PLATFORM) GOARCH=$(GOARCH) vgo build -o awsswitch

analysis: setup
	go vet ./...
	go get -u golang.org/x/lint/golint
	golint -set_exit_status $$(vgo list ./... | grep -v /vendor/)

test: setup
	go test -race ./...

.PHONY: setup build analysis test
