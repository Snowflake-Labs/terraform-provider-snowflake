package helpers

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

type TestClient struct {
	context *TestClientContext

	Application        *ApplicationClient
	ApplicationPackage *ApplicationPackageClient
	Context            *ContextClient
	Database           *DatabaseClient
	DatabaseRole       *DatabaseRoleClient
	DynamicTable       *DynamicTableClient
	Pipe               *PipeClient
	Role               *RoleClient
	Schema             *SchemaClient
	SessionPolicy      *SessionPolicyClient
	Stage              *StageClient
	Table              *TableClient
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
		Application:        NewApplicationClient(context),
		ApplicationPackage: NewApplicationPackageClient(context),
		Context:            NewContextClient(context),
		Database:           NewDatabaseClient(context),
		DatabaseRole:       NewDatabaseRoleClient(context),
		DynamicTable:       NewDynamicTableClient(context),
		Pipe:               NewPipeClient(context),
		Role:               NewRoleClient(context),
		Schema:             NewSchemaClient(context),
		SessionPolicy:      NewSessionPolicyClient(context),
		Stage:              NewStageClient(context),
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
