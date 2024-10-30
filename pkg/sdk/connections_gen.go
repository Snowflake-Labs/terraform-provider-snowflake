package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Connections interface {
	Create(ctx context.Context, request *CreateConnectionRequest) error
	Alter(ctx context.Context, request *AlterConnectionRequest) error
	Drop(ctx context.Context, request *DropConnectionRequest) error
	Show(ctx context.Context, request *ShowConnectionRequest) ([]Connection, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Connection, error)
}

// CreateConnectionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-connection.
type CreateConnectionOptions struct {
	create      bool                      `ddl:"static" sql:"CREATE"`
	connection  bool                      `ddl:"static" sql:"CONNECTION"`
	IfNotExists *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier   `ddl:"identifier"`
	AsReplicaOf *ExternalObjectIdentifier `ddl:"identifier" sql:"AS REPLICA OF"`
	Comment     *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterConnectionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-connection.
type AlterConnectionOptions struct {
	alter                     bool                       `ddl:"static" sql:"ALTER"`
	connection                bool                       `ddl:"static" sql:"CONNECTION"`
	IfExists                  *bool                      `ddl:"keyword" sql:"IF EXISTS"`
	name                      AccountObjectIdentifier    `ddl:"identifier"`
	EnableConnectionFailover  *EnableConnectionFailover  `ddl:"keyword" sql:"ENABLE FAILOVER TO ACCOUNTS"`
	DisableConnectionFailover *DisableConnectionFailover `ddl:"keyword" sql:"DISABLE FAILOVER"`
	Primary                   *bool                      `ddl:"keyword" sql:"PRIMARY"`
	Set                       *Set                       `ddl:"keyword" sql:"SET"`
	Unset                     *Unset                     `ddl:"keyword" sql:"UNSET"`
}
type EnableConnectionFailover struct {
	ToAccounts []AccountIdentifier `ddl:"list,no_parentheses"`
}
type DisableConnectionFailover struct {
	ToAccounts *ToAccounts `ddl:"keyword" sql:"TO ACCOUNTS"`
}
type ToAccounts struct {
	Accounts []AccountIdentifier `ddl:"list,no_parentheses"`
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
	RegionGroup               sql.NullString `db:"region_group"`
	SnowflakeRegion           string         `db:"snowflake_region"`
	CreatedOn                 time.Time      `db:"created_on"`
	AccountName               string         `db:"account_name"`
	Name                      string         `db:"name"`
	Comment                   sql.NullString `db:"comment"`
	IsPrimary                 string         `db:"is_primary"`
	Primary                   string         `db:"primary"`
	FailoverAllowedToAccounts string         `db:"failover_allowed_to_accounts"`
	ConnectionUrl             string         `db:"connection_url"`
	OrganizationName          string         `db:"organization_name"`
	AccountLocator            string         `db:"account_locator"`
}
type Connection struct {
	RegionGroup               *string
	SnowflakeRegion           string
	CreatedOn                 time.Time
	AccountName               string
	Name                      string
	Comment                   *string
	IsPrimary                 bool
	Primary                   string
	FailoverAllowedToAccounts []string
	ConnectionUrl             string
	OrganizationName          string
	AccountLocator            string
}

func (c *Connection) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(c.Name)
}

func (c *Connection) ObjectType() ObjectType {
	return ObjectTypeConnection
}
