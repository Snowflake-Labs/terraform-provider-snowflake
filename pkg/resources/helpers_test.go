package resources_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
)

// todo: remove the rest of these which are not used. also this file should be renamed for clarity to make it clear it is for testing only
// https://snowflakecomputing.atlassian.net/browse/SNOW-936093
type grantType int

const (
	normal grantType = iota
	onFuture
	onAll
)

func TestGetPropertyAsPointer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"integer": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"boolean": {
			Type:     schema.TypeBool,
			Required: true,
		},
	}, map[string]interface{}{
		"integer": 123,
		"string":  "some string",
		"boolean": true,
		"invalid": true,
	})

	assert.Equal(t, 123, *resources.GetPropertyAsPointer[int](d, "integer"))
	assert.Equal(t, "some string", *resources.GetPropertyAsPointer[string](d, "string"))
	assert.Equal(t, true, *resources.GetPropertyAsPointer[bool](d, "boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "invalid"))
}

func database(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

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

func managedAccount(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ManagedAccount().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func networkPolicy(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.NetworkPolicy().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func pipe(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Pipe().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func resourceMonitor(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ResourceMonitor().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func sequence(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Sequence().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func share(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Share().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func stage(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Stage().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func stream(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Stream().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func tag(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Tag().Schema, params)
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

func apiIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.APIIntegration().Schema, params)
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

func externalFunction(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ExternalFunction().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func procedure(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Procedure().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func storageIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.StorageIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func notificationIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.NotificationIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func table(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Table().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func externalTable(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.ExternalTable().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func user(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.User().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func view(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.View().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func materializedView(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.MaterializedView().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func warehouse(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Warehouse().Schema, params)
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

func rowAccessPolicy(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.RowAccessPolicy().Schema, params)
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

func function(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.Function().Schema, params)
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

// queriedAccountRolePrivilegesEqualTo will check if all the privileges specified in the argument are granted in Snowflake.
func queriedPrivilegesEqualTo(query func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error), privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()
		grants, err := query(client, ctx)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			if (grant.GrantTo == sdk.ObjectTypeDatabaseRole || grant.GrantedTo == sdk.ObjectTypeDatabaseRole) && grant.Privilege == "USAGE" {
				continue
			}
			if !slices.Contains(privileges, grant.Privilege) {
				return fmt.Errorf("grant not expected, grant: %v, not in %v", grants, privileges)
			}
		}

		return nil
	}
}

// queriedAccountRolePrivilegesContainAtLeast will check if all the privileges specified in the argument are granted in Snowflake.
// Any additional grants will be ignored.
func queriedPrivilegesContainAtLeast(query func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error), roleName sdk.ObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		grants, err := query(client, ctx)
		if err != nil {
			return err
		}
		var grantedPrivileges []string
		for _, grant := range grants {
			grantedPrivileges = append(grantedPrivileges, grant.Privilege)
		}
		notAllPrivilegesInGrantedPrivileges := slices.ContainsFunc(privileges, func(privilege string) bool {
			return !slices.Contains(grantedPrivileges, privilege)
		})
		if notAllPrivilegesInGrantedPrivileges {
			return fmt.Errorf("not every privilege from the list: %v was found in grant privileges: %v, for role name: %s", privileges, grantedPrivileges, roleName.FullyQualifiedName())
		}

		return nil
	}
}

// TODO(SNOW-936093): This function should be merged with testint/helpers_test.go updateAccountParameterTemporarily function which does the same thing.
// We cannot use it right now because it requires moving the function between the packages, so both tests will be able to see it.
func updateAccountParameter(t *testing.T, client *sdk.Client, parameter sdk.AccountParameter, temporarily bool, newValue string) func() {
	t.Helper()

	ctx := context.Background()

	param, err := client.Parameters.ShowAccountParameter(ctx, parameter)
	require.NoError(t, err)
	oldValue := param.Value

	if temporarily {
		t.Cleanup(func() {
			err = client.Parameters.SetAccountParameter(ctx, parameter, oldValue)
			require.NoError(t, err)
		})
	}

	return func() {
		err = client.Parameters.SetAccountParameter(ctx, parameter, newValue)
		require.NoError(t, err)
	}
}

type PlanCheck func(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse)

func (fn PlanCheck) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	fn(ctx, req, resp)
}

func ExpectsCreatePlan(resourceAddress string) PlanCheck {
	return func(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
		for _, rc := range req.Plan.ResourceChanges {
			if rc.Address == resourceAddress && rc.Change != nil && slices.Contains(rc.Change.Actions, tfjson.ActionCreate) {
				return
			}
		}

		resp.Error = fmt.Errorf("expected plan to contain create request for %s", resourceAddress)
	}
}
