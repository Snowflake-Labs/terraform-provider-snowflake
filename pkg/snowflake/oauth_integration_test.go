package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestOAuthIntegration(t *testing.T) {
	r := require.New(t)
	builder := snowflake.NewOAuthIntegrationBuilder("tableau_desktop")
	r.NotNil(builder)

	q := builder.Show()
	r.Equal("SHOW SECURITY INTEGRATIONS LIKE 'tableau_desktop'", q)

	q = builder.Describe()
	r.Equal("DESCRIBE SECURITY INTEGRATION \"tableau_desktop\"", q)

	c := builder.Create()
	c.SetRaw(`TYPE=oauth`)
	c.SetString(`oauth_client`, "tableau_desktop")
	q = c.Statement()
	r.Equal(`CREATE SECURITY INTEGRATION "tableau_desktop" TYPE=oauth OAUTH_CLIENT='tableau_desktop'`, q)

	e := builder.Drop()
	r.Equal(`DROP SECURITY INTEGRATION "tableau_desktop"`, e)
}
