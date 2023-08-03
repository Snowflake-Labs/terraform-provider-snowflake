export BASE_BINARY_NAME=terraform-provider-snowflake
export GO111MODULE=on
export TF_ACC_TERRAFORM_VERSION=1.4.1
export SKIP_EXTERNAL_TABLE_TESTS=true
export SKIP_SCIM_INTEGRATION_TESTS=true

all: test docs
.PHONY: all

setup: ## setup development dependencies
	curl -sfL https://raw.githubusercontent.com/chanzuckerberg/bff/main/download.sh | sh
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh| sh
.PHONY: setup

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	SNOWFLAKE_ENABLE_SWEEP=1 go test -timeout 300s -run ^TestSweepAll ./pkg/sdk -v

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
	go build -o $(BASE_BINARY_NAME) .
.PHONY: build

coverage: ## run the go coverage tool, reading file coverage.out
	go tool cover -html=coverage.txt
.PHONY: coverage

test:  ## run the tests (except sdk tests)
	CGO_ENABLED=1 go test -race -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) ./pkg/resources/...
	CGO_ENABLED=1 go test -race -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) ./pkg/provider/...
	CGO_ENABLED=1 go test -race -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) ./pkg/snowflake/...

.PHONY: test

test-acceptance: ## runs all tests, including the acceptance tests which create and destroys real resources
	SKIP_MANAGED_ACCOUNT_TEST=1 SKIP_EMAIL_INTEGRATION_TESTS=1 TF_ACC=1 go test -timeout 1200s -v -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) ./...
.PHONY: test-acceptance

deps:
	go mod tidy -compat=1.20
.PHONY: deps

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
	go mod tidy -compat=1.20
	git diff --exit-code -- go.mod go.sum
.PHONY: check-mod

.PHONY: fmt
fmt: ## Run linter and apply formatting autofix
	golangci-lint run ./... -v --fix
