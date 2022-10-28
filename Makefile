SHA=$(shell git rev-parse --short HEAD)
export DIRTY=$(shell if `git diff-index --quiet HEAD --`; then echo false; else echo true;  fi)
export BASE_BINARY_NAME=terraform-provider-snowflake
export GO111MODULE=on
export TF_ACC_TERRAFORM_VERSION=0.13.0
export SKIP_EXTERNAL_TABLE_TESTS=true
export SKIP_SCIM_INTEGRATION_TESTS=true

go_test ?= -
ifeq (, $(shell which gotest))
	go_test=go test
else
	go_test=gotest
endif

all: test docs install
.PHONY: all

setup: ## setup development dependencies
	curl -sfL https://raw.githubusercontent.com/chanzuckerberg/bff/main/download.sh | sh
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh| sh
.PHONY: setup

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test ./pkg/resources -v -timeout 10m -sweep=prod

lint:  ## run the fast go linters
	./bin/reviewdog -conf .reviewdog.yml  -diff "git diff main"
.PHONY: lint

lint-ci: ## run the fast go linters
	./bin/reviewdog -conf .reviewdog.yml -reporter=github-pr-review -tee -fail-on-error=true
.PHONY: lint-ci

lint-all:  ## run the fast go linters
	./bin/reviewdog -conf .reviewdog.yml  -filter-mode nofilter
.PHONY: lint-all

lint-missing-acceptance-tests:
	@for r in `ls pkg/resources/ | grep -v list_expansion | grep -v privileges | grep -v grant_helpers | grep -v test | xargs -I{} basename {} .go`; do \
		if [ ! -f pkg/resources/"$$r"_acceptance_test.go ]; then \
			echo $$r; \
		fi; \
	done
.PHONY: lint-missing-acceptance-tests

build: ## build the binary
	go build ${LDFLAGS} -o $(BASE_BINARY_NAME) .
.PHONY: build

coverage: ## run the go coverage tool, reading file coverage.out
	go tool cover -html=coverage.txt
.PHONY: coverage

test:  ## run the tests
	CGO_ENABLED=1 $(go_test) -race -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) ./...
.PHONY: test

test-acceptance: ## runs all tests, including the acceptance tests which create and destroys real resources
	SKIP_MANAGED_ACCOUNT_TEST=1 TF_ACC=1 $(go_test) -v -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) ./...
.PHONY: test-acceptance

deps:
	go mod tidy -compat=1.19
.PHONY: deps

install: ## install the terraform-provider-snowflake binary in $GOPATH/bin
	go install ${LDFLAGS} .
.PHONY: install

install-tf: build ## installs plugin where terraform can find it
	mkdir -p $(HOME)/.terraform.d/plugins
	cp ./$(BASE_BINARY_NAME) $(HOME)/.terraform.d/plugins/$(BASE_BINARY_NAME)
.PHONY: install-tf

uninstall-tf: build ## uninstalls plugin from where terraform can find it
	rm $(HOME)/.terraform.d/plugins/$(BASE_BINARY_NAME) 2>/dev/null
.PHONY: install-tf

help: ## display help for this makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

clean: ## clean the repo
	rm terraform-provider-snowflake 2>/dev/null || true
	go clean
	rm -rf dist
.PHONY: clean

docs:
	SNOWFLAKE_USER= SNOWFLAKE_ACCOUNT= go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
.PHONY: docs

check-docs: docs ## check that docs have been generated
	git diff --exit-code -- docs
.PHONY: check-docs

check-mod:
	go mod tidy -compat=1.19
	git diff --exit-code -- go.mod go.sum
.PHONY: check-mod

