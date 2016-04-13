#!/bin/sh

#app=`pwd | awk -F "/" '{print $(NF)}'`
app="domeize"

docker run --rm -e GOPATH=/go/$app -v "$PWD":/go/$app -w /go/$app golang:1.5-alpine go build -v
