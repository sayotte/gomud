SHELL=/bin/bash
ROOT = $(shell pwd)
export GOPATH := ${ROOT}
export PATH := ${PATH}:${ROOT}/bin
SOURCES := $(shell find src/ -name '*.go')

bin/gomud: $(SOURCES)
	go fmt ./...
	go vet ./...
	go install ./...

clean:
	rm -rf pkg bin/gomud
