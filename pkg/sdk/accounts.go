package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Accounts interface {
	// Create creates an account.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateAccountOptions) error
	// Alter modifies an existing account
	Alter(ctx context.Context, opts *AlterAccountOptions) error
	// Show returns a list of accounts.
	Show(ctx context.Context, opts *ShowAccountOptions) ([]*Account, error)
	// ShowByID returns an account by id
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Account, error)
}

var _ Accounts = (*accounts)(nil)

type accounts struct {
	client *Client
}

type AccountEdition string

var (
	EditionStandard         AccountEdition = "STANDARD"
	EditionEnterprise       AccountEdition = "ENTERPRISE"
	EditionBusinessCritical AccountEdition = "BUSINESS_CRITICAL"
)

type CreateAccountOptions struct {
	create  bool                    `ddl:"static" sql:"CREATE"`  //lint:ignore U1000 This is used in the ddl tag
	account bool                    `ddl:"static" sql:"ACCOUNT"` //lint:ignore U1000 This is used in the ddl tag
	name    AccountObjectIdentifier `ddl:"identifier"`

	// Object properties
	AdminName          string         `ddl:"parameter,single_quotes" sql:"ADMIN_NAME"`
	AdminPassword      *string        `ddl:"parameter,single_quotes" sql:"ADMIN_PASSWORD"`
	AdminRSAPublicKey  *string        `ddl:"parameter,single_quotes" sql:"ADMIN_RSA_PUBLIC_KEY"`
	FirstName          *string        `ddl:"parameter,single_quotes" sql:"FIRST_NAME"`
	LastName           *string        `ddl:"parameter,single_quotes" sql:"LAST_NAME"`
	Email              string         `ddl:"parameter,single_quotes" sql:"EMAIL"`
	MustChangePassword *bool          `ddl:"parameter" sql:"MUST_CHANGE_PASSWORD"`
	Edition            AccountEdition `ddl:"parameter" sql:"EDITION"`
	RegionGroup        *string        `ddl:"parameter,single_quotes" sql:"REGION_GROUP"`
	Region             *string        `ddl:"parameter,single_quotes" sql:"REGION"`
	Comment            *string        `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreateAccountOptions) validate() error {
	if opts.AdminName == "" {
		return fmt.Errorf("AdminName is required")
	}
	if !anyValueSet(opts.AdminPassword, opts.AdminRSAPublicKey) {
		return fmt.Errorf("at least one of AdminPassword or AdminRSAPublicKey must be set")
	}
	if opts.Email == "" {
		return fmt.Errorf("Email is required")
	}
	if opts.Edition == "" {
		return fmt.Errorf("Edition is required")
	}
	return nil
}

func (c *accounts) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateAccountOptions) error {
	if opts == nil {
		opts = &CreateAccountOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	stmt, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, stmt)
	return err
}

type AlterAccountOptions struct {
	alter   bool `ddl:"static" sql:"ALTER"`   //lint:ignore U1000 This is used in the ddl tag
	account bool `ddl:"static" sql:"ACCOUNT"` //lint:ignore U1000 This is used in the ddl tag

	Set    *AccountSet    `ddl:"keyword" sql:"SET"`
	Unset  *AccountUnset  `ddl:"list,no_parentheses" sql:"UNSET"`
	Rename *AccountRename `ddl:"-"`
	Drop   *AccountDrop   `ddl:"-"`
}

func (opts *AlterAccountOptions) validate() error {
	if ok := exactlyOneValueSet(
		opts.Set,
		opts.Unset,
		opts.Drop,
		opts.Rename); !ok {
		return fmt.Errorf("exactly one of Set, Unset, Drop, Rename  must be set")
	}
	if valueSet(opts.Set) {
		return opts.Set.validate()
	}
	if valueSet(opts.Unset) {
		return opts.Unset.validate()
	}
	if valueSet(opts.Drop) {
		return opts.Drop.validate()
	}
	if valueSet(opts.Rename) {
		return opts.Rename.validate()
	}
	return nil
}

type AccountLevelParameters struct {
	AccountParameters *AccountParameters `ddl:"list,no_parentheses"`
	SessionParameters *SessionParameters `ddl:"list,no_parentheses"`
	ObjectParameters  *ObjectParameters  `ddl:"list,no_parentheses"`
	UserParameters    *UserParameters    `ddl:"list,no_parentheses"`
}

func (opts *AccountLevelParameters) validate() error {
	if valueSet(opts.AccountParameters) {
		if err := opts.AccountParameters.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.SessionParameters) {
		if err := opts.SessionParameters.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.ObjectParameters) {
		if err := opts.ObjectParameters.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.UserParameters) {
		if err := opts.UserParameters.validate(); err != nil {
			return err
		}
	}
	return nil
}

type AccountSet struct {
	Parameters      *AccountLevelParameters `ddl:"list,no_parentheses"`
	ResourceMonitor AccountObjectIdentifier `ddl:"identifier,equals" sql:"RESOURCE_MONITOR"`
	PasswordPolicy  SchemaObjectIdentifier  `ddl:"identifier" sql:"PASSWORD POLICY"`
	SessionPolicy   SchemaObjectIdentifier  `ddl:"identifier" sql:"SESSION POLICY"`
	Tag             []TagAssociation        `ddl:"keyword" sql:"TAG"`
}

func (opts *AccountSet) validate() error {
	if !anyValueSet(opts.Parameters, opts.ResourceMonitor, opts.PasswordPolicy, opts.SessionPolicy, opts.Tag) {
		return fmt.Errorf("at least one of parameters, resource monitor, password policy, session policy, or tag must be set")
	}
	if valueSet(opts.Parameters) {
		if !everyValueNil(opts.ResourceMonitor, opts.PasswordPolicy, opts.SessionPolicy, opts.Tag) {
			return fmt.Errorf("cannot set both parameters and resource monitor, password policy, session policy, or tag")
		}
		return opts.Parameters.validate()
	}
	if valueSet(opts.ResourceMonitor) {
		if !everyValueNil(opts.PasswordPolicy, opts.SessionPolicy, opts.Tag) {
			return fmt.Errorf("cannot set both resource monitor and password policy, session policy, or tag")
		}
		return nil
	}
	if valueSet(opts.PasswordPolicy) {
		if !everyValueNil(opts.SessionPolicy, opts.Tag) {
			return fmt.Errorf("cannot set both password policy and session policy or tag")
		}
		return nil
	}
	if valueSet(opts.SessionPolicy) {
		if !everyValueNil(opts.Tag) {
			return fmt.Errorf("cannot set both session policy and tag")
		}
		return nil
	}
	return nil
}

type AccountLevelParametersUnset struct {
	AccountParameters *AccountParametersUnset `ddl:"list,no_parentheses"`
	SessionParameters *SessionParametersUnset `ddl:"list,no_parentheses"`
	ObjectParameters  *ObjectParametersUnset  `ddl:"list,no_parentheses"`
	UserParameters    *UserParametersUnset    `ddl:"list,no_parentheses"`
}

func (opts *AccountLevelParametersUnset) validate() error {
	if !anyValueSet(opts.AccountParameters, opts.SessionParameters, opts.ObjectParameters, opts.UserParameters) {
		return fmt.Errorf("at least one of account parameters, session parameters, object parameters, or user parameters must be set")
	}
	return nil
}

type AccountUnset struct {
	Parameters     *AccountLevelParametersUnset `ddl:"list,no_parentheses"`
	PasswordPolicy *bool                        `ddl:"keyword" sql:"PASSWORD POLICY"`
	SessionPolicy  *bool                        `ddl:"keyword" sql:"SESSION POLICY"`
	Tag            []ObjectIdentifier           `ddl:"keyword" sql:"TAG"`
}

func (opts *AccountUnset) validate() error {
	if !anyValueSet(opts.Parameters, opts.PasswordPolicy, opts.SessionPolicy, opts.Tag) {
		return fmt.Errorf("at least one of parameters, password policy, session policy, or tag must be set")
	}
	if valueSet(opts.Parameters) {
		if !everyValueNil(opts.PasswordPolicy, opts.SessionPolicy, opts.Tag) {
			return fmt.Errorf("cannot unset both parameters and password policy, session policy, or tag")
		}
		return opts.Parameters.validate()
	}
	if valueSet(opts.PasswordPolicy) {
		if !everyValueNil(opts.SessionPolicy, opts.Tag) {
			return fmt.Errorf("cannot unset both password policy and session policy or tag")
		}
		return nil
	}
	if valueSet(opts.SessionPolicy) {
		if !everyValueNil(opts.Tag) {
			return fmt.Errorf("cannot unset both session policy and tag")
		}
		return nil
	}
	return nil
}

type AccountRename struct {
	Name       AccountObjectIdentifier `ddl:"identifier"`
	NewName    AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SaveOldURL *bool                   `ddl:"parameter" sql:"SAVE_OLD_URL"`
}

func (opts *AccountRename) validate() error {
	if !validObjectidentifier(opts.Name) {
		return fmt.Errorf("Name must be set")
	}
	if !validObjectidentifier(opts.NewName) {
		return fmt.Errorf("NewName must be set")
	}
	return nil
}

type AccountDrop struct {
	Name   AccountObjectIdentifier `ddl:"identifier"`
	OldURL *bool                   `ddl:"keyword" sql:"DROP OLD URL"`
}

func (opts *AccountDrop) validate() error {
	if !validObjectidentifier(opts.Name) {
		return fmt.Errorf("Name must be set")
	}
	if valueSet(opts.OldURL) {
		if !*opts.OldURL {
			return fmt.Errorf("OldURL must be true")
		}
	} else {
		return fmt.Errorf("OldURL must be set")
	}
	return nil
}

func (c *accounts) Alter(ctx context.Context, opts *AlterAccountOptions) error {
	if opts == nil {
		opts = &AlterAccountOptions{}
	}
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, sql)
	return err
}

type ShowAccountOptions struct {
	show     bool  `ddl:"static" sql:"SHOW"`                  //lint:ignore U1000 This is used in the ddl tag
	accounts bool  `ddl:"static" sql:"ORGANIZATION ACCOUNTS"` //lint:ignore U1000 This is used in the ddl tag
	Like     *Like `ddl:"keyword" sql:"LIKE"`
}

func (opts *ShowAccountOptions) validate() error {
	return nil
}

type Account struct {
	OrganizationName                     string
	AccountName                          string
	RegionGroup                          string
	SnowflakeRegion                      string
	Edition                              AccountEdition
	AccountURL                           string
	CreatedOn                            time.Time
	Comment                              string
	AccountLocator                       string
	AccountLocatorURL                    string
	ManagedAccounts                      int
	ConsumptionBillingEntityName         string
	MarketplaceConsumerBillingEntityName string
	MarketplaceProviderBillingEntityName string
	OldAccountURL                        string
	IsOrgAdmin                           bool
}

func (v *Account) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.AccountName)
}

func (v *Account) AccountID() AccountIdentifier {
	return NewAccountIdentifier(v.OrganizationName, v.AccountName)
}

type accountDBRow struct {
	OrganizationName                     string         `db:"organization_name"`
	AccountName                          string         `db:"account_name"`
	RegionGroup                          sql.NullString `db:"region_group"`
	SnowflakeRegion                      string         `db:"snowflake_region"`
	Edition                              string         `db:"edition"`
	AccountURL                           string         `db:"account_url"`
	CreatedOn                            time.Time      `db:"created_on"`
	Comment                              string         `db:"comment"`
	AccountLocator                       string         `db:"account_locator"`
	AccountLocatorURL                    string         `db:"account_locator_url"`
	AccountOldURLSavedOn                 sql.NullString `db:"account_old_url_saved_on"`
	ManagedAccounts                      int            `db:"managed_accounts"`
	ConsumptionBillingEntityName         string         `db:"consumption_billing_entity_name"`
	MarketplaceConsumerBillingEntityName sql.NullString `db:"marketplace_consumer_billing_entity_name"`
	MarketplaceProviderBillingEntityName sql.NullString `db:"marketplace_provider_billing_entity_name"`
	OldAccountURL                        string         `db:"old_account_url"`
	IsOrgAdmin                           bool           `db:"is_org_admin"`
}

func (row accountDBRow) toAccount() *Account {
	acc := &Account{
		OrganizationName:                     row.OrganizationName,
		AccountName:                          row.AccountName,
		RegionGroup:                          "",
		SnowflakeRegion:                      row.SnowflakeRegion,
		Edition:                              AccountEdition(row.Edition),
		AccountURL:                           row.AccountURL,
		CreatedOn:                            row.CreatedOn,
		Comment:                              row.Comment,
		AccountLocator:                       row.AccountLocator,
		AccountLocatorURL:                    row.AccountLocatorURL,
		ManagedAccounts:                      row.ManagedAccounts,
		ConsumptionBillingEntityName:         row.ConsumptionBillingEntityName,
		MarketplaceConsumerBillingEntityName: "",
		MarketplaceProviderBillingEntityName: "",
		OldAccountURL:                        row.OldAccountURL,
		IsOrgAdmin:                           row.IsOrgAdmin,
	}
	if row.MarketplaceConsumerBillingEntityName.Valid {
		acc.MarketplaceConsumerBillingEntityName = row.MarketplaceConsumerBillingEntityName.String
	}
	if row.MarketplaceProviderBillingEntityName.Valid {
		acc.MarketplaceProviderBillingEntityName = row.MarketplaceProviderBillingEntityName.String
	}
	if row.RegionGroup.Valid {
		acc.SnowflakeRegion = row.RegionGroup.String
	}
	return acc
}

func (c *accounts) Show(ctx context.Context, opts *ShowAccountOptions) ([]*Account, error) {
	if opts == nil {
		opts = &ShowAccountOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []accountDBRow{}
	err = c.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]*Account, len(dest))
	for i, row := range dest {
		resultList[i] = row.toAccount()
	}

	return resultList, nil
}

func (c *accounts) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Account, error) {
	accounts, err := c.Show(ctx, &ShowAccountOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.AccountName == id.Name() || account.AccountLocator == id.Name() {
			return account, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}
