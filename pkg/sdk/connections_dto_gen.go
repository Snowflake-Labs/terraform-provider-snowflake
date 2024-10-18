package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateConnectionConnectionOptions]           = new(CreateConnectionConnectionRequest)
	_ optionsProvider[CreateReplicatedConnectionConnectionOptions] = new(CreateReplicatedConnectionConnectionRequest)
	_ optionsProvider[AlterConnectionFailoverConnectionOptions]    = new(AlterConnectionFailoverConnectionRequest)
	_ optionsProvider[AlterConnectionConnectionOptions]            = new(AlterConnectionConnectionRequest)
)

type CreateConnectionConnectionRequest struct {
	IfNotExists *bool
	name        AccountObjectIdentifier // required
	Comment     *string
}

type CreateReplicatedConnectionConnectionRequest struct {
	IfNotExists *bool
	name        AccountObjectIdentifier  // required
	ReplicaOf   ExternalObjectIdentifier // required
	Comment     *string
}

type AlterConnectionFailoverConnectionRequest struct {
	name                      AccountObjectIdentifier // required
	EnableConnectionFailover  *EnableConnectionFailoverRequest
	DisableConnectionFailover *DisableConnectionFailoverRequest
	Primary                   *PrimaryRequest
}

type EnableConnectionFailoverRequest struct {
	Accounts           []ExternalObjectIdentifier
	IgnoreEditionCheck *bool
}

type DisableConnectionFailoverRequest struct {
	ToAccounts *bool
	Accounts   []ExternalObjectIdentifier
}

type PrimaryRequest struct {
}

type AlterConnectionConnectionRequest struct {
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
