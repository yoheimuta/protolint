FROM alpine:3.17.0
RUN apk -U --no-cache upgrade; /bin/rm -rf /var/cache/apk/*
ENTRYPOINT ["/usr/local/bin/protolint"]
COPY protolint /usr/local/bin/protolint
