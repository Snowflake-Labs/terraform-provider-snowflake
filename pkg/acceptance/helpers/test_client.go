package helpers

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

type TestClient struct {
	context *TestClientContext

	Alert              *AlertClient
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
	PasswordPolicy     *PasswordPolicyClient
	Pipe               *PipeClient
	ResourceMonitor    *ResourceMonitorClient
	Role               *RoleClient
	Schema             *SchemaClient
	SessionPolicy      *SessionPolicyClient
	Stage              *StageClient
	Table              *TableClient
	Tag                *TagClient
	User               *UserClient
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
		Alert:              NewAlertClient(context),
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
		PasswordPolicy:     NewPasswordPolicyClient(context),
		Pipe:               NewPipeClient(context),
		ResourceMonitor:    NewResourceMonitorClient(context),
		Role:               NewRoleClient(context),
		Schema:             NewSchemaClient(context),
		SessionPolicy:      NewSessionPolicyClient(context),
		Stage:              NewStageClient(context),
		Tag:                NewTagClient(context),
		Table:              NewTableClient(context),
		User:               NewUserClient(context),
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
