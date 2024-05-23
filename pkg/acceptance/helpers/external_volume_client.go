package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ExternalVolumeClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewExternalVolumeClient(context *TestClientContext, idsGenerator *IdsGenerator) *ExternalVolumeClient {
	return &ExternalVolumeClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ExternalVolumeClient) exec(sql string) error {
	ctx := context.Background()
	_, err := c.context.client.ExecForTests(ctx, sql)
	return err
}

// TODO(SNOW-999142): Use SDK implementation for External Volume once it's available
func (c *ExternalVolumeClient) Create(t *testing.T) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	id := c.ids.RandomAccountObjectIdentifier()
	err := c.exec(fmt.Sprintf(`
create external volume %s
	storage_locations =
    	(
    		(
            	name = 'my-s3-us-west-2'
            	storage_provider = 's3'
            	storage_base_url = 's3://my_example_bucket/'
            	storage_aws_role_arn = 'arn:aws:iam::123456789012:role/myrole'
            	encryption=(type='aws_sse_kms' kms_key_id='1234abcd-12ab-34cd-56ef-1234567890ab')
        	)
      	);
`, id.FullyQualifiedName()))
	require.NoError(t, err)

	return id, c.DropFunc(t, id)
}

func (c *ExternalVolumeClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()

	return func() {
		err := c.exec(fmt.Sprintf(`drop external volume if exists %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
