package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ProcedureClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewProcedureClient(context *TestClientContext, idsGenerator *IdsGenerator) *ProcedureClient {
	return &ProcedureClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ProcedureClient) client() sdk.Procedures {
	return c.context.client.Procedures
}

func (c *ProcedureClient) Create(t *testing.T, arguments ...sdk.DataType) *sdk.Procedure {
	t.Helper()
	return c.CreateWithIdentifier(t, c.ids.RandomSchemaObjectIdentifierWithArguments(arguments...))
}

func (c *ProcedureClient) CreateWithIdentifier(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) *sdk.Procedure {
	t.Helper()
	ctx := context.Background()
	argumentRequests := make([]sdk.ProcedureArgumentRequest, len(id.ArgumentDataTypes()))
	for i, argumentDataType := range id.ArgumentDataTypes() {
		argumentRequests[i] = *sdk.NewProcedureArgumentRequest(c.ids.Alpha(), nil).WithArgDataTypeOld(argumentDataType)
	}
	err := c.client().CreateForSQL(ctx,
		sdk.NewCreateForSQLProcedureRequestDefinitionWrapped(
			id.SchemaObjectId(),
			*sdk.NewProcedureSQLReturnsRequest().WithResultDataType(*sdk.NewProcedureReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeInt)),
			`BEGIN RETURN 1; END`).WithArguments(argumentRequests),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, c.context.client.Procedures.Drop(ctx, sdk.NewDropProcedureRequest(id).WithIfExists(true)))
	})

	procedure, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return procedure
}

func (c *ProcedureClient) DropProcedureFunc(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropProcedureRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *ProcedureClient) Show(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.Procedure, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *ProcedureClient) DescribeDetails(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.ProcedureDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().DescribeDetails(ctx, id)
}

// Session argument is needed: https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java#data-access-example
// More references: https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java
func (c *ProcedureClient) SampleJavaDefinition(t *testing.T, className string, funcName string, argName string) string {
	t.Helper()

	return fmt.Sprintf(`
	import com.snowflake.snowpark_java.*;
	class %[1]s {
		public static String %[2]s(Session session, String %[3]s) {
			return %[3]s;
		}
	}
`, className, funcName, argName)
}

// For more references: https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-javascript
func (c *ProcedureClient) SampleJavascriptDefinition(t *testing.T, argName string) string {
	t.Helper()

	return fmt.Sprintf(`
	if (%[1]s <= 0) {
		return 1;
	} else {
		var result = 1;
		for (var i = 2; i <= %[1]s; i++) {
			result = result * i;
		}
		return result;
	}
`, argName)
}
