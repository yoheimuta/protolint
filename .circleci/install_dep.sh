#!/usr/bin/env bash

set -euxo pipefail

go get -u golang.org/x/tools/cmd/goimports
go get -u golang.org/x/lint/golint
go get -u github.com/kisielk/errcheck
go get -u golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
GO111MODULE=off go get -u github.com/haya14busa/gosum/cmd/gosumcheck
GO111MODULE=off go get -u github.com/gordonklaus/ineffassign
GO111MODULE=off go get -u github.com/opennota/check/cmd/varcheck
GO111MODULE=off go get -u github.com/opennota/check/cmd/aligncheck
GO111MODULE=off go get -u github.com/mdempsky/unconvert
