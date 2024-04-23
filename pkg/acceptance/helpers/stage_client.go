package helpers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type StageClient struct {
	context *TestClientContext
}

func NewStageClient(context *TestClientContext) *StageClient {
	return &StageClient{
		context: context,
	}
}

func (c *StageClient) client() sdk.Stages {
	return c.context.client.Stages
}

// TODO: use default schema
func (c *StageClient) CreateStageWithDirectory(t *testing.T, schema *sdk.Schema, name string) (*sdk.Stage, func()) {
	t.Helper()
	id := sdk.NewSchemaObjectIdentifier(c.context.database, schema.Name, name)
	return c.CreateStageWithOptions(t, id, func(request *sdk.CreateInternalStageRequest) *sdk.CreateInternalStageRequest {
		return request.WithDirectoryTableOptions(sdk.NewInternalDirectoryTableOptionsRequest().WithEnable(sdk.Bool(true)))
	})
}

func (c *StageClient) CreateStageWithURL(t *testing.T, id sdk.SchemaObjectIdentifier, url string) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()
	err := c.client().CreateOnS3(ctx, sdk.NewCreateOnS3StageRequest(id).
		WithExternalStageParams(sdk.NewExternalS3StageParamsRequest(url)))
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

func (c *StageClient) CreateStage(t *testing.T) (*sdk.Stage, func()) {
	t.Helper()
	return c.CreateStageInSchema(t, sdk.NewDatabaseObjectIdentifier(c.context.database, c.context.schema))
}

func (c *StageClient) CreateStageInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.Stage, func()) {
	t.Helper()
	return c.CreateStageWithOptions(t, sdk.NewSchemaObjectIdentifier(schemaId.DatabaseName(), schemaId.Name(), random.AlphaN(8)), func(request *sdk.CreateInternalStageRequest) *sdk.CreateInternalStageRequest { return request })
}

func (c *StageClient) CreateStageWithOptions(t *testing.T, id sdk.SchemaObjectIdentifier, reqMapping func(*sdk.CreateInternalStageRequest) *sdk.CreateInternalStageRequest) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()
	err := c.client().CreateInternal(ctx, reqMapping(sdk.NewCreateInternalStageRequest(id)))
	require.NoError(t, err)

	stage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, c.DropStageFunc(t, id)
}

func (c *StageClient) DropStageFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropStageRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}

func (c *StageClient) PutOnStage(t *testing.T, id sdk.SchemaObjectIdentifier, filename string) {
	t.Helper()
	ctx := context.Background()

	path, err := filepath.Abs("./testdata/" + filename)
	require.NoError(t, err)
	absPath := "file://" + path

	_, err = c.context.client.ExecForTests(ctx, fmt.Sprintf(`PUT '%s' @%s AUTO_COMPRESS = FALSE`, absPath, id.FullyQualifiedName()))
	require.NoError(t, err)
}

func (c *StageClient) PutOnStageWithContent(t *testing.T, id sdk.SchemaObjectIdentifier, filename string, content string) {
	t.Helper()
	ctx := context.Background()

	tf := fmt.Sprintf("/tmp/%s", filename)
	f, err := os.Create(tf)
	require.NoError(t, err)
	if content != "" {
		_, err = f.Write([]byte(content))
		require.NoError(t, err)
	}
	f.Close()
	defer os.Remove(f.Name())

	_, err = c.context.client.ExecForTests(ctx, fmt.Sprintf(`PUT file://%s @%s AUTO_COMPRESS = FALSE OVERWRITE = TRUE`, f.Name(), id.FullyQualifiedName()))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err = c.context.client.ExecForTests(ctx, fmt.Sprintf(`REMOVE @%s/%s`, id.FullyQualifiedName(), filename))
		require.NoError(t, err)
	})
}
