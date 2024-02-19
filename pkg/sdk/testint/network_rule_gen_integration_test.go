package testint

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInt_NetworkRules(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertValuesAndComment := func(id sdk.SchemaObjectIdentifier, expectedValueList []string, expectedComment string) {
		rule, err := client.NetworkRules.ShowByID(ctx, id)
		require.NoError(t, err)

		ruleDetails, err := client.NetworkRules.Describe(ctx, id)
		require.NoError(t, err)

		require.Len(t, expectedValueList, rule.EntriesInValueList)
		require.Equal(t, expectedValueList, ruleDetails.ValueList)
		require.Equal(t, expectedComment, rule.Comment)
	}

	t.Run("Create", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(20))
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress))
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.NetworkRules.Drop(ctx, sdk.NewDropNetworkRuleRequest(id))
			require.NoError(t, err)
		})

		_, err = client.NetworkRules.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Alter: set and unset", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(20))
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress))
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.NetworkRules.Drop(ctx, sdk.NewDropNetworkRuleRequest(id))
			require.NoError(t, err)
		})

		setReq := sdk.NewNetworkRuleSetRequest([]sdk.NetworkRuleValue{
			{"0.0.0.0"},
			{"1.1.1.1"},
		}).WithComment(sdk.String("some comment"))
		err = client.NetworkRules.Alter(ctx, sdk.NewAlterNetworkRuleRequest(id).WithSet(setReq))
		require.NoError(t, err)

		assertValuesAndComment(id, []string{"0.0.0.0", "1.1.1.1"}, "some comment")

		unsetReq := sdk.NewNetworkRuleUnsetRequest().
			WithValueList(sdk.Bool(true)).
			WithComment(sdk.Bool(true))
		err = client.NetworkRules.Alter(ctx, sdk.NewAlterNetworkRuleRequest(id).WithUnset(unsetReq))
		require.NoError(t, err)

		assertValuesAndComment(id, []string{}, "")
	})

	t.Run("Drop", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(20))
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress))
		require.NoError(t, err)

		_, err = client.NetworkRules.ShowByID(ctx, id)
		require.NoError(t, err)

		err = client.NetworkRules.Drop(ctx, sdk.NewDropNetworkRuleRequest(id))
		require.NoError(t, err)

		_, err = client.NetworkRules.ShowByID(ctx, id)
		require.Error(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("Show", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(20))
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress).WithComment(sdk.String("some comment")))
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.NetworkRules.Drop(ctx, sdk.NewDropNetworkRuleRequest(id))
			require.NoError(t, err)
		})

		networkRules, err := client.NetworkRules.Show(ctx, sdk.NewShowNetworkRuleRequest().WithIn(&sdk.In{
			Schema: sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName()),
		}).WithLike(&sdk.Like{
			Pattern: sdk.String(id.Name()),
		}))
		require.NoError(t, err)

		require.Equal(t, 1, len(networkRules))
		require.False(t, networkRules[0].CreatedOn.IsZero())
		require.Equal(t, id.Name(), networkRules[0].Name)
		require.Equal(t, id.DatabaseName(), networkRules[0].DatabaseName)
		require.Equal(t, id.SchemaName(), networkRules[0].SchemaName)
		require.Equal(t, "ACCOUNTADMIN", networkRules[0].Owner)
		require.Equal(t, "some comment", networkRules[0].Comment)
		require.Equal(t, sdk.NetworkRuleTypeIpv4, networkRules[0].Type)
		require.Equal(t, sdk.NetworkRuleModeIngress, networkRules[0].Mode)
		require.Equal(t, 0, networkRules[0].EntriesInValueList)
		require.Equal(t, "ROLE", networkRules[0].OwnerRoleType)
	})

	t.Run("Describe", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(20))
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress).WithComment(sdk.String("some comment")))
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.NetworkRules.Drop(ctx, sdk.NewDropNetworkRuleRequest(id))
			require.NoError(t, err)
		})

		ruleDetails, err := client.NetworkRules.Describe(ctx, id)
		require.NoError(t, err)
		assert.False(t, ruleDetails.CreatedOn.IsZero())
		assert.Equal(t, id.DatabaseName(), ruleDetails.DatabaseName)
		assert.Equal(t, id.SchemaName(), ruleDetails.SchemaName)
		assert.Equal(t, id.Name(), ruleDetails.Name)
		require.Equal(t, "ACCOUNTADMIN", ruleDetails.Owner)
		assert.Equal(t, "some comment", ruleDetails.Comment)
		assert.Empty(t, ruleDetails.ValueList)
		assert.Equal(t, sdk.NetworkRuleModeIngress, ruleDetails.Mode)
		assert.Equal(t, sdk.NetworkRuleTypeIpv4, ruleDetails.Type)
	})
}
