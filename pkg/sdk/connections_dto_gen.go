package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateConnectionOptions]           = new(CreateConnectionRequest)
	_ optionsProvider[CreateReplicatedConnectionOptions] = new(CreateReplicatedConnectionRequest)
	_ optionsProvider[AlterFailoverConnectionOptions]    = new(AlterFailoverConnectionRequest)
	_ optionsProvider[AlterConnectionOptions]            = new(AlterConnectionRequest)
	_ optionsProvider[DropConnectionOptions]             = new(DropConnectionRequest)
	_ optionsProvider[ShowConnectionOptions]             = new(ShowConnectionRequest)
)

type CreateConnectionRequest struct {
	IfNotExists *bool
	name        AccountObjectIdentifier // required
	Comment     *string
}

type CreateReplicatedConnectionRequest struct {
	IfNotExists *bool
	name        AccountObjectIdentifier  // required
	ReplicaOf   ExternalObjectIdentifier // required
	Comment     *string
}

type AlterFailoverConnectionRequest struct {
	name                      AccountObjectIdentifier // required
	EnableConnectionFailover  *EnableConnectionFailoverRequest
	DisableConnectionFailover *DisableConnectionFailoverRequest
	Primary                   *bool
}

type EnableConnectionFailoverRequest struct {
	ToAccounts         []AccountIdentifier
	IgnoreEditionCheck *bool
}

type DisableConnectionFailoverRequest struct {
	ToAccounts *bool
	Accounts   []AccountIdentifier
}

type AlterConnectionRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
	Set      *SetRequest
	Unset    *UnsetRequest
}

type SetRequest struct {
	Comment *string
}

type UnsetRequest struct {
	Comment *bool
}

type DropConnectionRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowConnectionRequest struct {
	Like *Like
}
