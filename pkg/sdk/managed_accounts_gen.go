package sdk

import (
	"context"
	"database/sql"
)

type ManagedAccounts interface {
	Create(ctx context.Context, request *CreateManagedAccountRequest) error
	Drop(ctx context.Context, request *DropManagedAccountRequest) error
	Show(ctx context.Context, request *ShowManagedAccountRequest) ([]ManagedAccount, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ManagedAccount, error)
}

// CreateManagedAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-managed-account.
type CreateManagedAccountOptions struct {
	create                     bool                       `ddl:"static" sql:"CREATE"`
	managedAccount             bool                       `ddl:"static" sql:"MANAGED ACCOUNT"`
	name                       AccountObjectIdentifier    `ddl:"identifier"`
	CreateManagedAccountParams CreateManagedAccountParams `ddl:"list,no_parentheses"`
}

type CreateManagedAccountParams struct {
	AdminName     string  `ddl:"parameter,single_quotes" sql:"ADMIN_NAME"`
	AdminPassword string  `ddl:"parameter,single_quotes" sql:"ADMIN_PASSWORD"`
	typeProvider  string  `ddl:"static" sql:"TYPE = READER"`
	Comment       *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// DropManagedAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-managed-account.
type DropManagedAccountOptions struct {
	drop           bool                    `ddl:"static" sql:"DROP"`
	managedAccount bool                    `ddl:"static" sql:"MANAGED ACCOUNT"`
	name           AccountObjectIdentifier `ddl:"identifier"`
}

// ShowManagedAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-managed-accounts.
type ShowManagedAccountOptions struct {
	show            bool  `ddl:"static" sql:"SHOW"`
	managedAccounts bool  `ddl:"static" sql:"MANAGED ACCOUNTS"`
	Like            *Like `ddl:"keyword" sql:"LIKE"`
}

type managedAccountDBRow struct {
	Name              string         `db:"name"`
	Cloud             string         `db:"cloud"`
	Region            string         `db:"region"`
	Locator           string         `db:"locator"`
	CreatedOn         string         `db:"created_on"`
	Url               string         `db:"url"`
	AccountLocatorUrl string         `db:"account_locator_url"`
	IsReader          bool           `db:"is_reader"`
	Comment           sql.NullString `db:"comment"`
}

type ManagedAccount struct {
	Name              string
	Cloud             string
	Region            string
	Locator           string
	CreatedOn         string
	URL               string
	AccountLocatorURL string
	IsReader          bool
	Comment           string
}
