FROM alpine:3.16.2
RUN apk -U --no-cache upgrade; /bin/rm -rf /var/cache/apk/*
ENTRYPOINT ["/usr/local/bin/protolint"]
COPY protolint /usr/local/bin/protolint
