package sdk

import (
	"context"
	"fmt"
	"time"
)

type Accounts interface {
	// Create creates an account.
	Create(ctx context.Context, name AccountObjectIdentifier, opts *AccountCreateOptions) error
	// Alter modifies an existing account
	Alter(ctx context.Context, opts *AccountAlterOptions) error
	// Show returns a list of accounts.
	Show(ctx context.Context, opts *AccountShowOptions) ([]*Account, error)
	// ShowByID returns an account by ID
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

type AccountCreateOptions struct {
	create  bool                    `ddl:"static" db:"CREATE"`  //lint:ignore U1000 This is used in the ddl tag
	account bool                    `ddl:"static" db:"ACCOUNT"` //lint:ignore U1000 This is used in the ddl tag
	name    AccountObjectIdentifier `ddl:"identifier"`

	// Object properties
	AdminName          string         `ddl:"parameter,single_quotes" db:"ADMIN_NAME"`
	AdminPassword      *string        `ddl:"parameter,single_quotes" db:"ADMIN_PASSWORD"`
	AdminRSAPublicKey  *string        `ddl:"parameter,single_quotes" db:"ADMIN_RSA_PUBLIC_KEY"`
	FirstName          *string        `ddl:"parameter,single_quotes" db:"FIRST_NAME"`
	LastName           *string        `ddl:"parameter,single_quotes" db:"LAST_NAME"`
	Email              string         `ddl:"parameter,single_quotes" db:"EMAIL"`
	MustChangePassword *bool          `ddl:"parameter" db:"MUST_CHANGE_PASSWORD"`
	Edition            AccountEdition `ddl:"parameter" db:"EDITION"`
	RegionGroup        *string        `ddl:"parameter,single_quotes" db:"REGION_GROUP"`
	Region             *string        `ddl:"parameter,single_quotes" db:"REGION"`
	Comment            *string        `ddl:"parameter,single_quotes" db:"COMMENT"`
}

func (opts *AccountCreateOptions) validate() error {
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

func (c *accounts) Create(ctx context.Context, name AccountObjectIdentifier, opts *AccountCreateOptions) error {
	if opts == nil {
		opts = &AccountCreateOptions{}
	}
	opts.name = name
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

type AccountAlterOptions struct {
	alter   bool                     `ddl:"static" db:"ALTER"`   //lint:ignore U1000 This is used in the ddl tag
	account bool                     `ddl:"static" db:"ACCOUNT"` //lint:ignore U1000 This is used in the ddl tag
	Name    *AccountObjectIdentifier `ddl:"identifier"`

	Set   *AccountSet   `ddl:"keyword" db:"SET"`
	Unset *AccountUnset `ddl:"list,no_parentheses" db:"UNSET"`

	ResourceMonitor     AccountObjectIdentifier `ddl:"identifier,equals" db:"SET RESOURCE_MONITOR"`
	PasswordPolicy      SchemaObjectIdentifier  `ddl:"identifier" db:"SET PASSWORD POLICY"`
	SessionPolicy       SchemaObjectIdentifier  `ddl:"identifier" db:"SET SESSION POLICY"`
	UnsetPasswordPolicy *bool                   `ddl:"keyword" db:"UNSET PASSWORD POLICY"`
	UnsetSessionPolicy  *bool                   `ddl:"keyword" db:"UNSET SESSION POLICY"`

	SetTag   []TagAssociation   `ddl:"keyword" db:"SET TAG"`
	UnsetTag []ObjectIdentifier `ddl:"keyword" db:"UNSET TAG"`

	NewName    AccountObjectIdentifier `ddl:"identifier" db:"RENAME TO"`
	SaveOldURL *bool                   `ddl:"parameter" db:"SAVE_OLD_URL"`
	DropOldURL *bool                   `ddl:"keyword" db:"DROP OLD URL"`
}

func (opts *AccountAlterOptions) validate() error {
	if ok := exactlyOneValueSet(
		opts.Set,
		opts.Unset,
		opts.ResourceMonitor,
		opts.PasswordPolicy,
		opts.SessionPolicy,
		opts.UnsetPasswordPolicy,
		opts.UnsetSessionPolicy,
		opts.SetTag,
		opts.UnsetTag,
		opts.NewName,
		opts.DropOldURL); !ok {
		return fmt.Errorf("exactly one of Set, Unset, ResourceMonitor, PasswordPolicy, SessionPolicy, UnsetPasswordPolicy, UnsetSessionPolicy, SetTag, UnsetTag, NewName or DropOldURL must be set")
	}
	if (valueSet(opts.NewName) || valueSet(opts.DropOldURL)) && !valueSet(opts.Name) {
		return fmt.Errorf("Name must be set when using NewName or DropOldURL")
	}

	return nil
}

type AccountSet struct {
	// Account params
	AllowIdToken                                 *bool   `ddl:"parameter" db:"ALLOW_ID_TOKEN"`
	ClientEncryptionKeySize                      *int    `ddl:"parameter" db:"CLIENT_ENCRYPTION_KEY_SIZE"`
	EnableInternalStagesPrivatelink              *bool   `ddl:"parameter" db:"ENABLE_INTERNAL_STAGES_PRIVATELINK"`
	ExternalOauthAddPrivilegedRolesToBlockedList *bool   `ddl:"parameter" db:"EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"`
	InitialReplicationSizeLimitInTb              *int    `ddl:"parameter" db:"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"`
	NetworkPolicy                                *string `ddl:"parameter,single_quotes" db:"NETWORK_POLICY"`
	PeriodicDataRekeying                         *bool   `ddl:"parameter" db:"PERIODIC_DATA_REKEYING"`
	PreventUnloadToInlineUrl                     *bool   `ddl:"parameter" db:"PREVENT_UNLOAD_TO_INLINE_URL"`
	PreventUnloadToInternalStages                *bool   `ddl:"parameter" db:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	RequireStorageIntegrationForStageCreation    *bool   `ddl:"parameter" db:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"`
	RequireStorageIntegrationForStageOperation   *bool   `ddl:"parameter" db:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"`
	SsoLoginPage                                 *bool   `ddl:"parameter" db:"SSO_LOGIN_PAGE"`

	// User params
	EnableUnredactedQuerySyntaxError *bool `ddl:"parameter" db:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`

	// Object params
	DataRetentionTimeInDays         *int    `ddl:"parameter" db:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays      *int    `ddl:"parameter" db:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDdlCollation             *string `ddl:"parameter,single_quotes" db:"DEFAULT_DDL_COLLATION"`
	MaxConcurrencyLevel             *int    `ddl:"parameter" db:"MAX_CONCURRENCY_LEVEL"`
	PipeExecutionPaused             *bool   `ddl:"parameter" db:"PIPE_EXECUTION_PAUSED"`
	StatementQueuedTimeoutInSeconds *int    `ddl:"parameter" db:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *int    `ddl:"parameter" db:"STATEMENT_TIMEOUT_IN_SECONDS"`

	// Session params
	AbortDetachedQuery               *bool   `ddl:"parameter" db:"ABORT_DETACHED_QUERY"`
	Autocommit                       *bool   `ddl:"parameter" db:"AUTOCOMMIT"`
	BinaryInputFormat                *string `ddl:"parameter,single_quotes" db:"BINARY_INPUT_FORMAT"`
	BinaryOutputFormat               *string `ddl:"parameter,single_quotes" db:"BINARY_OUTPUT_FORMAT"`
	DateInputFormat                  *string `ddl:"parameter,single_quotes" db:"DATE_INPUT_FORMAT"`
	DateOutputFormat                 *string `ddl:"parameter,single_quotes" db:"DATE_OUTPUT_FORMAT"`
	ErrorOnNondeterministicMerge     *bool   `ddl:"parameter" db:"ERROR_ON_NONDETERMINISTIC_MERGE"`
	ErrorOnNondeterministicUpdate    *bool   `ddl:"parameter" db:"ERROR_ON_NONDETERMINISTIC_UPDATE"`
	JsonIndent                       *int    `ddl:"parameter" db:"JSON_INDENT"`
	LockTimeout                      *int    `ddl:"parameter" db:"LOCK_TIMEOUT"`
	QueryTag                         *string `ddl:"parameter,single_quotes" db:"QUERY_TAG"`
	RowsPerResultset                 *int    `ddl:"parameter" db:"ROWS_PER_RESULTSET"`
	SimulatedDataSharingConsumer     *string `ddl:"parameter,single_quotes" db:"SIMULATED_DATA_SHARING_CONSUMER"`
	StrictJsonOutput                 *bool   `ddl:"parameter" db:"STRICT_JSON_OUTPUT"`
	TimestampDayIsAlways24h          *bool   `ddl:"parameter" db:"TIMESTAMP_DAY_IS_ALWAYS_24H"`
	TimestampInputFormat             *string `ddl:"parameter,single_quotes" db:"TIMESTAMP_INPUT_FORMAT"`
	TimestampLtzOutputFormat         *string `ddl:"parameter,single_quotes" db:"TIMESTAMP_LTZ_OUTPUT_FORMAT"`
	TimestampNtzOutputFormat         *string `ddl:"parameter,single_quotes" db:"TIMESTAMP_NTZ_OUTPUT_FORMAT"`
	TimestampOutputFormat            *string `ddl:"parameter,single_quotes" db:"TIMESTAMP_OUTPUT_FORMAT"`
	TimestampTypeMapping             *string `ddl:"parameter,single_quotes" db:"TIMESTAMP_TYPE_MAPPING"`
	TimestampTzOutputFormat          *string `ddl:"parameter,single_quotes" db:"TIMESTAMP_TZ_OUTPUT_FORMAT"`
	Timezone                         *string `ddl:"parameter,single_quotes" db:"TIMEZONE"`
	TimeInputFormat                  *string `ddl:"parameter,single_quotes" db:"TIME_INPUT_FORMAT"`
	TimeOutputFormat                 *string `ddl:"parameter,single_quotes" db:"TIME_OUTPUT_FORMAT"`
	TransactionDefaultIsolationLevel *string `ddl:"parameter,single_quotes" db:"TRANSACTION_DEFAULT_ISOLATION_LEVEL"`
	TwoDigitCenturyStart             *int    `ddl:"parameter" db:"TWO_DIGIT_CENTURY_START"`
	UnsupportedDdlAction             *string `ddl:"parameter,single_quotes" db:"UNSUPPORTED_DDL_ACTION"`
	UseCachedResult                  *bool   `ddl:"parameter" db:"USE_CACHED_RESULT"`
	WeekOfYearPolicy                 *int    `ddl:"parameter" db:"WEEK_OF_YEAR_POLICY"`
	WeekStart                        *int    `ddl:"parameter" db:"WEEK_START"`
}

type AccountUnset struct {
	// Account params
	AllowIdToken                                 *bool `ddl:"keyword" db:"ALLOW_ID_TOKEN"`
	ClientEncryptionKeySize                      *bool `ddl:"keyword" db:"CLIENT_ENCRYPTION_KEY_SIZE"`
	EnableInternalStagesPrivatelink              *bool `ddl:"keyword" db:"ENABLE_INTERNAL_STAGES_PRIVATELINK"`
	ExternalOauthAddPrivilegedRolesToBlockedList *bool `ddl:"keyword" db:"EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"`
	InitialReplicationSizeLimitInTb              *bool `ddl:"keyword" db:"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"`
	NetworkPolicy                                *bool `ddl:"keyword" db:"NETWORK_POLICY"`
	PeriodicDataRekeying                         *bool `ddl:"keyword" db:"PERIODIC_DATA_REKEYING"`
	PreventUnloadToInlineUrl                     *bool `ddl:"keyword" db:"PREVENT_UNLOAD_TO_INLINE_URL"`
	PreventUnloadToInternalStages                *bool `ddl:"keyword" db:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	RequireStorageIntegrationForStageCreation    *bool `ddl:"keyword" db:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"`
	RequireStorageIntegrationForStageOperation   *bool `ddl:"keyword" db:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"`
	SsoLoginPage                                 *bool `ddl:"keyword" db:"SSO_LOGIN_PAGE"`

	// Object params
	DataRetentionTimeInDays         *bool `ddl:"keyword" db:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays      *bool `ddl:"keyword" db:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDdlCollation             *bool `ddl:"keyword" db:"DEFAULT_DDL_COLLATION"`
	MaxConcurrencyLevel             *bool `ddl:"keyword" db:"MAX_CONCURRENCY_LEVEL"`
	PipeExecutionPaused             *bool `ddl:"keyword" db:"PIPE_EXECUTION_PAUSED"`
	StatementQueuedTimeoutInSeconds *bool `ddl:"keyword" db:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds       *bool `ddl:"keyword" db:"STATEMENT_TIMEOUT_IN_SECONDS"`

	// Session params
	AbortDetachedQuery               *bool `ddl:"keyword" db:"ABORT_DETACHED_QUERY"`
	Autocommit                       *bool `ddl:"keyword" db:"AUTOCOMMIT"`
	BinaryInputFormat                *bool `ddl:"keyword" db:"BINARY_INPUT_FORMAT"`
	BinaryOutputFormat               *bool `ddl:"keyword" db:"BINARY_OUTPUT_FORMAT"`
	DateInputFormat                  *bool `ddl:"keyword" db:"DATE_INPUT_FORMAT"`
	DateOutputFormat                 *bool `ddl:"keyword" db:"DATE_OUTPUT_FORMAT"`
	ErrorOnNondeterministicMerge     *bool `ddl:"keyword" db:"ERROR_ON_NONDETERMINISTIC_MERGE"`
	ErrorOnNondeterministicUpdate    *bool `ddl:"keyword" db:"ERROR_ON_NONDETERMINISTIC_UPDATE"`
	JsonIndent                       *bool `ddl:"keyword" db:"JSON_INDENT"`
	LockTimeout                      *bool `ddl:"keyword" db:"LOCK_TIMEOUT"`
	QueryTag                         *bool `ddl:"keyword" db:"QUERY_TAG"`
	RowsPerResultset                 *bool `ddl:"keyword" db:"ROWS_PER_RESULTSET"`
	SimulatedDataSharingConsumer     *bool `ddl:"keyword" db:"SIMULATED_DATA_SHARING_CONSUMER"`
	StrictJsonOutput                 *bool `ddl:"keyword" db:"STRICT_JSON_OUTPUT"`
	TimestampDayIsAlways24h          *bool `ddl:"keyword" db:"TIMESTAMP_DAY_IS_ALWAYS_24H"`
	TimestampInputFormat             *bool `ddl:"keyword" db:"TIMESTAMP_INPUT_FORMAT"`
	TimestampLtzOutputFormat         *bool `ddl:"keyword" db:"TIMESTAMP_LTZ_OUTPUT_FORMAT"`
	TimestampNtzOutputFormat         *bool `ddl:"keyword" db:"TIMESTAMP_NTZ_OUTPUT_FORMAT"`
	TimestampOutputFormat            *bool `ddl:"keyword" db:"TIMESTAMP_OUTPUT_FORMAT"`
	TimestampTypeMapping             *bool `ddl:"keyword" db:"TIMESTAMP_TYPE_MAPPING"`
	TimestampTzOutputFormat          *bool `ddl:"keyword" db:"TIMESTAMP_TZ_OUTPUT_FORMAT"`
	Timezone                         *bool `ddl:"keyword" db:"TIMEZONE"`
	TimeInputFormat                  *bool `ddl:"keyword" db:"TIME_INPUT_FORMAT"`
	TimeOutputFormat                 *bool `ddl:"keyword" db:"TIME_OUTPUT_FORMAT"`
	TransactionDefaultIsolationLevel *bool `ddl:"keyword" db:"TRANSACTION_DEFAULT_ISOLATION_LEVEL"`
	TwoDigitCenturyStart             *bool `ddl:"keyword" db:"TWO_DIGIT_CENTURY_START"`
	UnsupportedDdlAction             *bool `ddl:"keyword" db:"UNSUPPORTED_DDL_ACTION"`
	UseCachedResult                  *bool `ddl:"keyword" db:"USE_CACHED_RESULT"`
	WeekOfYearPolicy                 *bool `ddl:"keyword" db:"WEEK_OF_YEAR_POLICY"`
	WeekStart                        *bool `ddl:"keyword" db:"WEEK_START"`
}

func (c *accounts) Alter(ctx context.Context, opts *AccountAlterOptions) error {
	if opts == nil {
		opts = &AccountAlterOptions{}
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

type AccountShowOptions struct {
	show     bool  `ddl:"static" db:"SHOW"`                  //lint:ignore U1000 This is used in the ddl tag
	accounts bool  `ddl:"static" db:"ORGANIZATION ACCOUNTS"` //lint:ignore U1000 This is used in the ddl tag
	Like     *Like `ddl:"keyword" db:"LIKE"`
}

func (opts *AccountShowOptions) validate() error {
	return nil
}

type Account struct {
	OrganizationName                     string
	AccountName                          string
	RegionGroup                          string
	SnowflakeRegion                      string
	Edition                              string
	AccountUrl                           string
	CreatedOn                            time.Time
	Comment                              string
	AccountLocator                       string
	AccountLocatorUrl                    string
	ManagedAccounts                      string
	ConsumptionBillingEntityName         string
	MarketplaceConsumerBillingEntityName string
	MarketplaceProviderBillingEntityName string
	OldAccountUrl                        string
	IsOrgAdmin                           bool
}

type accountDBRow struct {
	OrganizationName                     string    `db:"ORGANIZATION_NAME"`
	AccountName                          string    `db:"ACCOUNT_NAME"`
	RegionGroup                          string    `db:"REGION_GROUP"`
	SnowflakeRegion                      string    `db:"SNOWFLAKE_REGION"`
	Edition                              string    `db:"EDITION "`
	AccountUrl                           string    `db:"ACCOUNT_URL"`
	CreatedOn                            time.Time `db:"CREATED_ON"`
	Comment                              string    `db:"COMMENT"`
	AccountLocator                       string    `db:"ACCOUNT_LOCATOR"`
	AccountLocatorUrl                    string    `db:"ACCOUNT_LOCATOR_URL"`
	ManagedAccounts                      string    `db:"MANAGED_ACCOUNTS"`
	ConsumptionBillingEntityName         string    `db:"CONSUMPTION_BILLING_ENTITY_NAME"`
	MarketplaceConsumerBillingEntityName string    `db:"MARKETPLACE_CONSUMER_BILLING_ENTITY_NAME"`
	MarketplaceProviderBillingEntityName string    `db:"MARKETPLACE_PROVIDER_BILLING_ENTITY_NAME"`
	OldAccountUrl                        string    `db:"OLD_ACCOUNT_URL"`
	IsOrgAdmin                           bool      `db:"IS_ORG_ADMIN"`
}

func (row accountDBRow) toAccount() *Account {
	acc := &Account{
		OrganizationName:                     row.OrganizationName,
		AccountName:                          row.AccountName,
		RegionGroup:                          row.RegionGroup,
		SnowflakeRegion:                      row.SnowflakeRegion,
		Edition:                              row.Edition,
		AccountUrl:                           row.AccountUrl,
		CreatedOn:                            row.CreatedOn,
		Comment:                              row.Comment,
		AccountLocator:                       row.AccountLocator,
		AccountLocatorUrl:                    row.AccountLocatorUrl,
		ManagedAccounts:                      row.ManagedAccounts,
		ConsumptionBillingEntityName:         row.ConsumptionBillingEntityName,
		MarketplaceConsumerBillingEntityName: row.MarketplaceConsumerBillingEntityName,
		MarketplaceProviderBillingEntityName: row.MarketplaceProviderBillingEntityName,
		OldAccountUrl:                        row.OldAccountUrl,
		IsOrgAdmin:                           row.IsOrgAdmin,
	}
	return acc
}

func (c *accounts) Show(ctx context.Context, opts *AccountShowOptions) ([]*Account, error) {
	if opts == nil {
		opts = &AccountShowOptions{}
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
	accounts, err := c.Show(ctx, &AccountShowOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.ID().name == id.Name() {
			return account, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

func (v *Account) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.AccountName)
}
