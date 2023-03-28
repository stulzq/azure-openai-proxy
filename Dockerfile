FROM alpine:3

EXPOSE 8080
COPY ./bin/azure-openai-proxy /usr/bin

ENTRYPOINT ["/usr/bin/azure-openai-proxy"]