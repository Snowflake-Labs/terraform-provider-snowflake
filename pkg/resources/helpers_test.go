package resources_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	d := schema.TestResourceDataRaw(t, resources.DatabaseGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func schemaGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SchemaGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func stageGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.StageGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func tableGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.TableGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func viewGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ViewGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func materializedViewGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.MaterializedViewGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func resourceMonitorGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ResourceMonitorGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func integrationGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.IntegrationGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func accountGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.AccountGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func managedAccount(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ManagedAccount().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func maskingPolicy(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.MaskingPolicy().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func networkPolicy(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.NetworkPolicy().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func pipe(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Pipe().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func resourceMonitor(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ResourceMonitor().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func sequence(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Sequence().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func share(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Share().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func stage(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Stage().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func stream(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Stream().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func tag(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Tag().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func providers() map[string]*schema.Provider {
	p := provider.Provider()
	return map[string]*schema.Provider{
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

func apiIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.APIIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func samlIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SAMLIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func scimIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SCIMIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func oauthIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.OAuthIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func externalFunction(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ExternalFunction().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func procedure(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Procedure().Schema, params)
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
func notificationIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.NotificationIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func table(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Table().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func externalTable(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ExternalTable().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func task(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Task().Schema, params)
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

func view(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.View().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func materializedView(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.MaterializedView().Schema, params)
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

func externalTableGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ExternalTableGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func fileFormatGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.FileFormatGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func sequenceGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SequenceGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func streamGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.StreamGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func functionGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.FunctionGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func procedureGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ProcedureGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func maskingPolicyGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.MaskingPolicyGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func pipeGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.PipeGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func taskGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.TaskGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func rowAccessPolicy(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RowAccessPolicy().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func rowAccessPolicyGrant(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RowAccessPolicyGrant().Resource.Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func function(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Function().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}
