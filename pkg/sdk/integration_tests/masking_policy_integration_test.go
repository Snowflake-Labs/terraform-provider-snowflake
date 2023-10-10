package sdk_integration_tests

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_MaskingPoliciesShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	maskingPolicyTest, maskingPolicyCleanup := sdk.createMaskingPolicy(t, client, databaseTest, schemaTest)
	t.Cleanup(maskingPolicyCleanup)

	maskingPolicy2Test, maskingPolicy2Cleanup := sdk.createMaskingPolicy(t, client, databaseTest, schemaTest)
	t.Cleanup(maskingPolicy2Cleanup)

	t.Run("without show options", func(t *testing.T) {
		useDatabaseCleanup := sdk.useDatabase(t, client, databaseTest.ID())
		t.Cleanup(useDatabaseCleanup)
		useSchemaCleanup := sdk.useSchema(t, client, schemaTest.ID())
		t.Cleanup(useSchemaCleanup)

		maskingPolicies, err := client.MaskingPolicies.Show(ctx, nil)
		require.NoError(t, err)
		assert.Equal(t, 2, len(maskingPolicies))
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &sdk.ShowMaskingPolicyOptions{
			In: &sdk.In{
				Schema: schemaTest.ID(),
			},
		}
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, maskingPolicies, *maskingPolicyTest)
		assert.Contains(t, maskingPolicies, *maskingPolicy2Test)
		assert.Equal(t, 2, len(maskingPolicies))
	})

	t.Run("with show options and like", func(t *testing.T) {
		showOptions := &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(maskingPolicyTest.Name),
			},
			In: &sdk.In{
				Database: databaseTest.ID(),
			},
		}
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, maskingPolicies, *maskingPolicyTest)
		assert.Equal(t, 1, len(maskingPolicies))
	})

	t.Run("when searching a non-existent masking policy", func(t *testing.T) {
		showOptions := &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(maskingPolicies))
	})

	/*
		// there appears to be a bug in the Snowflake API. LIMIT is not actually limiting the number of results
		t.Run("when limiting the number of results", func(t *testing.T) {
			showOptions := &MaskingPolicyShowOptions{
				In: &In{
					Schema: schemaTest.ID(),
				},
				Limit: Int(1),
			}
			maskingPolicies, err := client.MaskingPolicies.Show(ctx, showOptions)
			require.NoError(t, err)
			assert.Equal(t, 1, len(maskingPolicies))
		})
	*/
}

func TestInt_MaskingPolicyCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	t.Run("test complete case", func(t *testing.T) {
		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		signature := []sdk.TableColumnSignature{
			{
				Name: "col1",
				Type: sdk.DataTypeVARCHAR,
			},
			{
				Name: "col2",
				Type: sdk.DataTypeVARCHAR,
			},
		}
		expression := "REPLACE('X', 1, 2)"
		comment := randomComment(t)
		exemptOtherPolicies := randomBool(t)
		err := client.MaskingPolicies.Create(ctx, id, signature, sdk.DataTypeVARCHAR, expression, &sdk.CreateMaskingPolicyOptions{
			OrReplace:           sdk.Bool(true),
			IfNotExists:         sdk.Bool(false),
			Comment:             sdk.String(comment),
			ExemptOtherPolicies: sdk.Bool(exemptOtherPolicies),
		})
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, signature, maskingPolicyDetails.Signature)
		assert.Equal(t, sdk.DataTypeVARCHAR, maskingPolicyDetails.ReturnType)
		assert.Equal(t, expression, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(maskingPolicy))
		assert.Equal(t, name, maskingPolicy[0].Name)
		assert.Equal(t, comment, maskingPolicy[0].Comment)
		assert.Equal(t, exemptOtherPolicies, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("test if_not_exists", func(t *testing.T) {
		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		signature := []sdk.TableColumnSignature{
			{
				Name: "col1",
				Type: sdk.DataTypeVARCHAR,
			},
			{
				Name: "col2",
				Type: sdk.DataTypeVARCHAR,
			},
		}
		expression := "REPLACE('X', 1, 2)"
		comment := randomComment(t)
		err := client.MaskingPolicies.Create(ctx, id, signature, sdk.DataTypeVARCHAR, expression, &sdk.CreateMaskingPolicyOptions{
			OrReplace:           sdk.Bool(false),
			IfNotExists:         sdk.Bool(true),
			Comment:             sdk.String(comment),
			ExemptOtherPolicies: sdk.Bool(true),
		})
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, signature, maskingPolicyDetails.Signature)
		assert.Equal(t, sdk.DataTypeVARCHAR, maskingPolicyDetails.ReturnType)
		assert.Equal(t, expression, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(maskingPolicy))
		assert.Equal(t, name, maskingPolicy[0].Name)
		assert.Equal(t, comment, maskingPolicy[0].Comment)
		assert.Equal(t, true, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("test no options", func(t *testing.T) {
		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		signature := []sdk.TableColumnSignature{
			{
				Name: "col1",
				Type: sdk.DataTypeVARCHAR,
			},
		}
		expression := "REPLACE('X', 1, 2)"
		err := client.MaskingPolicies.Create(ctx, id, signature, sdk.DataTypeVARCHAR, expression, nil)
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, signature, maskingPolicyDetails.Signature)
		assert.Equal(t, sdk.DataTypeVARCHAR, maskingPolicyDetails.ReturnType)
		assert.Equal(t, expression, maskingPolicyDetails.Body)

		maskingPolicy, err := client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(maskingPolicy))
		assert.Equal(t, name, maskingPolicy[0].Name)
		assert.Equal(t, "", maskingPolicy[0].Comment)
		assert.Equal(t, false, maskingPolicy[0].ExemptOtherPolicies)
	})

	t.Run("test multiline expression", func(t *testing.T) {
		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		signature := []sdk.TableColumnSignature{
			{
				Name: "val",
				Type: sdk.DataTypeVARCHAR,
			},
		}
		expression := `
		case
			when current_role() in ('ROLE_A') then
				val
			when is_role_in_session( 'ROLE_B' ) then
				'ABC123'
			else
				'******'
		end
		`
		err := client.MaskingPolicies.Create(ctx, id, signature, sdk.DataTypeVARCHAR, expression, nil)
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, maskingPolicyDetails.Name)
		assert.Equal(t, signature, maskingPolicyDetails.Signature)
		assert.Equal(t, sdk.DataTypeVARCHAR, maskingPolicyDetails.ReturnType)
		assert.Equal(t, strings.TrimSpace(expression), maskingPolicyDetails.Body)
	})
}

