package helpers

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

type TestClient struct {
	context *TestClientContext

	Context      *ContextClient
	Database     *DatabaseClient
	DatabaseRole *DatabaseRoleClient
	Role         *RoleClient
	Schema       *SchemaClient
	User         *UserClient
	Warehouse    *WarehouseClient
}

func NewTestClient(c *sdk.Client, database string, schema string, warehouse string) *TestClient {
	context := &TestClientContext{
		client:    c,
		database:  database,
		schema:    schema,
		warehouse: warehouse,
	}
	return &TestClient{
		context:      context,
		Context:      NewContextClient(context),
		Database:     NewDatabaseClient(context),
		DatabaseRole: NewDatabaseRoleClient(context),
		Role:         NewRoleClient(context),
		Schema:       NewSchemaClient(context),
		User:         NewUserClient(context),
		Warehouse:    NewWarehouseClient(context),
	}
}

type TestClientContext struct {
	client    *sdk.Client
	database  string
	schema    string
	warehouse string
}
