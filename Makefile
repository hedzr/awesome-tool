include .env
# ref: https://kodfabrik.com/journal/a-good-makefile-for-go/

PROJECTNAME=$(shell basename "$(PWD)")
APPNAME=$(patsubst "%",%,$(shell grep -E "AppName[ \t]+=[ \t]+" doc.go|grep -Eo "\\\".+\\\""))
VERSION=$(shell grep -E "Version[ \t]+=[ \t]+" doc.go|grep -Eo "[0-9.]+")

# Go related variables.
GOBASE=$(shell pwd)
#GOPATH="$(GOBASE)/vendor:$(GOBASE)"
#GOPATH=$(GOBASE)/vendor:$(GOBASE):$(shell dirname $(GOBASE))
GOPATH2=$(shell dirname $(GOBASE))
GOPATH1=$(shell dirname $(GOPATH2))
GOPATH=$(shell dirname $(GOPATH1))
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

GO111MODULE=on
#GOPROXY=https://goproxy.io

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID=/tmp/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

#
LDFLAGS=
M = $(shell printf "\033[34;1m▶\033[0m")


## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: info go-get

info:
	@echo "     GOBASE: $(GOBASE)"
	@echo "      GOBIN: $(GOBIN)"
	@echo "     GOPATH: $(GOPATH)"
	@echo "GO111MODULE: $(GO111MODULE)"
	@echo "    GOPROXY: $(GOPROXY)"
	@echo "    APPNAME: $(APPNAME)"
	@echo "    VERSION: $(VERSION)"
	@echo

## start: Start in development mode. Auto-starts when code changes.
start:
	bash -c "trap 'make stop' EXIT; $(MAKE) compile start-server watch run='make compile start-server'"

## stop: Stop development mode.
stop: stop-server

start-server: stop-server
	@echo "  >  $(PROJECTNAME) is available at $(ADDR)"
	@-$(GOBIN)/$(PROJECTNAME) 2>&1 & echo $$! > $(PID)
	@cat $(PID) | sed "/^/s/^/  \>  PID: /"

stop-server:
	@-touch $(PID)
	@-kill `cat $(PID)` 2> /dev/null || true
	@-rm $(PID)

## watch: Run given command when code changes. e.g; make watch run="echo 'hey'"
watch:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) yolo -i . -e vendor -e bin -c "$(run)"

restart-server: stop-server start-server

## compile: Compile the binary.
compile:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## exec: Run given command, wrapped with custom GOPATH. e.g; make exec run="go test ./..."
exec:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(run)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@(MAKEFILE) go-clean

## build: Compile the binary.
build: go-compile

go-compile: go-clean go-get go-build

goarch=amd64
W_PKG=github.com/hedzr/cmdr/conf
TIMESTAMP=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
GITHASH=$(shell git rev-parse HEAD)
GOVERSION=$(shell go version)
LDFLAGS=-s -w -X '$(W_PKG).Buildstamp=$(TIMESTAMP)' -X '$(W_PKG).Githash=$(GITHASH)' -X '$(W_PKG).GoVersion=$(GOVERSION)' -X '$(W_PKG).Version=$(VERSION)' -X '$(W_PKG).AppName=$(APPNAME)'

build-ci:
	@echo "  >  Building binary..."
	@echo "  >  LDFLAGS = $(LDFLAGS)"
	$(foreach os, darwin linux windows, \
	  echo "     Building $(GOBIN)/$(PROJECTNAME)_$(os)_$(goarch)...$(os)"; \
	      GOARCH="$(goarch)" GOOS="$(os)" GOPATH="$(GOPATH)" GOBIN="$(GOBIN)" GO111MODULE="$(GO111MODULE)" GOPROXY="$(GOPROXY)" \
	        go build -ldflags "$(LDFLAGS)" -o $(GOBIN)/$(PROJECTNAME)_$(os)_$(goarch) $(GOBASE)/cli/main.go; \
	        gzip $(GOBIN)/$(PROJECTNAME)_$(os)_$(goarch); \
	)
	@ls -la $(GOBIN)/*

go-build:
	@echo "  >  Building binary..."
	@echo "  >  LDFLAGS = $(LDFLAGS)"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	  go build -ldflags "$(LDFLAGS)" -o $(GOBIN)/$(PROJECTNAME) $(GOBASE)/cli/main.go
	# go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)
	# chmod +x $(GOBIN)/*
	ls -la $(GOBIN)/$(PROJECTNAME)

go-generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go generate $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies...$(get)"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go get $(get)

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go clean

## go-format: run gofmt tool
go-format:
	@echo "  >  gofmt ..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	gofmt -l -w -s .

## go-lint: run golint tool
go-lint:
	@echo "  >  golint ..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	golint ./...

## go-coverage: run go coverage test
go-coverage:
	@echo "  >  gocov ..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go test -race -covermode=atomic -coverprofile cover.out
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go tool cover -html=cover.out -o cover.html
	@open cover.html

## go-codecov: run go test for codecov; (codecov.io)
go-codecov:
	# https://codecov.io/gh/hedzr/cmdr
	@echo "  >  codecov ..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go test -race -coverprofile=coverage.txt -covermode=atomic
	@bash <(curl -s https://codecov.io/bash) -t $(CODECOV_TOKEN)
	# open https://codecov.io/gh/hedzr/cmdr

## go-cyclo: run gocyclo tool
go-cyclo:
	@echo "  >  gocyclo ..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	gocyclo -top 20 .


## run: build, and test, ...
run: go-build
	@echo "  >  run ..."
	$(GOBIN)/$(PROJECTNAME) build one \
		--name=awesome-go \
		--source=https://github.com/avelino/awesome-go \
		--work-dir=./output


.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo