#!/bin/bash
VGOPATH=$(echo $GOPATH | sed -e 's/:.*$//')

if [ "x${1}" != "x" ]; then
	echo "Set GOPATH to ${1}"
	VGOPATH=${1}
    GOPATH=${1}
fi

if [ "x${VGOPATH}" == "x" ]; then
	echo "Set GOPATH to somewhere" >&2
	exit 0
fi

set -x
set -e

go get github.com/spf13/cobra
echo $?
ls -l ${GOPATH}"/src/github.com/spf13/cobra"
go get github.com/dgrijalva/jwt-go
echo $?
ls -l ${GOPATH}"/src/github.com/dgrijalva/jwt-go"
go get gopkg.in/yaml.v2
echo $?
go get go.uber.org/zap
echo $?
go get github.com/gorilla/mux
echo $?
go get gopkg.in/mgo.v2
echo $?

go get github.com/unectio/db
echo $?
go get github.com/unectio/api
echo $?
go get github.com/unectio/util
echo $?
