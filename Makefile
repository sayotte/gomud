SHELL=/bin/bash
ROOT = $(shell pwd)
export GOPATH := ${ROOT}
export PATH := ${PATH}:${ROOT}/bin
LIBS = pkg/darwin_amd64/gomud/model.a pkg/darwin_amd64/gomud/controller.a pkg/darwin_amd64/gomud/supervisor.a

bin/gomud: ${LIBS} src/main/gomud/*.go
	go fmt main/gomud
	go vet main/gomud
	go install main/gomud

pkg/darwin_amd64/gomud/model.a: src/gomud/model/*.go
	go fmt gomud/model
	go vet gomud/model
	go install gomud/model

pkg/darwin_amd64/gomud/controller.a: src/gomud/controller/*.go
	go fmt gomud/controller
	go vet gomud/controller
	go install gomud/controller

pkg/darwin_amd64/gomud/supervisor.a: src/gomud/supervisor/*.go
	go fmt gomud/supervisor
	go vet gomud/supervisor
	go install gomud/supervisor

clean:
	rm -rf pkg/darwin_amd64/gomud* bin/gomud
