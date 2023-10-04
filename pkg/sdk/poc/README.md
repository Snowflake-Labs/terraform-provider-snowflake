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

To invoke generation inside SDK package (with cleaning), e.g. for `session_policies` run (mind the `_`(underscore)):
```shell
make clean-generator-session_policies run-generator-session_policies
```

### Next steps
##### Essentials
- fix builder generation (`With`s for optional fields should have required param, optional fields should not be exported in `Request` structs)
- (?) add mapping from db to plain struct (for now only "// TODO: Mapping" comment is generated)
- add arguments to the generator, so we'll be able to specify which files should be generated / re-generated,
because after we fill things that need our input we don't want to re-generate those files and override the changes,
also adding small changes is very challenging, e.g. for new validation rule you have to re-generate unit-tests to get
one new function, revert to old tests (the one with filled tests), copy new test case (of course we could add that one by hand
but if we add one case, or modify more cases this becomes more challenging)
- add support for Enums
- generate `ShowID` function with 3 implementation variations (the last one is the rarest one and can be postponed)
  - use `Show` function with Like
  - use Show without any options and filter with Go for + if
  - in some cases we could need more filters -> see alerts.go (but we can implement it later)
- handle arrays
- handle more validation types
- write new `valueSet` function (see validations.go) that will have better defaults or more parameters that will determine 
checking behaviour which should get rid of edge cases that may cause bugs in the future
   - right now, we have `valueSet` function that doesn't take into consideration edge cases, e.g. with slice where sometimes
   we would like to do something like `alter x set y = ()` (set empty array to unset `y`). Those edge cases have cause on our
   validation, and it determines sometimes if we'll return an error or not, which can lead to bugs!
- refactor generation of `Describe`, so it will tak context and request as arguments
  - all the interface functions should have context and request as arguments for the sake of API consistency and generation simplicity
- split templates into multiple templates (e.g. ImplementationTemplate) to improve readability
  - example implementation - https://go.dev/play/p/Cgt0sISlzwK
  - divide implementation templates for Show, Describe and others
- check if SelfIdentifier implementation is correct (mostly type, because it's derived from interface obj) by checking
if there's a resource with different types of identifiers across queries (e.g. Create <AccountObjectIdentifier>, Alter <SchemaObjectIdentifier>) 
- we should specify prefix / postfix standard for top-level items in _def.go files to avoid any conflicts in the package
- remove name argument from QueryStruct in the Operation, because Opt structs in the Operation will have name from op name + interface field and not query struct itself
- Derive field name from QueryStruct, e.g. see network_policies_def where we can remove "Set" field, but we have to make a convention of creating nested struct with
name pattern like <interface name><name> e.g. NetworkPoliciesSet or NetworkPolicySet, then we could automatically remove prefix and we'll name field with postfix, so "Set" in this case
- Add more operations (every operation ?) in the database_role_def.go example
- Divide into packages or add common prefix for similar files (e.g. struct_plain.go, struct_db.go or builders_keyword.go, builders_parameter.go)
- Make a clear division between DSL files and model files (etc. QueryStruct(DSL) and Field(Model)) and divide them into separate packages (?)
- Simplify ImplementationTemplate (templates.go) and separate into multiple templates variables / definitions
- Add parameter to DtoTemplate (templates.go) to generate the right path to the dto generator's main.go file
- Right now to avoid generated structs duplication, arrays containing struct names have been introduced (template_executors.go),
find a better solution to solve the issue (add more logic to the templates ?)

##### Improvements
- automatic names of nested `struct`s (e.g. `DatabaseRoleRename`)
- check if generating with package name + invoking format removes unnecessary qualifier
- consider merging templates `StructTemplate` and `OptionsTemplate` (requires moving Doc to Field)
- expand unit tests generation
- experiment with Snowflake table (any table) representation in Go in order to implement DbStruct -> PlainStruct convert function
  - see if *string can have similar effect as sql.NullString (check go-snowflake connector ?)
     - if yes, then we should be using pointers instead of abstractions like sql.NullString and we can
     modify ShowMapping and DescribeMapping to generate convert function with automatic conversion (as we have in DTOs).
     warehouses.go is a good place to start with when planning mapping strategy, because there's a lot of different mapping cases.
- when calling .SelfIdentifier we can implicitly also add validateObjectIdentifier validation rule
- enforce user to use KindOf... functions with interface
  - example implementation - StringTyper implements Typer and all the KindOf... functions use StringTyper to return Typer easily - https://go.dev/play/p/TZZgSkkHw_M

##### Known issues
- generating two converts when Show and Desc use the same data structure
- wrong generated validations for validIdentifierIfSet for cases like
```go
A := QueryStruct("A").
	Identifier("Name").
	Validation(ValidIdentifierIfSet, "Name")
B := QueryStruct("B")
    .ListQueryStructField(A) // A []A - validations will be wrong because this is array
```
- cannot re-generate when client.go is using generated interface
- spaces in templates (especially nested validations)
- request mapping fails (`.toOpts()`) when nested object is not optional (pointer) e.g.
```go
type NestedReq struct {
}

type SomeReq struct {
    NestedReq NestedReq // Not a pointer and in toOpts right now we're always do a check if req.NestedReq != nil which is not correct for non pointer type
}
```

##### Known limitations
- automatic array conversion is not recursive, so we're only supporting one level mapping
  - []Request1{ foo Request2, bar int } won't be converted, but []Request1{ foo string, bar int } will