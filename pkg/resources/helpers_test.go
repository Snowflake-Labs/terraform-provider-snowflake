package resources_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"
)

func warehouse(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.Warehouse().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}

func database(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, params)
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
func providers() map[string]terraform.ResourceProvider {
	p := provider.Provider()
	return map[string]terraform.ResourceProvider{
		"snowflake": p,
	}
}
