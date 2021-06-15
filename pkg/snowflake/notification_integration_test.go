package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestNotificationIntegration(t *testing.T) {
	r := require.New(t)
	builder := snowflake.NotificationIntegration("azure")
	r.NotNil(builder)

	q := builder.Show()
	r.Equal("SHOW NOTIFICATION INTEGRATIONS LIKE 'azure'", q)

	c := builder.Create()

	c.SetString(`type`, `QUEUE`)
	c.SetString(`azure_storage_queue_primary_uri`, `azure://my-bucket/my-path/`)
	c.SetString(`azure_tenant_id`, `some-guid`)
	c.SetBool(`enabled`, true)
	q = c.Statement()

	r.Equal(`CREATE NOTIFICATION INTEGRATION "azure" AZURE_STORAGE_QUEUE_PRIMARY_URI='azure://my-bucket/my-path/' AZURE_TENANT_ID='some-guid' TYPE='QUEUE' ENABLED=true`, q)
}
