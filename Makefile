SHA=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)
DIRTY=$(shell if `git diff-index --quiet HEAD --`; then echo false; else echo true;  fi)
# TODO add release flag
LDFLAGS=-ldflags "-w -s -X github.com/chanzuckerberg/terraform-provider-snowflake/util.GitSha=${SHA} -X github.com/chanzuckerberg/terraform-provider-snowflake/util.Version=${VERSION} -X github.com/chanzuckerberg/terraform-provider-snowflake/util.Dirty=${DIRTY}"
GOTEST=gotest

all: test docs install
.PHONY: all

setup: ## setup development dependencies
	go get github.com/rakyll/gotest
	go install github.com/rakyll/gotest
	curl -L https://raw.githubusercontent.com/chanzuckerberg/bff/master/download.sh | sh
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
.PHONY: setup

lint: ## run the fast go linters
	golangci-lint run
.PHONY: lint

release: ## run a release
	./bin/bff bump
	git push
	goreleaser release
.PHONY: release

release-prerelease: build ## release to github as a 'pre-release'
	version=`./terraform-provider-snowflake version`; \
	git tag v"$$version"; \
	git push
	git push --tags
	goreleaser release -f .goreleaser.prerelease.yml --debug
.PHONY: release-prerelease

release-snapshot: ## run a release
	goreleaser release --snapshot
.PHONY: release-snapshot

build: ## build the binary
	go build ${LDFLAGS} .
.PHONY: build

coverage: ## run the go coverage tool, reading file coverage.out
	go tool cover -html=coverage.txt
.PHONY: coverage

test: ## run the tests
	${GOTEST} -race -coverprofile=coverage.txt -covermode=atomic ./...
.PHONY: test

test-acceptance: ## runs all tests, including the acceptance tests which create and destroys real resources
	TF_ACC=1 ${GOTEST} -v -race -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) ./...
.PHONY: test-acceptance

install: ## install the terraform-provider-snowflake binary in $GOPATH/bin
	go install ${LDFLAGS} .
.PHONY: install

install-tf: build ## installs plugin where terraform can find it
	mkdir -p $(HOME)/.terraform.d/plugins
	cp ./terraform-provider-snowflake $(HOME)/.terraform.d/plugins
.PHONY: install-tf

help: ## display help for this makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

clean: ## clean the repo
	rm terraform-provider-snowflake 2>/dev/null || true
	go clean
	rm -rf dist
.PHONY: clean

docs: build ## generate some docs
	./scripts/update-readme.sh update
.PHONY: docs

deps: ## installs and vendors dependencies
	go get .
	go mod tidy
	go mod vendor
.PHONY: deps

check-docs: build ## check that docs have been generated
	./scripts/update-readme.sh check
.PHONY: check-docs
