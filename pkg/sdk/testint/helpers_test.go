package testint

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
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

func createTag(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema) (*sdk.Tag, func()) {
	t.Helper()
	name := random.StringRange(8, 28)
	ctx := context.Background()
	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
	err := client.Tags.Create(ctx, sdk.NewCreateTagRequest(id))
	require.NoError(t, err)
	tag, err := client.Tags.ShowByID(ctx, id)
	require.NoError(t, err)
	return tag, func() {
		err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
		require.NoError(t, err)
	}
}

func createPipe(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, name string, copyStatement string) (*sdk.Pipe, func()) {
	t.Helper()
	require.NotNil(t, database, "database has to be created")
	require.NotNil(t, schema, "schema has to be created")

	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
	ctx := context.Background()

	pipeCleanup := func() {
		err := client.Pipes.Drop(ctx, id)
		require.NoError(t, err)
	}

	err := client.Pipes.Create(ctx, id, copyStatement, &sdk.CreatePipeOptions{})
	if err != nil {
		return nil, pipeCleanup
	}
	require.NoError(t, err)

	createdPipe, errDescribe := client.Pipes.Describe(ctx, id)
	if errDescribe != nil {
		return nil, pipeCleanup
	}
	require.NoError(t, errDescribe)

	return createdPipe, pipeCleanup
}

func createPasswordPolicy(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	return createPasswordPolicyWithOptions(t, client, database, schema, nil)
}

func createPasswordPolicyWithOptions(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, options *sdk.CreatePasswordPolicyOptions) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	var databaseCleanup func()
	if database == nil {
		database, databaseCleanup = testClientHelper().Database.CreateDatabase(t)
	}
	var schemaCleanup func()
	if schema == nil {
		schema, schemaCleanup = testClientHelper().Schema.CreateSchemaInDatabase(t, database.ID())
	}
	name := random.UUID()
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, name)
	ctx := context.Background()
	err := client.PasswordPolicies.Create(ctx, id, options)
	require.NoError(t, err)

	showOptions := &sdk.ShowPasswordPolicyOptions{
		Like: &sdk.Like{
			Pattern: sdk.String(name),
		},
		In: &sdk.In{
			Schema: schema.ID(),
		},
	}
	passwordPolicyList, err := client.PasswordPolicies.Show(ctx, showOptions)
	require.NoError(t, err)
	require.Equal(t, 1, len(passwordPolicyList))
	return &passwordPolicyList[0], func() {
		err := client.PasswordPolicies.Drop(ctx, id, nil)
		require.NoError(t, err)
		if schemaCleanup != nil {
			schemaCleanup()
		}
		if databaseCleanup != nil {
			databaseCleanup()
		}
	}
}

func createNetworkPolicy(t *testing.T, client *sdk.Client, req *sdk.CreateNetworkPolicyRequest) (error, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.NetworkPolicies.Create(ctx, req)
	return err, func() {
		err := client.NetworkPolicies.Drop(ctx, sdk.NewDropNetworkPolicyRequest(req.GetName()))
		require.NoError(t, err)
	}
}

func createResourceMonitor(t *testing.T, client *sdk.Client) (*sdk.ResourceMonitor, func()) {
	t.Helper()
	return createResourceMonitorWithOptions(t, client, &sdk.CreateResourceMonitorOptions{
		With: &sdk.ResourceMonitorWith{
			CreditQuota: sdk.Pointer(100),
			Triggers: []sdk.TriggerDefinition{
				{
					Threshold:     100,
					TriggerAction: sdk.TriggerActionSuspend,
				},
				{
					Threshold:     70,
					TriggerAction: sdk.TriggerActionSuspendImmediate,
				},
				{
					Threshold:     90,
					TriggerAction: sdk.TriggerActionNotify,
				},
			},
		},
	})
}

