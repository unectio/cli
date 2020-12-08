all: uctl

.PHONY: all

uctl: .FORCE
	go build -o uctl

uctl-full: .FORCE
	go build -tags batcher -o uctl

deps:
	go mod download

.PHONY: .FORCE
