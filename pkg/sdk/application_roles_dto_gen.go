package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[GrantApplicationRoleOptions]  = new(GrantApplicationRoleRequest)
	_ optionsProvider[RevokeApplicationRoleOptions] = new(RevokeApplicationRoleRequest)
	_ optionsProvider[ShowApplicationRoleOptions]   = new(ShowApplicationRoleRequest)
)

type GrantApplicationRoleRequest struct {
	name DatabaseObjectIdentifier // required
	To   KindOfRoleRequest
}

type KindOfRoleRequest struct {
	RoleName            *AccountObjectIdentifier
	ApplicationRoleName *DatabaseObjectIdentifier
	ApplicationName     *AccountObjectIdentifier
}

type RevokeApplicationRoleRequest struct {
	name DatabaseObjectIdentifier // required
	From KindOfRoleRequest
}

type ShowApplicationRoleRequest struct {
	ApplicationName AccountObjectIdentifier
	Limit           *LimitFrom
}
