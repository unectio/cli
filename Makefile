all: uctl

.PHONY: all

uctl: .FORCE
	GOPATH=$(GOPATH):$(shell pwd) go build -o uctl

uctl-full: .FORCE
	GOPATH=$(GOPATH):$(shell pwd) go build -tags batcher -o uctl

deps:
	GOPATH=$(GOPATH):$(shell pwd) ./scripts/prep-devel.sh

.PHONY: .FORCE
