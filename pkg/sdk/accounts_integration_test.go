package sdk

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_AccountShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	ok, err := client.ContextFunctions.IsRoleInSession(ctx, NewAccountObjectIdentifier("ORGADMIN"))
	require.NoError(t, err)
	if !ok {
		t.Skip("ORGADMIN role is not in current session")
	}
	currentAccount, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	opts := &ShowAccountOptions{
		Like: &Like{
			Pattern: String(currentAccount),
		},
	}
	accounts, err := client.Accounts.Show(ctx, opts)
	require.NoError(t, err)
	assert.NotEmpty(t, accounts)
	assert.Equal(t, 1, len(accounts))
	assert.Contains(t, []string{accounts[0].AccountLocator, accounts[0].AccountName}, currentAccount)
}

func TestInt_AccountCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	ok, err := client.ContextFunctions.IsRoleInSession(ctx, NewAccountObjectIdentifier("ORGADMIN"))
	require.NoError(t, err)
	if !ok {
		t.Skip("ORGADMIN role is not in current session")
	}
	if _, ok := os.LookupEnv("SNOWFLAKE_TEST_ACCOUNT_CREATE"); !ok {
		t.Skip("Skipping TestInt_AccountCreate")
	}
	t.Run("complete case", func(t *testing.T) {
		accountID := NewAccountObjectIdentifier("TF_" + strings.ToUpper(gofakeit.Fruit()) + "_" + fmt.Sprintf("%d", (randomIntRange(t, 100, 999))))
		region, err := client.ContextFunctions.CurrentRegion(ctx)
		require.NoError(t, err)

		opts := &CreateAccountOptions{
			AdminName:          "someadmin",
			AdminPassword:      String(randomStringN(t, 12)),
			FirstName:          String("Ad"),
			LastName:           String("Min"),
			Email:              "admin@example.com",
			MustChangePassword: Bool(false),
			Edition:            EditionBusinessCritical,
			Comment:            String("Please delete me!"),
			Region:             String(region),
		}
		err = client.Accounts.Create(ctx, accountID, opts)
		require.NoError(t, err)
		account, err := client.Accounts.ShowByID(ctx, accountID)
		require.NoError(t, err)
		assert.Equal(t, accountID.Name(), account.AccountName)
		assert.Equal(t, EditionBusinessCritical, account.Edition)
		assert.Equal(t, "Please delete me!", account.Comment)
		assert.Equal(t, region, account.SnowflakeRegion)

		// rename
		newAccountID := NewAccountObjectIdentifier("TF_" + strings.ToUpper(gofakeit.Animal()) + "_" + fmt.Sprintf("%d", (randomIntRange(t, 100, 999))))
		alterOpts := &AlterAccountOptions{
			Rename: &AccountRename{
				Name:       accountID,
				NewName:    newAccountID,
				SaveOldURL: Bool(true),
			},
		}
		err = client.Accounts.Alter(ctx, alterOpts)
		require.NoError(t, err)
		account, err = client.Accounts.ShowByID(ctx, newAccountID)
		require.NoError(t, err)
		assert.Equal(t, newAccountID.Name(), account.AccountName)

		// drop old url
		alterOpts = &AlterAccountOptions{
			Drop: &AccountDrop{
				Name:   newAccountID,
				OldURL: Bool(true),
			},
		}
		err = client.Accounts.Alter(ctx, alterOpts)
		require.NoError(t, err)
		_, err = client.Accounts.ShowByID(ctx, newAccountID)
		require.NoError(t, err)
	})
}

