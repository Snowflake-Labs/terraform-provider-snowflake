## SDK DTO builder generation

Short PoC of generating DTO constructors and builder methods.

### Usage

There is an example file ready for generation [pipes_dto.go](../pipes_dto.go) which creates file [pipes_dto_generated.go](../pipes_dto_generated.go).

Required fields should be marked with `// required` comment.

Output file will contain the same set of imports as the input file and will be formatted.

To invoke all generations run:
```shell
make generate-all-dto
```

To invoke only generation of given resource (usable later when we have more dtos to generate), run either:
```shell
make generate-dto-pipes
```
