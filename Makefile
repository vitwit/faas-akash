GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=pulsar

default: build

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

build: fmtcheck
	go build -o faas-akash

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

fmt:
	@echo "==> Fixing source code with gofmt..."
	@gofmt -s -w ./$(PKG_NAME)

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run -c golangci.yaml ./...

tools:
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: build test fmt fmtcheck lint tools