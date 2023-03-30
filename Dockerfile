FROM golang:1.19 AS building

COPY . /building
WORKDIR /building

RUN make build

FROM alpine:3

WORKDIR /app

EXPOSE 8080
COPY --from=building /building/bin .

ENTRYPOINT ["/app/azure-openai-proxy"]