LDFLAGS := -s -w
BIN_NAME := "azure-openai-proxy"
TARGETOS ?= linux
TARGETARCH ?= amd64

build:
	@env CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BIN_NAME) ./cmd

.PHONY: build