func createResourceMonitorWithOptions(t *testing.T, client *sdk.Client, opts *sdk.CreateResourceMonitorOptions) (*sdk.ResourceMonitor, func()) {
	t.Helper()
	id := sdk.RandomAccountObjectIdentifier()
	ctx := context.Background()
	err := client.ResourceMonitors.Create(ctx, id, opts)
	require.NoError(t, err)
	resourceMonitor, err := client.ResourceMonitors.ShowByID(ctx, id)
	require.NoError(t, err)
	return resourceMonitor, func() {
		err := client.ResourceMonitors.Drop(ctx, id)
		require.NoError(t, err)
	}
}

func createMaskingPolicy(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	signature := []sdk.TableColumnSignature{
		{
			Name: random.String(),
			Type: sdk.DataTypeVARCHAR,
		},
	}
	n := random.IntRange(0, 5)
	for i := 0; i < n; i++ {
		signature = append(signature, sdk.TableColumnSignature{
			Name: random.String(),
			Type: sdk.DataTypeVARCHAR,
		})
	}
	expression := "REPLACE('X', 1, 2)"
	return createMaskingPolicyWithOptions(t, client, database, schema, signature, sdk.DataTypeVARCHAR, expression, &sdk.CreateMaskingPolicyOptions{})
}

func createMaskingPolicyIdentity(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, columnType sdk.DataType) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	name := "a"
	signature := []sdk.TableColumnSignature{
		{
			Name: name,
			Type: columnType,
		},
	}
	expression := "a"
	return createMaskingPolicyWithOptions(t, client, database, schema, signature, columnType, expression, &sdk.CreateMaskingPolicyOptions{})
}

func createMaskingPolicyWithOptions(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, signature []sdk.TableColumnSignature, returns sdk.DataType, expression string, options *sdk.CreateMaskingPolicyOptions) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	var databaseCleanup func()
	if database == nil {
		database, databaseCleanup = testClientHelper().Database.CreateDatabase(t)
	}
	var schemaCleanup func()
	if schema == nil {
		schema, schemaCleanup = testClientHelper().Schema.CreateSchemaInDatabase(t, database.ID())
	}
	name := random.String()
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, name)
	ctx := context.Background()
	err := client.MaskingPolicies.Create(ctx, id, signature, returns, expression, options)
	require.NoError(t, err)

	showOptions := &sdk.ShowMaskingPolicyOptions{
		Like: &sdk.Like{
			Pattern: sdk.String(name),
		},
		In: &sdk.In{
			Schema: schema.ID(),
		},
	}
	maskingPolicyList, err := client.MaskingPolicies.Show(ctx, showOptions)
	require.NoError(t, err)
	require.Equal(t, 1, len(maskingPolicyList))
	return &maskingPolicyList[0], func() {
		err := client.MaskingPolicies.Drop(ctx, id)
		require.NoError(t, err)
		if schemaCleanup != nil {
			schemaCleanup()
		}
		if databaseCleanup != nil {
			databaseCleanup()
		}
	}
}

func createAlert(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, warehouse *sdk.Warehouse) (*sdk.Alert, func()) {
	t.Helper()
	schedule := "USING CRON * * * * * UTC"
	condition := "SELECT 1"
	action := "SELECT 1"
	return createAlertWithOptions(t, client, database, schema, warehouse, schedule, condition, action, &sdk.CreateAlertOptions{})
}

func createAlertWithOptions(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, warehouse *sdk.Warehouse, schedule string, condition string, action string, opts *sdk.CreateAlertOptions) (*sdk.Alert, func()) {
	t.Helper()
	var databaseCleanup func()
	if database == nil {
		database, databaseCleanup = testClientHelper().Database.CreateDatabase(t)
	}
	var schemaCleanup func()
	if schema == nil {
		schema, schemaCleanup = testClientHelper().Schema.CreateSchemaInDatabase(t, database.ID())
	}
	var warehouseCleanup func()
	if warehouse == nil {
		warehouse, warehouseCleanup = testClientHelper().Warehouse.CreateWarehouse(t)
	}

	name := random.String()
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, name)
	ctx := context.Background()
	err := client.Alerts.Create(ctx, id, warehouse.ID(), schedule, condition, action, opts)
	require.NoError(t, err)

	showOptions := &sdk.ShowAlertOptions{
		Like: &sdk.Like{
			Pattern: sdk.String(name),
		},
		In: &sdk.In{
			Schema: schema.ID(),
		},
	}
	alertList, err := client.Alerts.Show(ctx, showOptions)
	require.NoError(t, err)
	require.Equal(t, 1, len(alertList))
	return &alertList[0], func() {
		err := client.Alerts.Drop(ctx, id)
		require.NoError(t, err)
		if schemaCleanup != nil {
			schemaCleanup()
		}
		if databaseCleanup != nil {
			databaseCleanup()
		}
		if warehouseCleanup != nil {
			warehouseCleanup()
		}
	}
}

