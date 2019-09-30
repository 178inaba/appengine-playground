.PHONY: all mod-download fmt fmt-diff ci-lint lint vet test deploy install-tools

all: fmt-diff ci-lint lint vet test

mod-download:
	go mod download

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
	gcloud app deploy --version $$(git rev-parse --short HEAD) --no-promote

install-tools:
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports golang.org/x/lint/golint github.com/golangci/golangci-lint/cmd/golangci-lint
