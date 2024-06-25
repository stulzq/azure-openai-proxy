#!/bin/bash

set -e

VERSION=v1.3.10

rm -rf bin

export GOOS=linux
export GOARCH=amd64
go build -trimpath -ldflags "-s -w" -o bin/azure-openai-proxy ./cmd

docker buildx build --platform linux/amd64 -t stulzq/azure-openai-proxy:$VERSION .