func createFailoverGroup(t *testing.T, client *sdk.Client) (*sdk.FailoverGroup, func()) {
	t.Helper()
	objectTypes := []sdk.PluralObjectType{sdk.PluralObjectTypeRoles}
	currentAccount := testClientHelper().Context.CurrentAccount(t)
	accountID := sdk.NewAccountIdentifierFromAccountLocator(currentAccount)
	allowedAccounts := []sdk.AccountIdentifier{accountID}
	return createFailoverGroupWithOptions(t, client, objectTypes, allowedAccounts, nil)
}

func createFailoverGroupWithOptions(t *testing.T, client *sdk.Client, objectTypes []sdk.PluralObjectType, allowedAccounts []sdk.AccountIdentifier, opts *sdk.CreateFailoverGroupOptions) (*sdk.FailoverGroup, func()) {
	t.Helper()
	id := sdk.RandomAlphanumericAccountObjectIdentifier()
	ctx := context.Background()
	err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, opts)
	require.NoError(t, err)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	require.NoError(t, err)
	return failoverGroup, func() {
		err := client.FailoverGroups.Drop(ctx, id, nil)
		require.NoError(t, err)
	}
}

func createShare(t *testing.T, client *sdk.Client) (*sdk.Share, func()) {
	t.Helper()
	// TODO(SNOW-1058419): Try with identifier containing dot during identifiers rework
	id := sdk.RandomAlphanumericAccountObjectIdentifier()
	return createShareWithOptions(t, client, id, &sdk.CreateShareOptions{})
}

func createShareWithOptions(t *testing.T, client *sdk.Client, id sdk.AccountObjectIdentifier, opts *sdk.CreateShareOptions) (*sdk.Share, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.Shares.Create(ctx, id, opts)
	require.NoError(t, err)
	share, err := client.Shares.ShowByID(ctx, id)
	require.NoError(t, err)
	return share, func() {
		err := client.Shares.Drop(ctx, id)
		require.NoError(t, err)
	}
}

func createFileFormat(t *testing.T, client *sdk.Client, schema sdk.DatabaseObjectIdentifier) (*sdk.FileFormat, func()) {
	t.Helper()
	return createFileFormatWithOptions(t, client, schema, &sdk.CreateFileFormatOptions{
		Type: sdk.FileFormatTypeCSV,
	})
}

func createFileFormatWithOptions(t *testing.T, client *sdk.Client, schema sdk.DatabaseObjectIdentifier, opts *sdk.CreateFileFormatOptions) (*sdk.FileFormat, func()) {
	t.Helper()
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName(), schema.Name(), random.String())
	ctx := context.Background()
	err := client.FileFormats.Create(ctx, id, opts)
	require.NoError(t, err)
	fileFormat, err := client.FileFormats.ShowByID(ctx, id)
	require.NoError(t, err)
	return fileFormat, func() {
		err := client.FileFormats.Drop(ctx, id, nil)
		require.NoError(t, err)
	}
}

