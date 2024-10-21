package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Connections interface {
	CreateConnection(ctx context.Context, request *CreateConnectionConnectionRequest) error
	CreateReplicatedConnection(ctx context.Context, request *CreateReplicatedConnectionConnectionRequest) error
	AlterConnectionFailover(ctx context.Context, request *AlterConnectionFailoverConnectionRequest) error
	AlterConnection(ctx context.Context, request *AlterConnectionConnectionRequest) error
	Drop(ctx context.Context, request *DropConnectionRequest) error
	Show(ctx context.Context, request *ShowConnectionRequest) ([]Connection, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Connection, error)
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

// DropConnectionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-connection.
type DropConnectionOptions struct {
	drop       bool                    `ddl:"static" sql:"DROP"`
	connection bool                    `ddl:"static" sql:"CONNECTION"`
	IfExists   *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name       AccountObjectIdentifier `ddl:"identifier"`
}

// ShowConnectionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-connections.
type ShowConnectionOptions struct {
	show        bool  `ddl:"static" sql:"SHOW"`
	connections bool  `ddl:"static" sql:"CONNECTIONS"`
	Like        *Like `ddl:"keyword" sql:"LIKE"`
}
type connectionRow struct {
	SnowflakeRegion           string         `db:"snowflake_region"`
	CreatedOn                 time.Time      `db:"created_on"`
	AccountName               string         `db:"account_name"`
	Name                      string         `db:"name"`
	Comment                   sql.NullString `db:"comment"`
	IsPrimary                 string         `db:"is_primary"`
	Primary                   string         `db:"primary"`
	FailoverAllowedToAccounts string         `db:"failover_allowed_to_accounts"`
	ConnectionUrl             string         `db:"connection_url"`
	OrgnizationName           string         `db:"orgnization_name"`
	AccountLocator            string         `db:"account_locator"`
}
type Connection struct {
	SnowflakeRegion           string
	CreatedOn                 time.Time
	AccountName               string
	Name                      string
	Comment                   *string
	IsPrimary                 bool
	Primary                   string
	FailoverAllowedToAccounts []string
	ConnectionUrl             string
	OrgnizationName           string
	AccountLocator            string
}
