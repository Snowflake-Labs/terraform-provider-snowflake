package helpers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type TableClient struct {
	context *TestClientContext
}

func NewTableClient(context *TestClientContext) *TableClient {
	return &TableClient{
		context: context,
	}
}

func (c *TableClient) client() sdk.Tables {
	return c.context.client.Tables
}

func (c *TableClient) CreateTable(t *testing.T) (*sdk.Table, func()) {
	t.Helper()
	return c.CreateTableInSchema(t, c.context.schemaId())
}

func (c *TableClient) CreateTableInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.Table, func()) {
	t.Helper()

	columns := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
	}
	name := random.StringRange(8, 28)
	return c.CreateTableWithColumns(t, schemaId, name, columns)
}

func (c *TableClient) CreateTableWithColumns(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, name string, columns []sdk.TableColumnRequest) (*sdk.Table, func()) {
	t.Helper()

	id := sdk.NewSchemaObjectIdentifier(schemaId.DatabaseName(), schemaId.Name(), name)
	ctx := context.Background()

	dbCreateRequest := sdk.NewCreateTableRequest(id, columns)
	err := c.client().Create(ctx, dbCreateRequest)
	require.NoError(t, err)

	table, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return table, c.DropTableFunc(t, id)
}

func (c *TableClient) DropTableFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		// to prevent error when schema was removed before the table
		_, err := c.context.client.Schemas.ShowByID(ctx, sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName()))
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			return
		}

		dropErr := c.client().Drop(ctx, sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, dropErr)
	}
}

// GetTableColumnsFor is based on https://docs.snowflake.com/en/sql-reference/info-schema/columns.
// TODO: extract getting table columns as resource (like getting tag in system functions)
func (c *TableClient) GetTableColumnsFor(t *testing.T, tableId sdk.SchemaObjectIdentifier) []InformationSchemaColumns {
	t.Helper()
	ctx := context.Background()

	var columns []InformationSchemaColumns
	query := fmt.Sprintf("SELECT * FROM information_schema.columns WHERE table_schema = '%s' AND table_name = '%s' ORDER BY ordinal_position", tableId.SchemaName(), tableId.Name())
	err := c.context.client.QueryForTests(ctx, &columns, query)
	require.NoError(t, err)

	return columns
}

type InformationSchemaColumns struct {
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
