FROM alpine:3

COPY ./bin/azure-openai-proxy /usr/bin

ENTRYPOINT ["/usr/bin/azure-openai-proxy"]