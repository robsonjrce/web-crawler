
all:
	@echo "Tasks available:"
	@echo "  => build"
	@echo "  => test"

.PHONY: build
build: deps
	GOPATH=$(shell pwd)
	cd src/robsonjr.com.br/ && go build -o ../../bin/web-crawler main.go

.PHONY: deps
deps:
	GOPATH=$(shell pwd)
	cd src/robsonjr.com.br/ && go get -d ./...

.PHONY: test
test:
	GOPATH=$(shell pwd)
	cd src/robsonjr.com.br/ && GOPATH=$(GOPATH) go test -v robsonjr.com.br/utils/anchors robsonjr.com.br/utils/validation