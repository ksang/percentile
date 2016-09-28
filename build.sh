#!/bin/bash
. ./env.sh

go build -o build/percentile percentile/*.go

export GOOS=linux
export GOARCH=amd64
go build -o build/linux/amd64/percentile percentile/*.go

unset GOOS
unset GOARCH