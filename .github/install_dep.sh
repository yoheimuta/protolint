#!/bin/sh

set -eux

go install golang.org/x/tools/cmd/goimports@latest
go install github.com/kisielk/errcheck@latest
go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
go install github.com/gordonklaus/ineffassign@latest
go install github.com/opennota/check/cmd/varcheck@latest
go install github.com/opennota/check/cmd/aligncheck@latest
go install github.com/mdempsky/unconvert@latest
go install github.com/chavacava/garif@latest
