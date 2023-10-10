package sdk_integration_tests

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_PasswordPoliciesShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	passwordPolicyTest, passwordPolicyCleanup := createPasswordPolicy(t, client, databaseTest, schemaTest)
	t.Cleanup(passwordPolicyCleanup)

	passwordPolicy2Test, passwordPolicy2Cleanup := createPasswordPolicy(t, client, databaseTest, schemaTest)
	t.Cleanup(passwordPolicy2Cleanup)

	t.Run("without show options", func(t *testing.T) {
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(passwordPolicies))
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &sdk.ShowPasswordPolicyOptions{
			In: &sdk.In{
				Schema: schemaTest.ID(),
			},
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, passwordPolicies, *passwordPolicyTest)
		assert.Contains(t, passwordPolicies, *passwordPolicy2Test)
		assert.Equal(t, 2, len(passwordPolicies))
	})

	t.Run("with show options and like", func(t *testing.T) {
		showOptions := &sdk.ShowPasswordPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(passwordPolicyTest.Name),
			},
			In: &sdk.In{
				Database: databaseTest.ID(),
			},
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, passwordPolicies, *passwordPolicyTest)
		assert.Equal(t, 1, len(passwordPolicies))
	})

	t.Run("when searching a non-existent password policy", func(t *testing.T) {
		showOptions := &sdk.ShowPasswordPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(passwordPolicies))
	})

	/* there appears to be a bug in the Snowflake API. LIMIT is not actually limiting the number of results
	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &ShowPasswordPolicyOptions{
			In: &In{
				Schema: String(schemaTest.FullyQualifiedName()),
			},
			Limit: Int(1),
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(passwordPolicies))
	})*/
}

func TestInt_PasswordPolicyCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)
	t.Run("test complete", func(t *testing.T) {
		name := randomUUID(t)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.PasswordPolicies.Create(ctx, id, &sdk.CreatePasswordPolicyOptions{
			OrReplace:                 sdk.Bool(true),
			PasswordMinLength:         sdk.Int(10),
			PasswordMaxLength:         sdk.Int(20),
			PasswordMinUpperCaseChars: sdk.Int(1),
			PasswordMinLowerCaseChars: sdk.Int(1),
			PasswordMinNumericChars:   sdk.Int(1),
			PasswordMinSpecialChars:   sdk.Int(1),
			PasswordMaxAgeDays:        sdk.Int(30),
			PasswordMaxRetries:        sdk.Int(5),
			PasswordLockoutTimeMins:   sdk.Int(30),
			// todo [SNOW-928909]: uncomment this once comments are working again
			// Comment:                   String("test comment"),
		})
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, passwordPolicyDetails.Name.Value)
		// todo [SNOW-928909]: uncomment this once comments are working again
		// assert.Equal(t, "test comment", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, 10, *passwordPolicyDetails.PasswordMinLength.Value)
		assert.Equal(t, 20, *passwordPolicyDetails.PasswordMaxLength.Value)
		assert.Equal(t, 1, *passwordPolicyDetails.PasswordMinUpperCaseChars.Value)
		assert.Equal(t, 1, *passwordPolicyDetails.PasswordMinLowerCaseChars.Value)
		assert.Equal(t, 1, *passwordPolicyDetails.PasswordMinNumericChars.Value)
		assert.Equal(t, 1, *passwordPolicyDetails.PasswordMinSpecialChars.Value)
		assert.Equal(t, 30, *passwordPolicyDetails.PasswordMaxAgeDays.Value)
		assert.Equal(t, 5, *passwordPolicyDetails.PasswordMaxRetries.Value)
		assert.Equal(t, 30, *passwordPolicyDetails.PasswordLockoutTimeMins.Value)
	})

	t.Run("test if_not_exists", func(t *testing.T) {
		name := randomUUID(t)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.PasswordPolicies.Create(ctx, id, &sdk.CreatePasswordPolicyOptions{
			OrReplace:                 sdk.Bool(false),
			IfNotExists:               sdk.Bool(true),
			PasswordMinLength:         sdk.Int(10),
			PasswordMaxLength:         sdk.Int(20),
			PasswordMinUpperCaseChars: sdk.Int(5),
			// todo [SNOW-928909]: uncomment this once comments are working again
			// Comment:                   String("test comment"),
		})
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, passwordPolicyDetails.Name.Value)
		// todo [SNOW-928909]: uncomment this once comments are working again
		// assert.Equal(t, "test comment", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, 10, *passwordPolicyDetails.PasswordMinLength.Value)
		assert.Equal(t, 20, *passwordPolicyDetails.PasswordMaxLength.Value)
		assert.Equal(t, 5, *passwordPolicyDetails.PasswordMinUpperCaseChars.Value)
	})

	t.Run("test no options", func(t *testing.T) {
		name := randomUUID(t)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.PasswordPolicies.Create(ctx, id, nil)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, "", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, *passwordPolicyDetails.PasswordMinLength.Value, *passwordPolicyDetails.PasswordMinLength.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMaxLength.Value, *passwordPolicyDetails.PasswordMaxLength.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMinUpperCaseChars.Value, *passwordPolicyDetails.PasswordMinUpperCaseChars.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMinLowerCaseChars.Value, *passwordPolicyDetails.PasswordMinLowerCaseChars.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMinNumericChars.Value, *passwordPolicyDetails.PasswordMinNumericChars.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMinSpecialChars.Value, *passwordPolicyDetails.PasswordMinSpecialChars.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMaxAgeDays.Value, *passwordPolicyDetails.PasswordMaxAgeDays.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMaxRetries.Value, *passwordPolicyDetails.PasswordMaxRetries.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordLockoutTimeMins.Value, *passwordPolicyDetails.PasswordLockoutTimeMins.DefaultValue)
	})
}

