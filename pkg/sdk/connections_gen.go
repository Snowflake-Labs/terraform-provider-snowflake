package sdk

import "context"

type Conntections interface {
	CreateConnection(ctx context.Context, request *CreateConnectionConnectionRequest) error
	CreateReplicatedConnection(ctx context.Context, request *CreateReplicatedConnectionConnectionRequest) error
	AlterConnectionFailover(ctx context.Context, request *AlterConnectionFailoverConnectionRequest) error
	AlterConnection(ctx context.Context, request *AlterConnectionConnectionRequest) error
}

// CreateConnectionConnectionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-connection.
type CreateConnectionConnectionOptions struct {
	create      bool                    `ddl:"static" sql:"CREATE"`
	connection  bool                    `ddl:"static" sql:"CONNECTION"`
	IfNotExists *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier `ddl:"identifier"`
	Comment     *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// CreateReplicatedConnectionConnectionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-connection.
type CreateReplicatedConnectionConnectionOptions struct {
	create      bool                     `ddl:"static" sql:"CREATE"`
	connection  bool                     `ddl:"static" sql:"CONNECTION"`
	IfNotExists *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier  `ddl:"identifier"`
	asReplicaOf bool                     `ddl:"static" sql:"AS REPLICA OF"`
	ReplicaOf   ExternalObjectIdentifier `ddl:"identifier"`
	Comment     *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterConnectionFailoverConnectionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-connection.
type AlterConnectionFailoverConnectionOptions struct {
	alter                     bool                       `ddl:"static" sql:"ALTER"`
	connection                bool                       `ddl:"static" sql:"CONNECTION"`
	name                      AccountObjectIdentifier    `ddl:"identifier"`
	EnableConnectionFailover  *EnableConnectionFailover  `ddl:"keyword" sql:"ENABLE FAILOVER TO ACCOUNTS"`
	DisableConnectionFailover *DisableConnectionFailover `ddl:"keyword" sql:"DISABLE FAILOVER"`
	Primary                   *Primary                   `ddl:"keyword"`
}
type EnableConnectionFailover struct {
	Accounts           []ExternalObjectIdentifier `ddl:"list,no_parentheses"`
	IgnoreEditionCheck *bool                      `ddl:"keyword" sql:"IGNORE EDITION CHECK"`
}
type DisableConnectionFailover struct {
	ToAccounts *bool                      `ddl:"keyword" sql:"TO ACCOUNTS"`
	Accounts   []ExternalObjectIdentifier `ddl:"list,no_parentheses"`
}
type Primary struct {
	primary bool `ddl:"static" sql:"PRIMARY"`
}

// AlterConnectionConnectionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-connection.
type AlterConnectionConnectionOptions struct {
	alter      bool                    `ddl:"static" sql:"ALTER"`
	connection bool                    `ddl:"static" sql:"CONNECTION"`
	IfExists   *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name       AccountObjectIdentifier `ddl:"identifier"`
	Set        *Set                    `ddl:"keyword" sql:"SET"`
	Unset      *Unset                  `ddl:"keyword" sql:"UNSET"`
}
type Set struct {
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}
type Unset struct {
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}
