#!/usr/bin/env bash

set -euxo pipefail

go get -u golang.org/x/tools/cmd/goimports
go get -u golang.org/x/lint/golint
go get -u github.com/kisielk/errcheck
go get -u github.com/haya14busa/gosum/cmd/gosumcheck
go get -u github.com/gordonklaus/ineffassign
go get -u github.com/opennota/check/cmd/varcheck
go get -u github.com/opennota/check/cmd/aligncheck
go get -u github.com/mdempsky/unconvert