SHELL := /bin/bash
PLATFORM := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

default: build test

build:
	go fmt ./...
	GOOS=$(PLATFORM) GOARCH=$(GOARCH) vgo build -o awsswitch

install:
	#TODO

test:
	vgo test -race ./...

.PHONY: build install test run
