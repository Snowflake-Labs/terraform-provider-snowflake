package example

//go:generate go run ../../dto-builder-generator/main.go

var (
	_ optionsProvider[CreateDatabaseRoleOptions] = new(CreateDatabaseRoleRequest)
	_ optionsProvider[AlterDatabaseRoleOptions]  = new(AlterDatabaseRoleRequest)
)

type CreateDatabaseRoleRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        DatabaseObjectIdentifier // required
	Comment     *string
}

type AlterDatabaseRoleRequest struct {
	IfExists *bool
	name     DatabaseObjectIdentifier // required
	Rename   *DatabaseRoleRenameRequest
	Set      *DatabaseRoleSetRequest
	Unset    *DatabaseRoleUnsetRequest
}

type DatabaseRoleRenameRequest struct {
	Name DatabaseObjectIdentifier // required
}

type DatabaseRoleSetRequest struct {
	Comment          string // required
	NestedThirdLevel *NestedThirdLevelRequest
}

type NestedThirdLevelRequest struct {
	Field DatabaseObjectIdentifier // required
}

type DatabaseRoleUnsetRequest struct {
	Comment bool // required
}
