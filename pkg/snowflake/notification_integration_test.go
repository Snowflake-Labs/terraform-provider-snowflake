package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestNotificationIntegration(t *testing.T) {
	r := require.New(t)
	builder := snowflake.NotificationIntegration("gcp")
	r.NotNil(builder)

	q := builder.Show()
	r.Equal("SHOW NOTIFICATION INTEGRATIONS LIKE 'gcp'", q)

	c := builder.Create()

	c.SetString(`type`, `QUEUE`)
	c.SetString(`gcp_pubsub_subscription_name`, `<subscription_id>`)
	c.SetString(`notification_provider`, `GCP_PUBSUB`)
	c.SetBool(`enabled`, true)
	q = c.Statement()

	r.Equal(`CREATE NOTIFICATION INTEGRATION "gcp" GCP_PUBSUB_SUBSCRIPTION_NAME='<subscription_id>' NOTIFICATION_PROVIDER='GCP_PUBSUB' TYPE='QUEUE' ENABLED=true`, q)
}
