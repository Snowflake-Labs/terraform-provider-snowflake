package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateConnectionOptions] = new(CreateConnectionRequest)
	_ optionsProvider[AlterConnectionOptions]  = new(AlterConnectionRequest)
	_ optionsProvider[DropConnectionOptions]   = new(DropConnectionRequest)
	_ optionsProvider[ShowConnectionOptions]   = new(ShowConnectionRequest)
)

type CreateConnectionRequest struct {
	IfNotExists *bool
	name        AccountObjectIdentifier // required
	AsReplicaOf *AsReplicaOfRequest
	Comment     *string
}

type AsReplicaOfRequest struct {
	AsReplicaOf ExternalObjectIdentifier // required
}

type AlterConnectionRequest struct {
	IfExists                  *bool
	name                      AccountObjectIdentifier // required
	EnableConnectionFailover  *EnableConnectionFailoverRequest
	DisableConnectionFailover *DisableConnectionFailoverRequest
	Primary                   *bool
	Set                       *SetConnectionRequest
	Unset                     *UnsetConnectionRequest
}

type EnableConnectionFailoverRequest struct {
	ToAccounts []AccountIdentifier
}

type DisableConnectionFailoverRequest struct {
	ToAccounts *ToAccountsRequest
}

type ToAccountsRequest struct {
	Accounts []AccountIdentifier
}

type SetConnectionRequest struct {
	Comment *string
}

type UnsetConnectionRequest struct {
	Comment *bool
}

type DropConnectionRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowConnectionRequest struct {
	Like *Like
}
