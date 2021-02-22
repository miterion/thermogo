GOCMD=go
GOLINT=golint
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build

GOSRC:=$(shell find . -name '*.go')
GOSRC+=go.mod go.sum


all: build

run:
	$(GORUN) main.go

lint: $(GOSRC)
	$(GOLINT) 

build: thermogo

thermogo: $(GOSRC)
	$(GOBUILD) -o $@ main.go

build-static:
	CGO_ENABLED=0 $(GOBUILD) -o thermogo-static -tags netgo main.go