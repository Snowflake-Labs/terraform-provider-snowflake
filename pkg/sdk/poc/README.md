## SDK generator PoC

PoC of generating full object implementation based on object definition.

### Description

There is an example file ready for generation [database_role_def.go](example/database_role_def.go) which creates files:
- [database_role_gen.go](example/database_role_gen.go) - SDK interface, options structs
- [database_role_dto_gen.go](example/database_role_dto_gen.go) - SDK Request DTOs
- [database_role_dto_builders_gen.go](example/database_role_dto_builders_gen.go) - SDK Request DTOs constructors and builder methods (this file is generated using [dto-builder-generator](../dto-builder-generator/main.go))
- [database_role_validations_gen.go](example/database_role_validations_gen.go) - options structs validations
- [database_role_impl_gen.go](example/database_role_impl_gen.go) - SDK interface implementation
- [database_role_gen_test.go](example/database_role_gen_test.go) - unit tests placeholders with guidance comments (at least for now)
- [database_role_gen_integration_test.go](example/database_role_gen_integration_test.go) - integration test placeholder file

### How it works
##### Creating object generation definition

To create definition for object generation:

1. Create file `object_name_def.go` (like example [database_role_def.go](example/database_role_def.go) file).
2. Put go generate directive at the top: `//go:generate go run ../main.go`. Remember that you may have to change the path to [main.go](main.go) file.
3. Create object interface definition.
4. Add key-value entry to `definitionMapping` in [main.go](main.go):
   - key should be created file name (for [database_role_def.go](example/database_role_def.go) example file: `"database_role_def.go"`)
   - value should be created definition (like for [database_role_def.go](example/database_role_def.go) example file: `DatabaseRole`)
5. You are all set to run generation.

##### Invoking generation

To invoke example generation (with first cleaning all the generated files) run:
```shell
make clean-generator-poc run-generator-poc
```

### Next steps
##### Essentials
- use DSL to build object definitions (from branch [go-builder-dsl](https://github.com/Snowflake-Labs/terraform-provider-snowflake/tree/go-builder-dsl)) - ideally leave two options of defining objects and proceed with generation based on definition provided
- differentiate between different actions implementations (now only `Create` and `Alter` has been considered, `Show` on the other hand has totally different implementation)
- generate `struct`s for `Show` and `ShowID`
- handle arrays
- handle more validation types

##### Improvements
- automatic names of nested `struct`s (e.g. `DatabaseRoleRename`)
- check if generating with package name + invoking format removes unnecessary qualifier
- consider merging templates `StructTemplate` and `OptionsTemplate` (requires moving Doc to Field)
- add unit tests to this generator

##### Known issues
- spaces in templates (especially nested validations)
