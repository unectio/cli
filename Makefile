all: lets

.PHONY: all

lets: .FORCE
	GOPATH=$(GOPATH):$(shell pwd) go build -o lets

lets-full: .FORCE
	GOPATH=$(GOPATH):$(shell pwd) go build -tags batcher -o lets

deps:
	./scripts/prep-devel.sh

.PHONY: .FORCE
