LOCAL_PACKAGE_PREFIX := github.com/178inaba/appengine-playground

.PHONY: all fmt fmt-diff ci-lint ci-lint-fix vet test install-tools go-install-tools

all: ci-lint-fix fmt ci-lint vet test

fmt:
	goimports -local '$(LOCAL_PACKAGE_PREFIX)' -w .

fmt-diff:
	test -z $$(goimports -local '$(LOCAL_PACKAGE_PREFIX)' -l .) || (goimports -local '$(LOCAL_PACKAGE_PREFIX)' -d . && exit 1)

ci-lint:
	golangci-lint run

ci-lint-fix:
	golangci-lint run --fix

vet:
	go vet ./...

test:
	go test -race -count 1 -cover ./...

install-tools: go-install-tools
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.51.2

go-install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
