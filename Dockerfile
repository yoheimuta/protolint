FROM alpine:3.12
ENTRYPOINT ["/usr/local/bin/protolint"]
COPY protolint /usr/local/bin/protolint