func TestInt_MaskingPolicyDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	maskingPolicy, maskingPolicyCleanup := sdk.createMaskingPolicy(t, client, databaseTest, schemaTest)
	t.Cleanup(maskingPolicyCleanup)

	t.Run("when masking policy exists", func(t *testing.T) {
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, maskingPolicy.ID())
		require.NoError(t, err)
		assert.Equal(t, maskingPolicy.Name, maskingPolicyDetails.Name)
	})

	t.Run("when masking policy does not exist", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, "does_not_exist")
		_, err := client.MaskingPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})
}

func TestInt_MaskingPolicyAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	t.Run("when setting and unsetting a value", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := sdk.createMaskingPolicy(t, client, databaseTest, schemaTest)
		t.Cleanup(maskingPolicyCleanup)
		comment := randomComment(t)
		alterOptions := &sdk.AlterMaskingPolicyOptions{
			Set: &sdk.MaskingPolicySet{
				Comment: sdk.String(comment),
			},
		}
		err := client.MaskingPolicies.Alter(ctx, maskingPolicy.ID(), alterOptions)
		require.NoError(t, err)
		maskingPolicies, err := client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(maskingPolicy.Name),
			},
			In: &sdk.In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(maskingPolicies))
		assert.Equal(t, comment, maskingPolicies[0].Comment)

		err = client.MaskingPolicies.Alter(ctx, maskingPolicy.ID(), alterOptions)
		require.NoError(t, err)
		alterOptions = &sdk.AlterMaskingPolicyOptions{
			Unset: &sdk.MaskingPolicyUnset{
				Comment: sdk.Bool(true),
			},
		}
		err = client.MaskingPolicies.Alter(ctx, maskingPolicy.ID(), alterOptions)
		require.NoError(t, err)
		maskingPolicies, err = client.MaskingPolicies.Show(ctx, &sdk.ShowMaskingPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(maskingPolicy.Name),
			},
			In: &sdk.In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(maskingPolicies))
		assert.Equal(t, "", maskingPolicies[0].Comment)
	})

	t.Run("when renaming", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := sdk.createMaskingPolicy(t, client, databaseTest, schemaTest)
		oldID := maskingPolicy.ID()
		t.Cleanup(maskingPolicyCleanup)
		newName := randomString(t)
		newID := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, newName)
		alterOptions := &sdk.AlterMaskingPolicyOptions{
			NewName: newID,
		}
		err := client.MaskingPolicies.Alter(ctx, oldID, alterOptions)
		require.NoError(t, err)
		maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newName, maskingPolicyDetails.Name)
		// rename back to original name so it can be cleaned up
		alterOptions = &sdk.AlterMaskingPolicyOptions{
			NewName: oldID,
		}
		err = client.MaskingPolicies.Alter(ctx, newID, alterOptions)
		require.NoError(t, err)
	})

	t.Run("setting and unsetting tags", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := sdk.createMaskingPolicy(t, client, databaseTest, schemaTest)
		id := maskingPolicy.ID()
		t.Cleanup(maskingPolicyCleanup)

		tag, tagCleanup := createTag(t, client, databaseTest, schemaTest)
		t.Cleanup(tagCleanup)

		tag2, tag2Cleanup := createTag(t, client, databaseTest, schemaTest)
		t.Cleanup(tag2Cleanup)

		tagAssociations := []sdk.TagAssociation{{Name: tag.ID(), Value: "value1"}, {Name: tag2.ID(), Value: "value2"}}
		alterOptions := &sdk.AlterMaskingPolicyOptions{
			Set: &sdk.MaskingPolicySet{
				Tag: tagAssociations,
			},
		}
		err := client.MaskingPolicies.Alter(ctx, id, alterOptions)
		require.NoError(t, err)
		tagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeMaskingPolicy)
		require.NoError(t, err)
		assert.Equal(t, tagAssociations[0].Value, tagValue)
		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2.ID(), id, sdk.ObjectTypeMaskingPolicy)
		require.NoError(t, err)
		assert.Equal(t, tagAssociations[1].Value, tag2Value)

		// unset tag
		alterOptions = &sdk.AlterMaskingPolicyOptions{
			Unset: &sdk.MaskingPolicyUnset{
				Tag: []sdk.ObjectIdentifier{tag.ID()},
			},
		}
		err = client.MaskingPolicies.Alter(ctx, id, alterOptions)
		require.NoError(t, err)
		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeMaskingPolicy)
		assert.Error(t, err)
	})
}

func TestInt_MaskingPolicyDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	t.Run("when masking policy exists", func(t *testing.T) {
		maskingPolicy, _ := sdk.createMaskingPolicy(t, client, databaseTest, schemaTest)
		id := maskingPolicy.ID()
		err := client.MaskingPolicies.Drop(ctx, id)
		require.NoError(t, err)
		_, err = client.MaskingPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("when masking policy does not exist", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, "does_not_exist")
		err := client.MaskingPolicies.Drop(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})
}
