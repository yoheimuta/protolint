## test/all runs all related tests.
test/all: test/lint test

## test runs `go test`
test:
	go test -v -p 2 -count 1 -timeout 240s -race ./...

## test runs `go test -run $(RUN)`
test/run:
	go test -v -p 2 -count 1 -timeout 240s -race ./... -run $(RUN)

## test/lint/all runs all linter
test/lint/all: test/lint test/lint/nosupport/gomod

## test/lint runs linter
test/lint:
	# checks the coding style.
	(! gofmt -s -d `find . -name vendor -prune -type f -o -name '*.go'` | grep '^')
	golint -set_exit_status `go list ./...`
	# checks the import format.
	#(! goimports -l `find . -name vendor -prune -type f -o -name '*.go'` | grep 'go')
	# checks the error the compiler can't find.
	go vet ./...
	# checks shadowed variables.
	go vet -shadow ./...
	# checks no used assigned value.
	ineffassign .

## test/lint/nosupport/gomod runs linter without support of go mod
test/lint/nosupport/gomod:
	# checks not to ignore the error.
	errcheck ./...
	# checks unused global variables and constants.
	varcheck ./...
	# checks dispensable type conversions.
	unconvert -v ./...

## dev/install/dep installs depenencies required for development.
dev/install/dep:
	./.circleci/install_dep.sh

## ARG is command arguments.
ARG=lint

## run/cmd/pl runs pl with ARG
run/cmd/pl:
	go run cmd/pl/main.go $(ARG)

## build/cmd/pl builds pl
build/cmd/pl:
	go build -o pl cmd/pl/main.go