package helpers

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

type TestClient struct {
	context *TestClientContext

	Database *DatabaseClient
}

func NewTestClient(c *sdk.Client, database string, schema string, warehouse string) *TestClient {
	context := &TestClientContext{
		client:    c,
		database:  database,
		schema:    schema,
		warehouse: warehouse,
	}
	return &TestClient{
		context:  context,
		Database: NewDatabaseClient(context),
	}
}

type TestClientContext struct {
	client    *sdk.Client
	database  string
	schema    string
	warehouse string
}
