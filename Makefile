export SKIP_EMAIL_INTEGRATION_TESTS=true
export SKIP_EXTERNAL_TABLE_TEST=true
export SKIP_NOTIFICATION_INTEGRATION_TESTS=true
export SKIP_SAML_INTEGRATION_TESTS=true
export SKIP_STREAM_TEST=true

default: help

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate

docs-check: docs ## check that docs have been generated
	git diff --exit-code -- docs

fmt: ## Run gofumpt
	@echo "==> Fixing source code with gofumpt..."
	gofumpt -l -w .

fumpt: fmt

# Generate docs, terraform fmt the examples folder
generate:
	cd tools && go generate ./...

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-23s\033[0m %s\n", $$1, $$2}'

install:
	go install -v ./...

# See https://golangci-lint.run/
lint:
	golangci-lint run ./... -v

lint-fix: ## Run static code analysis, check formatting and try to fix findings
	golangci-lint run ./... -v --fix

mod: ## add missing and remove unused modules
	go mod tidy -compat=1.20

mod-check: mod ## check if there are any missing/unused modules
	git diff --exit-code -- go.mod go.sum

pre-push: fmt lint mod docs ## Run a few checks before pushing a change (docs, fmt, mod, etc.)

pre-push-check: docs-check lint-check mod-check; ## Run a few checks before pushing a change (docs, fmt, mod, etc.)

sweep: ## destroy the whole architecture; USE ONLY FOR DEVELOPMENT ACCOUNTS
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	@echo "Are you sure? [y/n]" >&2
	@read -r REPLY; \
		if echo "$$REPLY" | grep -qG "^[yY]$$"; then \
			SNOWFLAKE_ENABLE_SWEEP=1 go test -timeout 300s -run ^TestSweepAll ./pkg/sdk -v; \
			else echo "Aborting..."; \
		fi;

test:
	go test -v -cover -timeout=30m -parallel=4 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 30m -parallel=4  `go list ./... | grep -v pkg/sdk`

tools:
	cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	cd tools && go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	cd tools && go install github.com/hashicorp/copywrite
	cd tools && go install mvdan.cc/gofumpt

generate-all-dto: ## Generate all DTOs for SDK interfaces
	go generate ./internal/sdk/*_dto.go

generate-dto-%: ./internal/sdk/%_dto.go ## Generate DTO for given SDK interface
	go generate $<

clean-generator-poc:
	rm -f ./internal/sdk/poc/example/*_gen.go
	rm -f ./internal/sdk/poc/example/*_gen_test.go

clean-generator-%: ## Clean generated files for specified resource
	rm -f ./internal/sdk/$**_gen.go
	rm -f ./internal/sdk/$**_gen_*test.go

run-generator-poc:
	go generate ./internal/sdk/poc/example/*_def.go
	go generate ./internal/sdk/poc/example/*_dto_gen.go

run-generator-%: ./internal/sdk/%_def.go ## Run go generate on given object definition
	go generate $<
	go generate ./internal/sdk/$*_dto_gen.go

.PHONY: build clean-generator-poc clean-generator-% docs docs-check fmt fumpt generate generate-all-dto generate-dto-% help install lint lint-fix mod mod-check pre-push pre-push-check run-generator-poc run-generator-% sweep test testacc tools
