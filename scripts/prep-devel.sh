#!/bin/bash
VGOPATH=$(echo $GOPATH | sed -e 's/:.*$//')

if [ "x${VGOPATH}" == "x" ]; then
	echo "Set GOPATH to somewhere" >&2
	exit 0
fi

set -x
set -e

git --version
git config --global core.protectNTFS false

go get github.com/spf13/cobra
go get github.com/dgrijalva/jwt-go
go get gopkg.in/yaml.v2
go get go.uber.org/zap
go get github.com/gorilla/mux
go get gopkg.in/mgo.v2

go get github.com/unectio/util
go get github.com/unectio/db
go get github.com/unectio/api
