FROM golang:1.22-bookworm AS builder

COPY . /builder
WORKDIR /builder

RUN make build

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /app

EXPOSE 8080
COPY --from=builder /builder/bin .

ENTRYPOINT ["/app/azure-openai-proxy"]
