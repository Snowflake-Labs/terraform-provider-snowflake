## SDK DTO builder generation

Short PoC of generating DTO constructors and builder methods.

### Description

There is an example file ready for generation [pipes_dto.go](example/pipes_dto.go) which creates file [pipes_dto_generated.go](example/pipes_dto_generated.go).

Required fields should be marked with `// required` comment.

To mark file inside `pkg/sdk/` directory as ready for generation add to the file: `//go:generate go run ./dto-builder-generator/main.go`.

Output file will contain the same set of imports as the input file and will be formatted.

### Usage

To invoke example generation run:
```shell
go generate pkg/sdk/dto-builder-generator/example/pipes_dto.go
```

To invoke all generations run:
```shell
make generate-all-dto
```

To invoke only generation of given resource (e.g. pipes), run:
```shell
make generate-dto-pipes
```

### Next steps
- if comments are not enough, use different method to mark required fields (e.g. struct tags)
- generate mappings between dto and Options struct
- add more meta info to generated file header comment (e.g. time of generation etc.)
