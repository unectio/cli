VGOPATH=$(shell echo $$GOPATH | sed -e 's/:.*$$//')

all: .api
	GOPATH=$(shell pwd):$(GOPATH) go build -o function

.api:
	# A link to where langlet API files are (should be)
	ln -s $(VGOPATH)/src/github.com/unectio/api/langlet/golang $@
