#!/bin/sh

set -eux

go get -u golang.org/x/tools/cmd/goimports
go get -u github.com/kisielk/errcheck
go get -u golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
go get -u github.com/gordonklaus/ineffassign
go get -u github.com/opennota/check/cmd/varcheck
go get -u github.com/opennota/check/cmd/aligncheck
go get -u github.com/mdempsky/unconvert
go get -u github.com/chavacava/garif
