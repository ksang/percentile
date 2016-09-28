#!/bin/bash
. ./env.sh
export GOOS=darwin
export GOARCH=amd64
go build -o build/darwin/amd64/percentile percentile/*.go

export GOOS=linux
export GOARCH=amd64
go build -o build/linux/amd64/percentile percentile/*.go
