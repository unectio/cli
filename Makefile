lets: .FORCE
	go build

lets-full: .FORCE
	go build -tags batcher

.PHONY: .FORCE
