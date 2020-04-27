package resources_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/require"
)

func database(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func databaseGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.DatabaseGrant().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func resourceMonitorGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ResourceMonitorGrant().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func integrationGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.IntegrationGrant().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func accountGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.AccountGrant().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func providers() map[string]terraform.ResourceProvider {
	p := provider.Provider()
	return map[string]terraform.ResourceProvider{
		"snowflake": p,
	}
}

func role(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Role().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func roleGrants(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RoleGrants().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func storageIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.StorageIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func user(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.User().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func warehouse(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Warehouse().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}