func createView(t *testing.T, client *sdk.Client, viewId sdk.SchemaObjectIdentifier, asQuery string) func() {
	t.Helper()
	ctx := context.Background()
	_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE VIEW %s AS %s`, viewId.FullyQualifiedName(), asQuery))
	require.NoError(t, err)

	return func() {
		_, err := client.ExecForTests(ctx, fmt.Sprintf(`DROP VIEW %s`, viewId.FullyQualifiedName()))
		require.NoError(t, err)
	}
}

func createRowAccessPolicy(t *testing.T, client *sdk.Client, schema *sdk.Schema) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, random.String())

	arg := sdk.NewCreateRowAccessPolicyArgsRequest("A", sdk.DataTypeNumber)
	body := "true"
	createRequest := sdk.NewCreateRowAccessPolicyRequest(id, []sdk.CreateRowAccessPolicyArgsRequest{*arg}, body)
	err := client.RowAccessPolicies.Create(ctx, createRequest)
	require.NoError(t, err)

	return id, func() {
		err := client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id))
		require.NoError(t, err)
	}
}

func createApiIntegration(t *testing.T, client *sdk.Client) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	id := sdk.NewAccountObjectIdentifier(random.String())
	apiAllowedPrefixes := []sdk.ApiIntegrationEndpointPrefix{{Path: "https://xyz.execute-api.us-west-2.amazonaws.com/production"}}
	req := sdk.NewCreateApiIntegrationRequest(id, apiAllowedPrefixes, true)
	req.WithAwsApiProviderParams(sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiGateway, "arn:aws:iam::123456789012:role/hello_cloud_account_role"))
	err := client.ApiIntegrations.Create(ctx, req)
	require.NoError(t, err)

	return id, func() {
		err := client.ApiIntegrations.Drop(ctx, sdk.NewDropApiIntegrationRequest(id))
		require.NoError(t, err)
	}
}

func createExternalFunction(t *testing.T, client *sdk.Client, schema *sdk.Schema) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	apiIntegration, cleanupApiIntegration := createApiIntegration(t, client)
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, random.StringN(4))
	argument := sdk.NewExternalFunctionArgumentRequest("x", sdk.DataTypeVARCHAR)
	argumentsRequest := []sdk.ExternalFunctionArgumentRequest{*argument}
	as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
	request := sdk.NewCreateExternalFunctionRequest(id, sdk.DataTypeVariant, &apiIntegration, as).
		WithOrReplace(sdk.Bool(true)).
		WithArguments(argumentsRequest)
	err := client.ExternalFunctions.Create(ctx, request)
	require.NoError(t, err)
	return id, func() {
		cleanupApiIntegration()
		err = client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id, []sdk.DataType{sdk.DataTypeVARCHAR}))
		require.NoError(t, err)
	}
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

func updateAccountParameterTemporarily(t *testing.T, client *sdk.Client, parameter sdk.AccountParameter, newValue string) func() {
	t.Helper()
	ctx := context.Background()

	param, err := client.Parameters.ShowAccountParameter(ctx, parameter)
	require.NoError(t, err)
	oldValue := param.Value

	err = client.Parameters.SetAccountParameter(ctx, parameter, newValue)
	require.NoError(t, err)

	return func() {
		err = client.Parameters.SetAccountParameter(ctx, parameter, oldValue)
		require.NoError(t, err)
	}
}

func createTaskWithRequest(t *testing.T, client *sdk.Client, request *sdk.CreateTaskRequest) (*sdk.Task, func()) {
	t.Helper()
	ctx := context.Background()

	id := request.GetName()

	err := client.Tasks.Create(ctx, request)
	require.NoError(t, err)

	task, err := client.Tasks.ShowByID(ctx, id)
	require.NoError(t, err)

	return task, func() {
		err = client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(id))
		require.NoError(t, err)
	}
}

func createTask(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema) (*sdk.Task, func()) {
	t.Helper()
	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, random.AlphaN(20))
	warehouseReq := sdk.NewCreateTaskWarehouseRequest().WithWarehouse(sdk.Pointer(testWarehouse(t).ID()))
	return createTaskWithRequest(t, client, sdk.NewCreateTaskRequest(id, "SELECT CURRENT_TIMESTAMP").WithSchedule(sdk.String("60 minutes")).WithWarehouse(warehouseReq))
}
