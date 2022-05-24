package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestExternalOauthIntegration(t *testing.T) {
	r := require.New(t)
	builder := snowflake.ExternalOauthIntegration("azure")
	r.NotNil(builder)

	q := builder.Show()
	r.Equal("SHOW SECURITY INTEGRATIONS LIKE 'azure'", q)

	q = builder.Describe()
	r.Equal("DESCRIBE SECURITY INTEGRATION \"azure\"", q)

	c := builder.Create()
	c.SetRaw(`TYPE=EXTERNAL_OAUTH`)
	c.SetString(`EXTERNAL_OAUTH_TYPE`, "AZURE")
	q = c.Statement()
	r.Equal(`CREATE SECURITY INTEGRATION "azure" TYPE=EXTERNAL_OAUTH EXTERNAL_OAUTH_TYPE='AZURE'`, q)

	e := builder.Drop()
	r.Equal(`DROP SECURITY INTEGRATION "azure"`, e)
}
