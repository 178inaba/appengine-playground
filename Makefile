.PHONY: all fmt fmt-diff ci-lint lint vet test install-tools

all: fmt-diff ci-lint lint vet test

fmt:
	goimports -w .

fmt-diff:
	test -z $$(goimports -l .) || (goimports -d . && exit 1)

ci-lint:
	golangci-lint run

lint:
	golint -set_exit_status ./...

vet:
	go vet ./...

test:
	go test -race -count 1 ./...

deploy:
	echo TODO

install-tools:
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports golang.org/x/lint/golint github.com/golangci/golangci-lint/cmd/golangci-lint