func TestInt_PasswordPolicyDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	passwordPolicy, passwordPolicyCleanup := createPasswordPolicy(t, client, databaseTest, schemaTest)
	t.Cleanup(passwordPolicyCleanup)

	t.Run("when password policy exists", func(t *testing.T) {
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, passwordPolicy.ID())
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, passwordPolicy.Comment, passwordPolicyDetails.Comment.Value)
	})

	t.Run("when password policy does not exist", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, "does_not_exist")
		_, err := client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})
}

func TestInt_PasswordPolicyAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	t.Run("when setting new values", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := createPasswordPolicy(t, client, databaseTest, schemaTest)
		t.Cleanup(passwordPolicyCleanup)
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMinLength: sdk.Int(10),
				PasswordMaxLength: sdk.Int(20),
				Comment:           sdk.String("new comment"),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, passwordPolicy.ID(), alterOptions)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, passwordPolicy.ID())
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, 10, *passwordPolicyDetails.PasswordMinLength.Value)
		assert.Equal(t, 20, *passwordPolicyDetails.PasswordMaxLength.Value)
		assert.Equal(t, "new comment", passwordPolicyDetails.Comment.Value)
	})

	t.Run("when renaming", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := createPasswordPolicy(t, client, databaseTest, schemaTest)
		oldID := passwordPolicy.ID()
		t.Cleanup(passwordPolicyCleanup)
		newName := randomUUID(t)
		newID := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, newName)
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			NewName: newID,
		}
		err := client.PasswordPolicies.Alter(ctx, oldID, alterOptions)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, newID)
		require.NoError(t, err)
		// rename back to original name so it can be cleaned up
		assert.Equal(t, newName, passwordPolicyDetails.Name.Value)
		alterOptions = &sdk.AlterPasswordPolicyOptions{
			NewName: oldID,
		}
		err = client.PasswordPolicies.Alter(ctx, newID, alterOptions)
		require.NoError(t, err)
	})

	t.Run("when unsetting values", func(t *testing.T) {
		createOptions := &sdk.CreatePasswordPolicyOptions{
			PasswordMaxAgeDays: sdk.Int(20),
			PasswordMaxRetries: sdk.Int(10),
			// todo [SNOW-928909]: uncomment this once comments are working again
			// Comment: String("test comment")
		}
		passwordPolicy, passwordPolicyCleanup := createPasswordPolicyWithOptions(t, client, databaseTest, schemaTest, createOptions)
		id := passwordPolicy.ID()
		t.Cleanup(passwordPolicyCleanup)
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Unset: &sdk.PasswordPolicyUnset{
				PasswordMaxRetries: sdk.Bool(true),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, id, alterOptions)
		require.NoError(t, err)
		alterOptions = &sdk.AlterPasswordPolicyOptions{
			Unset: &sdk.PasswordPolicyUnset{
				PasswordMaxAgeDays: sdk.Bool(true),
				// todo [SNOW-928909]: uncomment this once comments are working again
				// Comment: Bool("true")
			},
		}
		err = client.PasswordPolicies.Alter(ctx, id, alterOptions)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, "", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, *passwordPolicyDetails.PasswordMaxRetries.Value, *passwordPolicyDetails.PasswordMaxRetries.DefaultValue)
	})

	t.Run("when unsetting multiple values at same time", func(t *testing.T) {
		createOptions := &sdk.CreatePasswordPolicyOptions{
			PasswordMaxAgeDays: sdk.Int(20),
			PasswordMaxRetries: sdk.Int(10),
			// todo [SNOW-928909]: uncomment this once comments are working again
			// Comment: String("test comment")
		}
		passwordPolicy, passwordPolicyCleanup := createPasswordPolicyWithOptions(t, client, databaseTest, schemaTest, createOptions)
		id := passwordPolicy.ID()
		t.Cleanup(passwordPolicyCleanup)
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Unset: &sdk.PasswordPolicyUnset{
				PasswordMaxAgeDays: sdk.Bool(true),
				PasswordMaxRetries: sdk.Bool(true),
				// todo [SNOW-928909]: uncomment this once comments are working again
				// Comment: Bool("true")
			},
		}
		err := client.PasswordPolicies.Alter(ctx, id, alterOptions)
		require.Error(t, err)
	})
}

func TestInt_PasswordPolicyDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	t.Run("when password policy exists", func(t *testing.T) {
		passwordPolicy, _ := createPasswordPolicy(t, client, databaseTest, schemaTest)
		id := passwordPolicy.ID()
		err := client.PasswordPolicies.Drop(ctx, id, nil)
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("when password policy does not exist", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, "does_not_exist")
		err := client.PasswordPolicies.Drop(ctx, id, nil)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("when password policy exists and if exists is true", func(t *testing.T) {
		passwordPolicy, _ := createPasswordPolicy(t, client, databaseTest, schemaTest)
		id := passwordPolicy.ID()
		dropOptions := &sdk.DropPasswordPolicyOptions{IfExists: sdk.Bool(true)}
		err := client.PasswordPolicies.Drop(ctx, id, dropOptions)
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})
}