func TestInt_AccountAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	ok, err := client.ContextFunctions.IsRoleInSession(ctx, NewAccountObjectIdentifier("ACCOUNTADMIN"))
	require.NoError(t, err)
	if !ok {
		t.Skip("ACCOUNTADMIN role is not in current session")
	}
	t.Run("set and unset params", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				Parameters: &AccountLevelParameters{
					AccountParameters: &AccountParameters{
						ClientEncryptionKeySize:       Int(128),
						PreventUnloadToInternalStages: Bool(true),
					},
					SessionParameters: &SessionParameters{
						JSONIndent: Int(16),
					},
					ObjectParameters: &ObjectParameters{
						MaxDataExtensionTimeInDays: Int(30),
					},
				},
			},
		}
		err := client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
		p, err := client.Sessions.ShowAccountParameter(ctx, AccountParameterClientEncryptionKeySize)
		require.NoError(t, err)
		assert.Equal(t, 128, toInt(p.Value))
		p, err = client.Sessions.ShowAccountParameter(ctx, AccountParameterPreventUnloadToInternalStages)
		require.NoError(t, err)
		assert.Equal(t, true, toBool(p.Value))
		p, err = client.Sessions.ShowAccountParameter(ctx, AccountParameterJSONIndent)
		require.NoError(t, err)
		assert.Equal(t, 16, toInt(p.Value))
		p, err = client.Sessions.ShowAccountParameter(ctx, AccountParameterMaxDataExtensionTimeInDays)
		require.NoError(t, err)
		assert.Equal(t, 30, toInt(p.Value))

		opts = &AlterAccountOptions{
			Unset: &AccountUnset{
				Parameters: &AccountLevelParametersUnset{
					AccountParameters: &AccountParametersUnset{
						ClientEncryptionKeySize:       Bool(true),
						PreventUnloadToInternalStages: Bool(true),
					},
					SessionParameters: &SessionParametersUnset{
						JSONIndent: Bool(true),
					},
					ObjectParameters: &ObjectParametersUnset{
						MaxDataExtensionTimeInDays: Bool(true),
					},
				},
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
	})

	t.Run("set resource monitor", func(t *testing.T) {
		resourceMonitorTest, resourceMonitorCleanup := createResourceMonitor(t, client)
		t.Cleanup(resourceMonitorCleanup)
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				ResourceMonitor: resourceMonitorTest.ID(),
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
	})

	t.Run("set and unset password policy", func(t *testing.T) {
		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)
		schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
		t.Cleanup(schemaCleanup)
		passwordPolicyTest, passwordPolicyCleanup := createPasswordPolicy(t, client, databaseTest, schemaTest)
		t.Cleanup(passwordPolicyCleanup)
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				PasswordPolicy: passwordPolicyTest.ID(),
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)

		// now unset
		opts = &AlterAccountOptions{
			Unset: &AccountUnset{
				PasswordPolicy: Bool(true),
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
	})

	t.Run("set and unset session policy", func(t *testing.T) {
		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)
		schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
		t.Cleanup(schemaCleanup)
		sessionPolicyTest, sessionPolicyCleanup := createSessionPolicy(t, client, databaseTest, schemaTest)
		t.Cleanup(sessionPolicyCleanup)
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				SessionPolicy: sessionPolicyTest.ID(),
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)

		// now unset
		opts = &AlterAccountOptions{
			Unset: &AccountUnset{
				SessionPolicy: Bool(true),
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
	})

	t.Run("set and unset tag", func(t *testing.T) {
		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)
		schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
		t.Cleanup(schemaCleanup)
		tagTest1, tagCleanup1 := createTag(t, client, databaseTest, schemaTest)
		t.Cleanup(tagCleanup1)
		tagTest2, tagCleanup2 := createTag(t, client, databaseTest, schemaTest)
		t.Cleanup(tagCleanup2)

		opts := &AlterAccountOptions{
			Set: &AccountSet{
				Tag: []TagAssociation{
					{
						Name:  tagTest1.ID(),
						Value: "abc",
					},
					{
						Name:  tagTest2.ID(),
						Value: "123",
					},
				},
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
		currentAccount, err := client.ContextFunctions.CurrentAccount(ctx)
		require.NoError(t, err)
		tagValue, err := client.SystemFunctions.GetTag(ctx, tagTest1.ID(), NewAccountObjectIdentifier(currentAccount), ObjectTypeAccount)
		require.NoError(t, err)
		assert.Equal(t, "abc", tagValue)
		tagValue, err = client.SystemFunctions.GetTag(ctx, tagTest2.ID(), NewAccountObjectIdentifier(currentAccount), ObjectTypeAccount)
		require.NoError(t, err)
		assert.Equal(t, "123", tagValue)
	})
}
