#!/bin/sh
set -eux

dep ensure
export CGO_ENABLED=0
export GOARCH=amd64
mkdir -p bin
rm -f bin/*

GOOS=darwin  go build -o bin/packer-post-processor-teamcity.macos
GOOS=linux   go build -o bin/packer-post-processor-teamcity.linux
GOOS=windows go build -o bin/packer-post-processor-teamcity.exe
