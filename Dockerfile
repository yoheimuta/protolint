FROM alpine:3.14.3
ENTRYPOINT ["/usr/local/bin/protolint"]
COPY protolint /usr/local/bin/protolint