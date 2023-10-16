package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	_ validatable = new(CreateAccountOptions)
	_ validatable = new(AlterAccountOptions)
	_ validatable = new(ShowAccountOptions)
)

type Accounts interface {
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateAccountOptions) error
	Alter(ctx context.Context, opts *AlterAccountOptions) error
	Show(ctx context.Context, opts *ShowAccountOptions) ([]Account, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Account, error)
	Drop(ctx context.Context, id AccountObjectIdentifier, gracePeriodInDays int, opts *DropAccountOptions) error
	Undrop(ctx context.Context, id AccountObjectIdentifier) error
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

// CreateAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-account.
type CreateAccountOptions struct {
	create  bool                    `ddl:"static" sql:"CREATE"`
	account bool                    `ddl:"static" sql:"ACCOUNT"`
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
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if opts.AdminName == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "AdminName"))
	}
	if !anyValueSet(opts.AdminPassword, opts.AdminRSAPublicKey) {
		errs = append(errs, errAtLeastOneOf("CreateAccountOptions", "AdminPassword", "AdminRSAPublicKey"))
	}
	if opts.Email == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "Email"))
	}
	if opts.Edition == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "Edition"))
	}
	return errors.Join(errs...)
}

func (c *accounts) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateAccountOptions) error {
	if opts == nil {
		opts = &CreateAccountOptions{}
	}
	opts.name = id
	return validateAndExec(c.client, ctx, opts)
}

// AlterAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-account.
type AlterAccountOptions struct {
	alter   bool `ddl:"static" sql:"ALTER"`
	account bool `ddl:"static" sql:"ACCOUNT"`

	Set    *AccountSet    `ddl:"keyword" sql:"SET"`
	Unset  *AccountUnset  `ddl:"list,no_parentheses" sql:"UNSET"`
	Rename *AccountRename `ddl:"-"`
	Drop   *AccountDrop   `ddl:"-"`
}

func (opts *AlterAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.Drop, opts.Rename) {
		errs = append(errs, errExactlyOneOf("CreateAccountOptions", "Set", "Unset", "Drop", "Rename"))
	}
	if valueSet(opts.Set) {
		errs = append(errs, opts.Set.validate())
	}
	if valueSet(opts.Unset) {
		errs = append(errs, opts.Unset.validate())
	}
	if valueSet(opts.Drop) {
		errs = append(errs, opts.Drop.validate())
	}
	if valueSet(opts.Rename) {
		errs = append(errs, opts.Rename.validate())
	}
	return errors.Join(errs...)
}

type AccountLevelParameters struct {
	AccountParameters *AccountParameters `ddl:"list,no_parentheses"`
	SessionParameters *SessionParameters `ddl:"list,no_parentheses"`
	ObjectParameters  *ObjectParameters  `ddl:"list,no_parentheses"`
	UserParameters    *UserParameters    `ddl:"list,no_parentheses"`
}

