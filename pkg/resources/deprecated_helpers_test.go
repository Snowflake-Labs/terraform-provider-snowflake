package resources_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

/**
 * Most functions from this file will be removed while removing old grant resources.
 * The rest will be removed while adding security integrations to the SDK.
 */

func databaseGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.DatabaseGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func schemaGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SchemaGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func stageGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.StageGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func tableGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.TableGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func viewGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ViewGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func materializedViewGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.MaterializedViewGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func resourceMonitorGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ResourceMonitorGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func integrationGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.IntegrationGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func accountGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.AccountGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func roleGrants(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RoleGrants().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func userOwnershipGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.UserOwnershipGrant().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func roleOwnershipGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RoleOwnershipGrant().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func samlIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SAMLIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func scimIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SCIMIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func oauthIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.OAuthIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func externalTableGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ExternalTableGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func fileFormatGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.FileFormatGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func sequenceGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SequenceGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func streamGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.StreamGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func maskingPolicyGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.MaskingPolicyGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func pipeGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.PipeGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func taskGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.TaskGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func rowAccessPolicyGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RowAccessPolicyGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func tagGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.TagGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}
