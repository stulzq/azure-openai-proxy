FROM golang:1.19 AS builder

COPY . /builder
WORKDIR /builder

RUN make build

FROM scratch

WORKDIR /app

EXPOSE 8080
COPY --from=builder /builder/bin .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/azure-openai-proxy"]