func (opts *AccountLevelParameters) validate() error {
	var errs []error
	if valueSet(opts.AccountParameters) {
		if err := opts.AccountParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.SessionParameters) {
		if err := opts.SessionParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.ObjectParameters) {
		if err := opts.ObjectParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.UserParameters) {
		if err := opts.UserParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type AccountSet struct {
	Parameters      *AccountLevelParameters `ddl:"list,no_parentheses"`
	ResourceMonitor AccountObjectIdentifier `ddl:"identifier,equals" sql:"RESOURCE_MONITOR"`
	PasswordPolicy  SchemaObjectIdentifier  `ddl:"identifier" sql:"PASSWORD POLICY"`
	SessionPolicy   SchemaObjectIdentifier  `ddl:"identifier" sql:"SESSION POLICY"`
	Tag             []TagAssociation        `ddl:"keyword" sql:"TAG"`
}

func (opts *AccountSet) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.Parameters, opts.ResourceMonitor, opts.PasswordPolicy, opts.SessionPolicy, opts.Tag) {
		errs = append(errs, errExactlyOneOf("AccountSet", "Parameters", "ResourceMonitor", "PasswordPolicy", "SessionPolicy", "Tag"))
	}
	if valueSet(opts.Parameters) {
		if err := opts.Parameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type AccountLevelParametersUnset struct {
	AccountParameters *AccountParametersUnset `ddl:"list,no_parentheses"`
	SessionParameters *SessionParametersUnset `ddl:"list,no_parentheses"`
	ObjectParameters  *ObjectParametersUnset  `ddl:"list,no_parentheses"`
	UserParameters    *UserParametersUnset    `ddl:"list,no_parentheses"`
}

func (opts *AccountLevelParametersUnset) validate() error {
	if !anyValueSet(opts.AccountParameters, opts.SessionParameters, opts.ObjectParameters, opts.UserParameters) {
		return errAtLeastOneOf("AccountLevelParametersUnset", "AccountParameters", "SessionParameters", "ObjectParameters", "UserParameters")
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
	var errs []error
	if !exactlyOneValueSet(opts.Parameters, opts.PasswordPolicy, opts.SessionPolicy, opts.Tag) {
		errs = append(errs, errExactlyOneOf("AccountUnset", "Parameters", "PasswordPolicy", "SessionPolicy", "Tag"))
	}
	if valueSet(opts.Parameters) {
		if err := opts.Parameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type AccountRename struct {
	Name       AccountObjectIdentifier `ddl:"identifier"`
	NewName    AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SaveOldURL *bool                   `ddl:"parameter" sql:"SAVE_OLD_URL"`
}

func (opts *AccountRename) validate() error {
	var errs []error
	if !ValidObjectIdentifier(opts.Name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.NewName) {
		errs = append(errs, errInvalidIdentifier("AccountRename", "NewName"))
	}
	return errors.Join(errs...)
}

type AccountDrop struct {
	Name   AccountObjectIdentifier `ddl:"identifier"`
	OldURL *bool                   `ddl:"keyword" sql:"DROP OLD URL"`
}

func (opts *AccountDrop) validate() error {
	var errs []error
	if !ValidObjectIdentifier(opts.Name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.OldURL) {
		// TODO: Should this really be validated to be true ?
		if !*opts.OldURL {
			errs = append(errs, fmt.Errorf("OldURL must be true"))
		}
	} else {
		errs = append(errs, errNotSet("AccountDrop", "OldURL"))
	}
	return errors.Join(errs...)
}

func (c *accounts) Alter(ctx context.Context, opts *AlterAccountOptions) error {
	if opts == nil {
		opts = &AlterAccountOptions{}
	}
	return validateAndExec(c.client, ctx, opts)
}

// ShowAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-organisation-accounts.
type ShowAccountOptions struct {
	show     bool  `ddl:"static" sql:"SHOW"`
	accounts bool  `ddl:"static" sql:"ORGANIZATION ACCOUNTS"`
	Like     *Like `ddl:"keyword" sql:"LIKE"`
}

func (opts *ShowAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
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
	Comment                              sql.NullString `db:"comment"`
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

func (row accountDBRow) convert() *Account {
	acc := &Account{
		OrganizationName:                     row.OrganizationName,
		AccountName:                          row.AccountName,
		RegionGroup:                          "",
		SnowflakeRegion:                      row.SnowflakeRegion,
		Edition:                              AccountEdition(row.Edition),
		AccountURL:                           row.AccountURL,
		CreatedOn:                            row.CreatedOn,
		Comment:                              row.Comment.String,
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

func (c *accounts) Show(ctx context.Context, opts *ShowAccountOptions) ([]Account, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[accountDBRow](c.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[accountDBRow, Account](dbRows)
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
			return &account, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

// DropAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-account.
type DropAccountOptions struct {
	drop              bool                    `ddl:"static" sql:"DROP"`
	account           bool                    `ddl:"static" sql:"ACCOUNT"`
	IfExists          *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name              AccountObjectIdentifier `ddl:"identifier"`
	gracePeriodInDays int                     `ddl:"parameter" sql:"GRACE_PERIOD_IN_DAYS"`
}

func (opts *DropAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !validateIntGreaterThanOrEqual(opts.gracePeriodInDays, 3) {
		errs = append(errs, errIntValue("DropAccountOptions", "gracePeriodInDays", IntErrGreaterOrEqual, 3))
	}
	return errors.Join(errs...)
}

func (c *accounts) Drop(ctx context.Context, id AccountObjectIdentifier, gracePeriodInDays int, opts *DropAccountOptions) error {
	if opts == nil {
		opts = &DropAccountOptions{}
	}
	opts.name = id
	opts.gracePeriodInDays = gracePeriodInDays
	return validateAndExec(c.client, ctx, opts)
}

// undropAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/undrop-account.
type undropAccountOptions struct {
	undrop  bool                    `ddl:"static" sql:"UNDROP"`
	account bool                    `ddl:"static" sql:"ACCOUNT"`
	name    AccountObjectIdentifier `ddl:"identifier"`
}

func (c *accounts) Undrop(ctx context.Context, id AccountObjectIdentifier) error {
	opts := &undropAccountOptions{
		name: id,
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, sql)
	return err
}
