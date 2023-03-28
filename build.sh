#!/bin/bash

set -e

VERSION=v1.1.0

rm -rf bin

export GOOS=linux
export GOARCH=amd64
go build -trimpath -ldflags "-s -w" -o bin/azure-openai-proxy ./cmd

docker build -t stulzq/azure-openai-proxy:$VERSION .