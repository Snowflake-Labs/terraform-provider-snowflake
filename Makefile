export TEST_SF_TF_SKIP_SAML_INTEGRATION_TEST=true
export TEST_SF_TF_SKIP_MANAGED_ACCOUNT_TEST=true
export BASE_BINARY_NAME=terraform-provider-snowflake
export TERRAFORM_PLUGINS_DIR=$(HOME)/.terraform.d/plugins
export TERRAFORM_PLUGIN_LOCAL_INSTALL=$(TERRAFORM_PLUGINS_DIR)/$(BASE_BINARY_NAME)

default: help

dev-setup: ## setup development dependencies
	@which ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.55.2
	cd tools && mkdir -p bin/
	cd tools && env GOBIN=$$PWD/bin go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	cd tools && env GOBIN=$$PWD/bin go install mvdan.cc/gofumpt

dev-cleanup: ## cleanup development dependencies
	rm -rf bin/*
	rm -rf tools/bin/*

docs: generate-doc-templates ## generate docs
	tools/bin/tfplugindocs generate

docs-check: docs ## check that docs have been generated
	git diff --exit-code -- docs

fmt: terraform-fmt ## Run terraform fmt and gofumpt
	tools/bin/gofumpt -l -w .

terraform-fmt: ## Run terraform fmt
	terraform fmt -recursive ./examples/
	terraform fmt -recursive ./pkg/resources/testdata/
	terraform fmt -recursive ./pkg/datasources/testdata/

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-23s\033[0m %s\n", $$1, $$2}'

install: ## install the binary
	go install -v ./...

lint: # Run static code analysis, check formatting. See https://golangci-lint.run/
	./bin/golangci-lint run ./... -v

lint-fix: ## Run static code analysis, check formatting and try to fix findings
	./bin/golangci-lint run ./... -v --fix

mod: ## add missing and remove unused modules
	go mod tidy -compat=1.21

mod-check: mod ## check if there are any missing/unused modules
	git diff --exit-code -- go.mod go.sum

pre-push: mod fmt generate-doc-templates docs lint test-architecture ## Run a few checks before pushing a change (docs, fmt, mod, etc.)

pre-push-check: pre-push mod-check generate-doc-templates-check docs-check ## Run checks before pushing a change (docs, fmt, mod, etc.)

sweep: ## destroy the whole architecture; USE ONLY FOR DEVELOPMENT ACCOUNTS
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	@echo "Are you sure? [y/n]" >&2
	@read -r REPLY; \
		if echo "$$REPLY" | grep -qG "^[yY]$$"; then \
			TEST_SF_TF_ENABLE_SWEEP=1 go test -timeout 300s -run ^TestSweepAll ./pkg/sdk -v; \
			else echo "Aborting..."; \
		fi;

test: test-client ## run unit and integration tests
	go test -v -cover -timeout=30m ./...

test-acceptance: ## run acceptance tests
	TF_ACC=1 SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE=true go test -run "^TestAcc_" -v -cover -timeout=60m ./...

test-integration: ## run SDK integration tests
	go test -run "^TestInt_" -v -cover -timeout=30m ./...

test-architecture: ## check architecture constraints between packages
	go test ./pkg/architests/... -v

test-client: ## runs test that checks sdk.Client without instrumentedsql
	SF_TF_NO_INSTRUMENTED_SQL=1 SF_TF_GOSNOWFLAKE_LOG_LEVEL=debug go test ./pkg/sdk/internal/client/... -v

build-local: ## build the binary locally
	go build -o $(BASE_BINARY_NAME) .

install-tf: build-local ## installs plugin where terraform can find it
	mkdir -p $(TERRAFORM_PLUGINS_DIR)
	cp ./$(BASE_BINARY_NAME) $(TERRAFORM_PLUGIN_LOCAL_INSTALL)

uninstall-tf: ## uninstalls plugin from where terraform can find it
	rm -f $(TERRAFORM_PLUGIN_LOCAL_INSTALL)

generate-all-dto: ## Generate all DTOs for SDK interfaces
	go generate ./pkg/sdk/*_dto.go

generate-dto-%: ./pkg/sdk/%_dto.go ## Generate DTO for given SDK interface
	go generate $<

clean-generator-poc:
	rm -f ./pkg/sdk/poc/example/*_gen.go
	rm -f ./pkg/sdk/poc/example/*_gen_test.go

clean-generator-%: ## Clean generated files for specified resource
	rm -f ./pkg/sdk/$**_gen.go
	rm -f ./pkg/sdk/$**_gen_*test.go

run-generator-poc:
	go generate ./pkg/sdk/poc/example/*_def.go
	go generate ./pkg/sdk/poc/example/*_dto_gen.go

run-generator-%: ./pkg/sdk/%_def.go ## Run go generate on given object definition
	go generate $<
	go generate ./pkg/sdk/$*_dto_gen.go

generate-doc-templates: ## generate doc templates
	go run ./pkg/internal/tools/doc-gen-helper/ $$PWD

generate-doc-templates-check: generate-doc-templates ## check that doc templates have been generated
	git diff --exit-code -- templates

.PHONY: build-local clean-generator-poc dev-setup dev-cleanup docs docs-check fmt fmt-check fumpt help install lint lint-fix mod mod-check pre-push pre-push-check sweep test test-acceptance uninstall-tf
