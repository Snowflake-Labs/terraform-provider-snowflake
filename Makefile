SHA=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)
export DIRTY=$(shell if `git diff-index --quiet HEAD --`; then echo false; else echo true;  fi)
# TODO add release flag
LDFLAGS=-ldflags "-w -s -X github.com/chanzuckerberg/terraform-provider-snowflake/pkg/version.GitSha=${SHA} -X github.com/chanzuckerberg/terraform-provider-snowflake/pkg/version.Version=${VERSION} -X github.com/chanzuckerberg/terraform-provider-snowflake/pkg/version.Dirty=${DIRTY}"
export BASE_BINARY_NAME=terraform-provider-snowflake_v$(VERSION)
export GO111MODULE=on

all: test docs install
.PHONY: all

setup: ## setup development dependencies
	curl -sfL https://raw.githubusercontent.com/chanzuckerberg/bff/master/download.sh | sh
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh| sh
.PHONY: setup

lint: ## run the fast go linters
	./bin/reviewdog -conf .reviewdog.yml  -diff "git diff master"
.PHONY: lint

lint-ci: ## run the fast go linters
	./bin/reviewdog -conf .reviewdog.yml  -reporter=github-pr-review
.PHONY: lint-ci

lint-all: ## run the fast go linters
	# doesn't seem to be a way to get reviewdog to not filter by diff
	./bin/golangci-lint run
.PHONY: lint-all

release: ## run a release
	./bin/bff bump
	git push
	goreleaser release
.PHONY: release

release-prerelease: build ## release to github as a 'pre-release'
	version=`./$(BASE_BINARY_NAME) -version`; \
	git tag v"$$version"; \
	git push
	git push --tags
	goreleaser release -f .goreleaser.prerelease.yml --debug
.PHONY: release-prerelease

release-snapshot: ## run a release
	goreleaser release --snapshot
.PHONY: release-snapshot

build: ## build the binary
	go build ${LDFLAGS} -o $(BASE_BINARY_NAME) .
.PHONY: build

coverage: ## run the go coverage tool, reading file coverage.out
	go tool cover -html=coverage.txt
.PHONY: coverage

test: deps ## run the tests
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
.PHONY: test

test-acceptance: deps ## runs all tests, including the acceptance tests which create and destroys real resources
	SKIP_WAREHOUSE_GRANT_TESTS=1 SKIP_SHARE_TESTS=1 SKIP_MANAGED_ACCOUNT_TEST=1 TF_ACC=1 go test -v -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) ./...
.PHONY: test-acceptance

test-acceptance-ci: ## runs all tests, including the acceptance tests which create and destroys real resources
	SKIP_WAREHOUSE_GRANT_TESTS=1 SKIP_SHARE_TESTS=1 SKIP_MANAGED_ACCOUNT_TEST=1 TF_ACC=1 go test -v -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) ./...
.PHONY: test-acceptance

deps:
	go mod tidy
.PHONY: deps

install: ## install the terraform-provider-snowflake binary in $GOPATH/bin
	go install ${LDFLAGS} .
.PHONY: install

install-tf: build ## installs plugin where terraform can find it
	mkdir -p $(HOME)/.terraform.d/plugins
	cp ./$(BASE_BINARY_NAME) $(HOME)/.terraform.d/plugins/$(BASE_BINARY_NAME)
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

check-docs: build ## check that docs have been generated
	./scripts/update-readme.sh check
.PHONY: check-docs

check-mod:
	go mod tidy
	git diff --exit-code -- go.mod go.sum
.PHONY: check-mod