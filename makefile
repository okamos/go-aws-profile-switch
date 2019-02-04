SHELL := /bin/bash
PLATFORM := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

setup:
	go get -u golang.org/x/vgo
	vgo install

build:
	go fmt ./...
	GOOS=$(PLATFORM) GOARCH=$(GOARCH) vgo build -o awsswitch

analysis: setup
	vgo vet ./...
	vgo get -u golang.org/x/lint/golint
	golint -set_exit_status $$(vgo list ./... | grep -v /vendor/)

test: setup
	vgo test -race ./...

.PHONY: setup build analysis test
