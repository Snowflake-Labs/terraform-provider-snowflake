package testint

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

const (
	nycWeatherDataURL = "s3://snowflake-workshop-lab/weather-nyc"
)

// there is no direct way to get the account identifier from Snowflake API, but you can get it if you know
// the account locator and by filtering the list of accounts in replication accounts by the account locator
func getAccountIdentifier(t *testing.T, client *sdk.Client) sdk.AccountIdentifier {
	t.Helper()
	ctx := context.Background()
	// TODO: replace later (incoming clients differ)
	currentAccountLocator, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	replicationAccounts, err := client.ReplicationFunctions.ShowReplicationAccounts(ctx)
	require.NoError(t, err)
	for _, replicationAccount := range replicationAccounts {
		if replicationAccount.AccountLocator == currentAccountLocator {
			return sdk.NewAccountIdentifier(replicationAccount.OrganizationName, replicationAccount.AccountName)
		}
	}
	return sdk.AccountIdentifier{}
}

// TODO: extract getting row access policies as resource (like getting tag in system functions)
// getRowAccessPolicyFor is based on https://docs.snowflake.com/en/user-guide/security-row-intro#obtain-database-objects-with-a-row-access-policy.
func getRowAccessPolicyFor(t *testing.T, client *sdk.Client, id sdk.SchemaObjectIdentifier, objectType sdk.ObjectType) (*policyReference, error) {
	t.Helper()
	ctx := context.Background()

	s := &policyReference{}
	policyReferencesId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), "INFORMATION_SCHEMA", "POLICY_REFERENCES")
	err := client.QueryOneForTests(ctx, s, fmt.Sprintf(`SELECT * FROM TABLE(%s(REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => '%v'))`, policyReferencesId.FullyQualifiedName(), id.FullyQualifiedName(), objectType))

	return s, err
}

type policyReference struct {
	PolicyDb          string         `db:"POLICY_DB"`
	PolicySchema      string         `db:"POLICY_SCHEMA"`
	PolicyName        string         `db:"POLICY_NAME"`
	PolicyKind        string         `db:"POLICY_KIND"`
	RefDatabaseName   string         `db:"REF_DATABASE_NAME"`
	RefSchemaName     string         `db:"REF_SCHEMA_NAME"`
	RefEntityName     string         `db:"REF_ENTITY_NAME"`
	RefEntityDomain   string         `db:"REF_ENTITY_DOMAIN"`
	RefColumnName     sql.NullString `db:"REF_COLUMN_NAME"`
	RefArgColumnNames string         `db:"REF_ARG_COLUMN_NAMES"`
	TagDatabase       sql.NullString `db:"TAG_DATABASE"`
	TagSchema         sql.NullString `db:"TAG_SCHEMA"`
	TagName           sql.NullString `db:"TAG_NAME"`
	PolicyStatus      string         `db:"POLICY_STATUS"`
}

// TODO: extract getting table columns as resource (like getting tag in system functions)
// getTableColumnsFor is based on https://docs.snowflake.com/en/sql-reference/info-schema/columns.
func getTableColumnsFor(t *testing.T, client *sdk.Client, tableId sdk.SchemaObjectIdentifier) []informationSchemaColumns {
	t.Helper()
	ctx := context.Background()

	var columns []informationSchemaColumns
	query := fmt.Sprintf("SELECT * FROM information_schema.columns WHERE table_schema = '%s'  AND table_name = '%s' ORDER BY ordinal_position", tableId.SchemaName(), tableId.Name())
	err := client.QueryForTests(ctx, &columns, query)
	require.NoError(t, err)

	return columns
}

type informationSchemaColumns struct {
	TableCatalog           string         `db:"TABLE_CATALOG"`
	TableSchema            string         `db:"TABLE_SCHEMA"`
	TableName              string         `db:"TABLE_NAME"`
	ColumnName             string         `db:"COLUMN_NAME"`
	OrdinalPosition        string         `db:"ORDINAL_POSITION"`
	ColumnDefault          sql.NullString `db:"COLUMN_DEFAULT"`
	IsNullable             string         `db:"IS_NULLABLE"`
	DataType               string         `db:"DATA_TYPE"`
	CharacterMaximumLength sql.NullString `db:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   sql.NullString `db:"CHARACTER_OCTET_LENGTH"`
	NumericPrecision       sql.NullString `db:"NUMERIC_PRECISION"`
	NumericPrecisionRadix  sql.NullString `db:"NUMERIC_PRECISION_RADIX"`
	NumericScale           sql.NullString `db:"NUMERIC_SCALE"`
	DatetimePrecision      sql.NullString `db:"DATETIME_PRECISION"`
	IntervalType           sql.NullString `db:"INTERVAL_TYPE"`
	IntervalPrecision      sql.NullString `db:"INTERVAL_PRECISION"`
	CharacterSetCatalog    sql.NullString `db:"CHARACTER_SET_CATALOG"`
	CharacterSetSchema     sql.NullString `db:"CHARACTER_SET_SCHEMA"`
	CharacterSetName       sql.NullString `db:"CHARACTER_SET_NAME"`
	CollationCatalog       sql.NullString `db:"COLLATION_CATALOG"`
	CollationSchema        sql.NullString `db:"COLLATION_SCHEMA"`
	CollationName          sql.NullString `db:"COLLATION_NAME"`
	DomainCatalog          sql.NullString `db:"DOMAIN_CATALOG"`
	DomainSchema           sql.NullString `db:"DOMAIN_SCHEMA"`
	DomainName             sql.NullString `db:"DOMAIN_NAME"`
	UdtCatalog             sql.NullString `db:"UDT_CATALOG"`
	UdtSchema              sql.NullString `db:"UDT_SCHEMA"`
	UdtName                sql.NullString `db:"UDT_NAME"`
	ScopeCatalog           sql.NullString `db:"SCOPE_CATALOG"`
	ScopeSchema            sql.NullString `db:"SCOPE_SCHEMA"`
	ScopeName              sql.NullString `db:"SCOPE_NAME"`
	MaximumCardinality     sql.NullString `db:"MAXIMUM_CARDINALITY"`
	DtdIdentifier          sql.NullString `db:"DTD_IDENTIFIER"`
	IsSelfReferencing      string         `db:"IS_SELF_REFERENCING"`
	IsIdentity             string         `db:"IS_IDENTITY"`
	IdentityGeneration     sql.NullString `db:"IDENTITY_GENERATION"`
	IdentityStart          sql.NullString `db:"IDENTITY_START"`
	IdentityIncrement      sql.NullString `db:"IDENTITY_INCREMENT"`
	IdentityMaximum        sql.NullString `db:"IDENTITY_MAXIMUM"`
	IdentityMinimum        sql.NullString `db:"IDENTITY_MINIMUM"`
	IdentityCycle          sql.NullString `db:"IDENTITY_CYCLE"`
	IdentityOrdered        sql.NullString `db:"IDENTITY_ORDERED"`
	Comment                sql.NullString `db:"COMMENT"`
}
