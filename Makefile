lets: .FORCE
	go build -o lets

lets-full: .FORCE
	go build -tags batcher -o lets

.PHONY: .FORCE
