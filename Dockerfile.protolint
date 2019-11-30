FROM golang:1.13.4-alpine3.10 AS builder
RUN apk update
RUN apk add --virtual build-dependencies build-base git
ADD . ./protolint
WORKDIR ./protolint
RUN go mod download
RUN make build/cmd/protolint && mv protolint /go/bin/protolint

FROM alpine:3.10
RUN apk add --update --no-cache \
    ca-certificates && \
  rm -rf /var/cache/apk/*
COPY --from=builder /go/bin/protolint /usr/local/bin/protolint
ENTRYPOINT ["/usr/local/bin/protolint"]
