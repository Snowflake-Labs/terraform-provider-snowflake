package helpers

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

type TestClient struct {
	context *TestClientContext

	Account            *AccountClient
	Alert              *AlertClient
	ApiIntegration     *ApiIntegrationClient
	Application        *ApplicationClient
	ApplicationPackage *ApplicationPackageClient
	Context            *ContextClient
	Database           *DatabaseClient
	DatabaseRole       *DatabaseRoleClient
	DynamicTable       *DynamicTableClient
	FailoverGroup      *FailoverGroupClient
	FileFormat         *FileFormatClient
	MaskingPolicy      *MaskingPolicyClient
	NetworkPolicy      *NetworkPolicyClient
	Parameter          *ParameterClient
	PasswordPolicy     *PasswordPolicyClient
	Pipe               *PipeClient
	ResourceMonitor    *ResourceMonitorClient
	Role               *RoleClient
	RowAccessPolicy    *RowAccessPolicyClient
	Schema             *SchemaClient
	SessionPolicy      *SessionPolicyClient
	Share              *ShareClient
	Stage              *StageClient
	Table              *TableClient
	Tag                *TagClient
	Task               *TaskClient
	User               *UserClient
	View               *ViewClient
	Warehouse          *WarehouseClient
}

func NewTestClient(c *sdk.Client, database string, schema string, warehouse string) *TestClient {
	context := &TestClientContext{
		client:    c,
		database:  database,
		schema:    schema,
		warehouse: warehouse,
	}
	return &TestClient{
		context:            context,
		Account:            NewAccountClient(context),
		Alert:              NewAlertClient(context),
		ApiIntegration:     NewApiIntegrationClient(context),
		Application:        NewApplicationClient(context),
		ApplicationPackage: NewApplicationPackageClient(context),
		Context:            NewContextClient(context),
		Database:           NewDatabaseClient(context),
		DatabaseRole:       NewDatabaseRoleClient(context),
		DynamicTable:       NewDynamicTableClient(context),
		FailoverGroup:      NewFailoverGroupClient(context),
		FileFormat:         NewFileFormatClient(context),
		MaskingPolicy:      NewMaskingPolicyClient(context),
		NetworkPolicy:      NewNetworkPolicyClient(context),
		Parameter:          NewParameterClient(context),
		PasswordPolicy:     NewPasswordPolicyClient(context),
		Pipe:               NewPipeClient(context),
		ResourceMonitor:    NewResourceMonitorClient(context),
		Role:               NewRoleClient(context),
		RowAccessPolicy:    NewRowAccessPolicyClient(context),
		Schema:             NewSchemaClient(context),
		SessionPolicy:      NewSessionPolicyClient(context),
		Share:              NewShareClient(context),
		Stage:              NewStageClient(context),
		Table:              NewTableClient(context),
		Tag:                NewTagClient(context),
		Task:               NewTaskClient(context),
		User:               NewUserClient(context),
		View:               NewViewClient(context),
		Warehouse:          NewWarehouseClient(context),
	}
}

type TestClientContext struct {
	client    *sdk.Client
	database  string
	schema    string
	warehouse string
}

func (c *TestClientContext) databaseId() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.database)
}

func (c *TestClientContext) schemaId() sdk.DatabaseObjectIdentifier {
	return sdk.NewDatabaseObjectIdentifier(c.database, c.schema)
}

func (c *TestClientContext) warehouseId() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.warehouse)
}

func (c *TestClientContext) newSchemaObjectIdentifier(name string) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifier(c.database, c.schema, name)
}
