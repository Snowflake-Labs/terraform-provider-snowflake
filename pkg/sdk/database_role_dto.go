package sdk

var (
	_ optionsProvider[createDatabaseRoleOptions] = new(CreateDatabaseRoleRequest)
	_ optionsProvider[alterDatabaseRoleOptions]  = new(AlterDatabaseRoleRequest)
	_ optionsProvider[dropDatabaseRoleOptions]   = new(DropDatabaseRoleRequest)
	_ optionsProvider[showDatabaseRoleOptions]   = new(ShowDatabaseRoleRequest)
)

type CreateDatabaseRoleRequest struct {
	orReplace   bool
	ifNotExists bool
	name        DatabaseObjectIdentifier // required
	comment     *string
}

type AlterDatabaseRoleRequest struct {
	ifExists bool
	name     DatabaseObjectIdentifier // required

	// One of
	rename *DatabaseRoleRenameRequest
	set    *DatabaseRoleSetRequest
	unset  *DatabaseRoleUnsetRequest
}

type DatabaseRoleRenameRequest struct {
	name DatabaseObjectIdentifier // required
}

type DatabaseRoleSetRequest struct {
	comment string // required
}

type DatabaseRoleUnsetRequest struct{}

type DropDatabaseRoleRequest struct {
	ifExists bool
	name     DatabaseObjectIdentifier // required
}

type ShowDatabaseRoleRequest struct {
	like     *Like
	database AccountObjectIdentifier // required
}
