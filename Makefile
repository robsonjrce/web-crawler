
all:
	@echo "Tasks available:"
	@echo "  => build"
	@echo "  => test"

.PHONY: build
build: deps
	GOPATH=$(shell pwd)
	cd src/br.com.robsonjr/ && go build -o ../../bin/web-crawler main.go

.PHONY: deps
deps:
	GOPATH=$(shell pwd)
	cd src/br.com.robsonjr/ && go get -d ./...

.PHONY: test
test: