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
	AsReplicaOf *ExternalObjectIdentifier
	Comment     *string
}

type AlterConnectionRequest struct {
	IfExists                  *bool
	name                      AccountObjectIdentifier // required
	EnableConnectionFailover  *EnableConnectionFailoverRequest
	DisableConnectionFailover *DisableConnectionFailoverRequest
	Primary                   *bool
	Set                       *ConnectionSetRequest
	Unset                     *ConnectionUnsetRequest
}

type EnableConnectionFailoverRequest struct {
	ToAccounts []AccountIdentifier // required
}

type DisableConnectionFailoverRequest struct {
	ToAccounts *ToAccountsRequest
}

type ToAccountsRequest struct {
	Accounts []AccountIdentifier // required
}

type ConnectionSetRequest struct {
	Comment *string
}

type ConnectionUnsetRequest struct {
	Comment *bool
}

type DropConnectionRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowConnectionRequest struct {
	Like *Like
}
