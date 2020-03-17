package resources_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func database(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}

func databaseGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.DatabaseGrant().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}

func accountGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.AccountGrant().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}

func fixture(name string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", name))
	return string(b), err
}

func providers() map[string]terraform.ResourceProvider {
	p := provider.Provider()
	return map[string]terraform.ResourceProvider{
		"snowflake": p,
	}
}

func role(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.Role().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}

func roleGrants(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.RoleGrants().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}

func storageIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.StorageIntegration().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}

func user(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.User().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}

func warehouse(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.Warehouse().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}
