package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ExternalTableClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewExternalTableClient(context *TestClientContext, idsGenerator *IdsGenerator) *ExternalTableClient {
	return &ExternalTableClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ExternalTableClient) client() sdk.ExternalTables {
	return c.context.client.ExternalTables
}

func (c *ExternalTableClient) PublishDataToStage(t *testing.T, stageId sdk.SchemaObjectIdentifier, data []byte) {
	ctx := context.Background()

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf(`copy into @%s/external_tables_test_data/test_data from (select parse_json('%s')) overwrite = true`, stageId.FullyQualifiedName(), string(data)))
	require.NoError(t, err)
}
