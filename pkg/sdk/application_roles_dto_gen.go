package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateApplicationRoleOptions] = new(CreateApplicationRoleRequest)
	_ optionsProvider[AlterApplicationRoleOptions]  = new(AlterApplicationRoleRequest)
	_ optionsProvider[DropApplicationRoleOptions]   = new(DropApplicationRoleRequest)
	_ optionsProvider[ShowApplicationRoleOptions]   = new(ShowApplicationRoleRequest)
	_ optionsProvider[GrantApplicationRoleOptions]  = new(GrantApplicationRoleRequest)
	_ optionsProvider[RevokeApplicationRoleOptions] = new(RevokeApplicationRoleRequest)
)

type CreateApplicationRoleRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        DatabaseObjectIdentifier // required
	Comment     *string
}

type AlterApplicationRoleRequest struct {
	IfExists     *bool
	name         DatabaseObjectIdentifier // required
	RenameTo     *DatabaseObjectIdentifier
	SetComment   *string
	UnsetComment *bool
}

type DropApplicationRoleRequest struct {
	IfExists *bool
	name     DatabaseObjectIdentifier // required
}

type ShowByIDApplicationRoleRequest struct {
	name            DatabaseObjectIdentifier // required
	ApplicationName AccountObjectIdentifier  // required
}

type ShowApplicationRoleRequest struct {
	ApplicationName AccountObjectIdentifier // required
	Limit           *LimitFrom
}

type GrantApplicationRoleRequest struct {
	name    DatabaseObjectIdentifier       // required
	GrantTo ApplicationGrantOptionsRequest // required
}

type ApplicationGrantOptionsRequest struct {
	ParentRole      *AccountObjectIdentifier
	ApplicationRole *DatabaseObjectIdentifier
	Application     *AccountObjectIdentifier
}

type RevokeApplicationRoleRequest struct {
	name       DatabaseObjectIdentifier       // required
	RevokeFrom ApplicationGrantOptionsRequest // required
}
