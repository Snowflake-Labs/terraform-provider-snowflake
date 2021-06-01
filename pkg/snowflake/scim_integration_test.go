package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestScimIntegration(t *testing.T) {
	r := require.New(t)
	builder := snowflake.ScimIntegration("aad_provisioning")
	r.NotNil(builder)

	q := builder.Show()
	r.Equal("SHOW SECURITY INTEGRATIONS LIKE 'aad_provisioning'", q)

	q = builder.Describe()
	r.Equal("DESCRIBE SECURITY INTEGRATION \"aad_provisioning\"", q)

	c := builder.Create()
	c.SetRaw(`type=scim`)
	c.SetString(`scim_client`, "azure")
	c.SetString(`run_as_role`, "AAD_PROVISIONER")
	q = c.Statement()
	r.Equal(`CREATE SECURITY INTEGRATION "aad_provisioning" type=scim RUN_AS_ROLE='AAD_PROVISIONER' SCIM_CLIENT='azure'`, q)

	d := builder.Alter()
	d.SetRaw(`type=scim`)
	d.SetString(`scim_client`, "azure")
	d.SetString(`run_as_role`, "AAD_PROVISIONER")
	d.SetString(`network_policy`, "aad_policy")
	q = d.Statement()
	r.Equal(`ALTER SECURITY INTEGRATION "aad_provisioning" SET type=scim SCIM_CLIENT='azure' RUN_AS_ROLE='AAD_PROVISIONER' NETWORK_POLICY='aad_policy'`, q)

	e := builder.Drop()
	r.Equal(`DROP SECURITY INTEGRATION "aad_provisioning"`, e)
}
