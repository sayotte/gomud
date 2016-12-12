SHELL=/bin/bash
ROOT = $(shell pwd)
export GOPATH := ${ROOT}
export PATH := ${PATH}:${ROOT}/bin

bin/gomud: pkg/darwin_amd64/gomud.a src/main/gomud/main.go
	go fmt src/main/gomud/*
	go vet src/main/gomud/*
	go install main/gomud

pkg/darwin_amd64/gomud.a: src/gomud/*.go
	go fmt src/gomud/*.go
	go vet src/gomud/*.go
	go install gomud

clean:
	rm -rf pkg/darwin_amd64/gomud.a bin/gomud
