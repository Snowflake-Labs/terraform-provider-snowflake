package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type DataMetricFunctionReferencesClient struct {
	context *TestClientContext
}

func NewDataMetricFunctionReferencesClient(context *TestClientContext) *DataMetricFunctionReferencesClient {
	return &DataMetricFunctionReferencesClient{
		context: context,
	}
}

// GetDataMetricFunctionReferences is based on https://docs.snowflake.com/en/sql-reference/functions/data_metric_function_references.
func (c *DataMetricFunctionReferencesClient) GetDataMetricFunctionReferences(t *testing.T, id sdk.SchemaObjectIdentifier, objectType sdk.ObjectType) ([]DataMetricFunctionReference, error) {
	t.Helper()
	ctx := context.Background()

	s := []DataMetricFunctionReference{}
	dmfReferencesId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), "INFORMATION_SCHEMA", "DATA_METRIC_FUNCTION_REFERENCES")
	err := c.context.client.QueryForTests(ctx, &s, fmt.Sprintf(`SELECT * FROM TABLE(%s(REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => '%v'))`, dmfReferencesId.FullyQualifiedName(), id.FullyQualifiedName(), objectType))

	return s, err
}

type DataMetricFunctionReference struct {
	MetricDatabaseName    string `db:"METRIC_DATABASE_NAME"`
	MetricSchemaName      string `db:"METRIC_SCHEMA_NAME"`
	MetricName            string `db:"METRIC_NAME"`
	MetricSignature       string `db:"METRIC_SIGNATURE"`
	MetricDataType        string `db:"METRIC_DATA_TYPE"`
	RefEntityDatabaseName string `db:"REF_ENTITY_DATABASE_NAME"`
	RefEntitySchemaName   string `db:"REF_ENTITY_SCHEMA_NAME"`
	RefEntityName         string `db:"REF_ENTITY_NAME"`
	RefEntityDomain       string `db:"REF_ENTITY_DOMAIN"`
	RefArguments          string `db:"REF_ARGUMENTS"`
	RefId                 string `db:"REF_ID"`
	Schedule              string `db:"SCHEDULE"`
	ScheduleStatus        string `db:"SCHEDULE_STATUS"`
}
