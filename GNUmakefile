
default: build

build:
	go build -v ./...

install: build
	go install -v ./...

# See https://golangci-lint.run/
lint:
	golangci-lint run

fmt: ## Run gofumpt
	@echo "==> Fixing source code with gofumpt..."
	gofumpt -l -w .

fumpt: fmt

test:
	go test -v -cover -timeout=120s -parallel=4 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate

# Generate docs, terraform fmt the examples folder, and create copywrite headers
generate:
	cd tools && go generate ./...

tools:
	cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	cd tools && go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	cd tools && go install github.com/hashicorp/copywrite
	cd tools && go install go install mvdan.cc/gofumpt@latest

generate-all-dto: ## Generate all DTOs for SDK interfaces
	go generate ./internal/sdk/*_dto.go

generate-dto-%: ./internal/sdk/%_dto.go ## Generate DTO for given SDK interface
	go generate $<

run-generator-poc:
	go generate ./internal/sdk/poc/example/*_def.go
	go generate ./internal/sdk/poc/example/*_dto_gen.go

clean-generator-poc:
	rm -f ./internal/sdk/poc/example/*_gen.go
	rm -f ./internal/sdk/poc/example/*_gen_test.go

run-generator-%: ./internal/sdk/%_def.go ## Run go generate on given object definition
	go generate $<
	go generate ./internal/sdk/$*_dto_gen.go

clean-generator-%: ## Clean generated files for specified resource
	rm -f ./internal/sdk/$**_gen.go
	rm -f ./internal/sdk/$**_gen_*test.go

.PHONY: build install lint generate fmt test testacc tools docs
