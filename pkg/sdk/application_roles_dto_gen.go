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
	name        AccountObjectIdentifier // required
	Comment     *string
}

type AlterApplicationRoleRequest struct {
	IfExists     *bool
	name         AccountObjectIdentifier // required
	RenameTo     *AccountObjectIdentifier
	SetComment   *string
	UnsetComment *bool
}

type DropApplicationRoleRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowApplicationRoleRequest struct {
	ApplicationName AccountObjectIdentifier
	LimitFrom       *LimitFromApplicationRoleRequest
}

type LimitFromApplicationRoleRequest struct {
	Rows int // required
	From *string
}

type GrantApplicationRoleRequest struct {
	name    AccountObjectIdentifier        // required
	GrantTo ApplicationGrantOptionsRequest // required
}

type ApplicationGrantOptionsRequest struct {
	ParentRole      *AccountObjectIdentifier
	ApplicationRole *AccountObjectIdentifier
	Application     *AccountObjectIdentifier
}

type RevokeApplicationRoleRequest struct {
	name       AccountObjectIdentifier        // required
	RevokeFrom ApplicationGrantOptionsRequest // required
}
