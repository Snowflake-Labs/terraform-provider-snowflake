package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordPoliciesShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
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
		showOptions := &PasswordPolicyShowOptions{
			In: &In{
				Schema: schemaTest.Identifier(),
			},
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, passwordPolicies, passwordPolicyTest)
		assert.Contains(t, passwordPolicies, passwordPolicy2Test)
		assert.Equal(t, 2, len(passwordPolicies))
	})

	t.Run("with show options and like", func(t *testing.T) {
		showOptions := &PasswordPolicyShowOptions{
			Like: &Like{
				Pattern: String(passwordPolicyTest.Name),
			},
			In: &In{
				Database: databaseTest.Identifier(),
			},
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, passwordPolicies, passwordPolicyTest)
		assert.Equal(t, 1, len(passwordPolicies))
	})

	t.Run("when searching a non-existent password policy", func(t *testing.T) {
		showOptions := &PasswordPolicyShowOptions{
			Like: &Like{
				Pattern: String("non-existent"),
			},
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(passwordPolicies))
	})

	/* there appears to be a bug in the Snowflake API. LIMIT is not actually limiting the number of results
	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &PasswordPolicyShowOptions{
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

func TestPasswordPolicyCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	t.Run("test complete case", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.PasswordPolicies.Create(ctx, id, &PasswordPolicyCreateOptions{
			OrReplace:                 Bool(true),
			PasswordMinLength:         Int(10),
			PasswordMaxLength:         Int(20),
			PasswordMinUpperCaseChars: Int(1),
			PasswordMinLowerCaseChars: Int(1),
			PasswordMinNumericChars:   Int(1),
			PasswordMinSpecialChars:   Int(1),
			PasswordMaxAgeDays:        Int(30),
			PasswordMaxRetries:        Int(5),
			PasswordLockoutTimeMins:   Int(30),
			Comment:                   String("test comment"),
		})
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, "test comment", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, 10, passwordPolicyDetails.PasswordMinLength.Value)
		assert.Equal(t, 20, passwordPolicyDetails.PasswordMaxLength.Value)
		assert.Equal(t, 1, passwordPolicyDetails.PasswordMinUpperCaseChars.Value)
		assert.Equal(t, 1, passwordPolicyDetails.PasswordMinLowerCaseChars.Value)
		assert.Equal(t, 1, passwordPolicyDetails.PasswordMinNumericChars.Value)
		assert.Equal(t, 1, passwordPolicyDetails.PasswordMinSpecialChars.Value)
		assert.Equal(t, 30, passwordPolicyDetails.PasswordMaxAgeDays.Value)
		assert.Equal(t, 5, passwordPolicyDetails.PasswordMaxRetries.Value)
		assert.Equal(t, 30, passwordPolicyDetails.PasswordLockoutTimeMins.Value)
	})

	t.Run("test no on_on replace", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.PasswordPolicies.Create(ctx, id, &PasswordPolicyCreateOptions{
			OrReplace:                 Bool(false),
			PasswordMinLength:         Int(10),
			PasswordMaxLength:         Int(20),
			PasswordMinUpperCaseChars: Int(5),
			Comment:                   String("test comment"),
		})
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, "test comment", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, 10, passwordPolicyDetails.PasswordMinLength.Value)
		assert.Equal(t, 20, passwordPolicyDetails.PasswordMaxLength.Value)
		assert.Equal(t, 5, passwordPolicyDetails.PasswordMinUpperCaseChars.Value)
	})

	t.Run("test no options", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.PasswordPolicies.Create(ctx, id, nil)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, "", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, passwordPolicyDetails.PasswordMinLength.Value, passwordPolicyDetails.PasswordMinLength.DefaultValue)
		assert.Equal(t, passwordPolicyDetails.PasswordMaxLength.Value, passwordPolicyDetails.PasswordMaxLength.DefaultValue)
		assert.Equal(t, passwordPolicyDetails.PasswordMinUpperCaseChars.Value, passwordPolicyDetails.PasswordMinUpperCaseChars.DefaultValue)
		assert.Equal(t, passwordPolicyDetails.PasswordMinLowerCaseChars.Value, passwordPolicyDetails.PasswordMinLowerCaseChars.DefaultValue)
		assert.Equal(t, passwordPolicyDetails.PasswordMinNumericChars.Value, passwordPolicyDetails.PasswordMinNumericChars.DefaultValue)
		assert.Equal(t, passwordPolicyDetails.PasswordMinSpecialChars.Value, passwordPolicyDetails.PasswordMinSpecialChars.DefaultValue)
		assert.Equal(t, passwordPolicyDetails.PasswordMaxAgeDays.Value, passwordPolicyDetails.PasswordMaxAgeDays.DefaultValue)
		assert.Equal(t, passwordPolicyDetails.PasswordMaxRetries.Value, passwordPolicyDetails.PasswordMaxRetries.DefaultValue)
		assert.Equal(t, passwordPolicyDetails.PasswordLockoutTimeMins.Value, passwordPolicyDetails.PasswordLockoutTimeMins.DefaultValue)
	})
}

func TestPasswordPolicyDescribe(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	passwordPolicy, passwordPolicyCleanup := createPasswordPolicy(t, client, databaseTest, schemaTest)
	t.Cleanup(passwordPolicyCleanup)

	t.Run("when password policy exists", func(t *testing.T) {
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, passwordPolicy.Identifier())
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, passwordPolicy.Comment, passwordPolicyDetails.Comment.Value)
	})

	t.Run("when password policy does not exist", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, "does_not_exist")
		_, err := client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}

func TestPasswordPolicyAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	t.Run("when setting new values", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := createPasswordPolicy(t, client, databaseTest, schemaTest)
		t.Cleanup(passwordPolicyCleanup)
		alterOptions := &PasswordPolicyAlterOptions{
			Set: &PasswordPolicyAlterSet{
				PasswordMinLength: Int(10),
				PasswordMaxLength: Int(20),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, passwordPolicy.Identifier(), alterOptions)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, passwordPolicy.Identifier())
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, 10, passwordPolicyDetails.PasswordMinLength.Value)
		assert.Equal(t, 20, passwordPolicyDetails.PasswordMaxLength.Value)
	})

	t.Run("when renaming", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := createPasswordPolicy(t, client, databaseTest, schemaTest)
		oldID := passwordPolicy.Identifier()
		t.Cleanup(passwordPolicyCleanup)
		newName := randomString(t)
		newID := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, newName)
		alterOptions := &PasswordPolicyAlterOptions{
			NewName: newID,
		}
		err := client.PasswordPolicies.Alter(ctx, oldID, alterOptions)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, newID)
		require.NoError(t, err)
		// rename back to original name so it can be cleaned up
		assert.Equal(t, newName, passwordPolicyDetails.Name.Value)
		alterOptions = &PasswordPolicyAlterOptions{
			NewName: oldID,
		}
		err = client.PasswordPolicies.Alter(ctx, newID, alterOptions)
		require.NoError(t, err)
	})

	t.Run("when unsetting values", func(t *testing.T) {
		createOptions := &PasswordPolicyCreateOptions{
			Comment:            String("test comment"),
			PasswordMaxRetries: Int(10),
		}
		passwordPolicy, passwordPolicyCleanup := createPasswordPolicyWithOptions(t, client, databaseTest, schemaTest, createOptions)
		id := passwordPolicy.Identifier()
		t.Cleanup(passwordPolicyCleanup)
		alterOptions := &PasswordPolicyAlterOptions{
			Unset: &PasswordPolicyAlterUnset{
				PasswordMaxRetries: Bool(true),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, id, alterOptions)
		require.NoError(t, err)
		alterOptions = &PasswordPolicyAlterOptions{
			Unset: &PasswordPolicyAlterUnset{
				Comment: Bool(true),
			},
		}
		err = client.PasswordPolicies.Alter(ctx, id, alterOptions)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, "", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, passwordPolicyDetails.PasswordMaxRetries.Value, passwordPolicyDetails.PasswordMaxRetries.DefaultValue)
	})

	t.Run("when unsetting multiple values at same time", func(t *testing.T) {
		createOptions := &PasswordPolicyCreateOptions{
			Comment:            String("test comment"),
			PasswordMaxRetries: Int(10),
		}
		passwordPolicy, passwordPolicyCleanup := createPasswordPolicyWithOptions(t, client, databaseTest, schemaTest, createOptions)
		id := passwordPolicy.Identifier()
		t.Cleanup(passwordPolicyCleanup)
		alterOptions := &PasswordPolicyAlterOptions{
			Unset: &PasswordPolicyAlterUnset{
				Comment:            Bool(true),
				PasswordMaxRetries: Bool(true),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, id, alterOptions)
		require.Error(t, err)
	})
}

func TestPasswordPolicyDrop(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	t.Run("when password policy exists", func(t *testing.T) {
		passwordPolicy, _ := createPasswordPolicy(t, client, databaseTest, schemaTest)
		id := passwordPolicy.Identifier()
		err := client.PasswordPolicies.Drop(ctx, id, nil)
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("when password policy does not exist", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, "does_not_exist")
		err := client.PasswordPolicies.Drop(ctx, id, nil)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("when password policy exists and if exists is true", func(t *testing.T) {
		passwordPolicy, _ := createPasswordPolicy(t, client, databaseTest, schemaTest)
		id := passwordPolicy.Identifier()
		dropOptions := &PasswordPolicyDropOptions{IfExists: Bool(true)}
		err := client.PasswordPolicies.Drop(ctx, id, dropOptions)
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}
