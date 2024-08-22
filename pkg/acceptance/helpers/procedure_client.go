package helpers

import (
	"context"
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
		argumentRequests[i] = *sdk.NewProcedureArgumentRequest(c.ids.Alpha(), argumentDataType)
	}
	err := c.client().CreateForSQL(ctx,
		sdk.NewCreateForSQLProcedureRequest(
			id.SchemaObjectId(),
			*sdk.NewProcedureSQLReturnsRequest().WithResultDataType(*sdk.NewProcedureReturnsResultDataTypeRequest(sdk.DataTypeInt)),
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
