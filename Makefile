LDFLAGS := -s -w
BIN_NAME := "azure-openai-proxy"

build:
	@env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BIN_NAME) ./cmd

fmt:
	go fmt ./...

vet:
	go vet ./...

image:
	docker build -t azure-openai-proxy:latest .

.PHONY: build fmt vet