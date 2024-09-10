package testint

import (
	"log"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/avast/retry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_AccountShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	ok := testClientHelper().Context.IsRoleInSession(t, snowflakeroles.Orgadmin)
	if !ok {
		t.Skip("ORGADMIN role is not in current session")
	}
	currentAccount := testClientHelper().Context.CurrentAccount(t)
	opts := &sdk.ShowAccountOptions{
		Like: &sdk.Like{
			Pattern: sdk.String(currentAccount),
		},
	}
	accounts, err := client.Accounts.Show(ctx, opts)
	require.NoError(t, err)
	assert.NotEmpty(t, accounts)
	assert.Equal(t, 1, len(accounts))
	assert.Contains(t, []string{accounts[0].AccountLocator, accounts[0].AccountName}, currentAccount)
}

func TestInt_AccountShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	ok := testClientHelper().Context.IsRoleInSession(t, snowflakeroles.Orgadmin)
	if !ok {
		t.Skip("ORGADMIN role is not in current session")
	}
	_, err := client.Accounts.ShowByID(ctx, NonExistingAccountObjectIdentifier)
	require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}

func TestInt_AccountCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	ok := testClientHelper().Context.IsRoleInSession(t, snowflakeroles.Orgadmin)
	if !ok {
		t.Skip("ORGADMIN role is not in current session")
	}
	t.Run("complete case", func(t *testing.T) {
		accountID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		region := testClientHelper().Context.CurrentRegion(t)

		opts := &sdk.CreateAccountOptions{
			AdminName:          "someadmin",
			AdminPassword:      sdk.String(random.Password()),
			FirstName:          sdk.String("Ad"),
			LastName:           sdk.String("Min"),
			Email:              "admin@example.com",
			MustChangePassword: sdk.Bool(false),
			Edition:            sdk.EditionBusinessCritical,
			Comment:            sdk.String("Please delete me!"),
			Region:             sdk.String(region),
		}
		err := client.Accounts.Create(ctx, accountID, opts)
		require.NoError(t, err)

		var account *sdk.Account
		err = retry.Do(
			func() error {
				account, err = client.Accounts.ShowByID(ctx, accountID)
				return err
			},
			retry.OnRetry(func(n uint, err error) {
				log.Printf("[DEBUG] Retrying client.Accounts.ShowByID: #%d", n+1)
			}),
			retry.Delay(1*time.Second),
			retry.Attempts(3),
		)
		require.NoError(t, err)
		assert.Equal(t, accountID.Name(), account.AccountName)
		assert.Equal(t, sdk.EditionBusinessCritical, account.Edition)
		assert.Equal(t, "Please delete me!", account.Comment)
		assert.Equal(t, region, account.SnowflakeRegion)

		// rename
		newAccountID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		alterOpts := &sdk.AlterAccountOptions{
			Rename: &sdk.AccountRename{
				Name:       accountID,
				NewName:    newAccountID,
				SaveOldURL: sdk.Bool(true),
			},
		}
		err = client.Accounts.Alter(ctx, alterOpts)
		require.NoError(t, err)

		err = retry.Do(
			func() error {
				account, err = client.Accounts.ShowByID(ctx, newAccountID)
				return err
			},
			retry.OnRetry(func(n uint, err error) {
				log.Printf("[DEBUG] Retrying client.Accounts.ShowByID: #%d", n+1)
			}),
			retry.Delay(1*time.Second),
			retry.Attempts(3),
		)
		require.NoError(t, err)
		assert.Equal(t, newAccountID.Name(), account.AccountName)

		// drop old url
		alterOpts = &sdk.AlterAccountOptions{
			Drop: &sdk.AccountDrop{
				Name:   newAccountID,
				OldURL: sdk.Bool(true),
			},
		}
		err = client.Accounts.Alter(ctx, alterOpts)
		require.NoError(t, err)
		_, err = client.Accounts.ShowByID(ctx, newAccountID)
		require.NoError(t, err)

		// drop account
		err = client.Accounts.Drop(ctx, newAccountID, 3, &sdk.DropAccountOptions{
			IfExists: sdk.Bool(true),
		})
		require.NoError(t, err)

		// check if account is dropped
		_, err = client.Accounts.ShowByID(ctx, newAccountID)
		require.Error(t, err)

		// undrop account
		err = client.Accounts.Undrop(ctx, newAccountID)
		require.NoError(t, err)

		// check if account is undropped
		_, err = client.Accounts.ShowByID(ctx, newAccountID)
		require.NoError(t, err)

		// drop account again
		err = client.Accounts.Drop(ctx, newAccountID, 3, nil)
		require.NoError(t, err)

		// check if account is dropped
		_, err = client.Accounts.ShowByID(ctx, newAccountID)
		require.Error(t, err)
	})
}

func TestInt_AccountAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	ok := testClientHelper().Context.IsRoleInSession(t, snowflakeroles.Accountadmin)
	if !ok {
		t.Skip("ACCOUNTADMIN role is not in current session")
	}
	t.Run("set and unset params", func(t *testing.T) {
		opts := &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				Parameters: &sdk.AccountLevelParameters{
					AccountParameters: &sdk.AccountParameters{
						ClientEncryptionKeySize:       sdk.Int(128),
						PreventUnloadToInternalStages: sdk.Bool(true),
					},
					SessionParameters: &sdk.SessionParameters{
						JSONIndent: sdk.Int(16),
					},
					ObjectParameters: &sdk.ObjectParameters{
						MaxDataExtensionTimeInDays: sdk.Int(30),
					},
				},
			},
		}
		err := client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
		p, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterClientEncryptionKeySize)
		require.NoError(t, err)
		assert.Equal(t, 128, sdk.ToInt(p.Value))
		p, err = client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterPreventUnloadToInternalStages)
		require.NoError(t, err)
		assert.Equal(t, true, sdk.ToBool(p.Value))
		p, err = client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterJSONIndent)
		require.NoError(t, err)
		assert.Equal(t, 16, sdk.ToInt(p.Value))
		p, err = client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterMaxDataExtensionTimeInDays)
		require.NoError(t, err)
		assert.Equal(t, 30, sdk.ToInt(p.Value))

		opts = &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				Parameters: &sdk.AccountLevelParametersUnset{
					AccountParameters: &sdk.AccountParametersUnset{
						ClientEncryptionKeySize:       sdk.Bool(true),
						PreventUnloadToInternalStages: sdk.Bool(true),
					},
					SessionParameters: &sdk.SessionParametersUnset{
						JSONIndent: sdk.Bool(true),
					},
					ObjectParameters: &sdk.ObjectParametersUnset{
						MaxDataExtensionTimeInDays: sdk.Bool(true),
					},
				},
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
	})

	t.Run("set resource monitor", func(t *testing.T) {
		resourceMonitorTest, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)
		opts := &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				ResourceMonitor: resourceMonitorTest.ID(),
			},
		}
		err := client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
	})

	t.Run("set and unset password policy", func(t *testing.T) {
		passwordPolicyTest, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)
		opts := &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				PasswordPolicy: passwordPolicyTest.ID(),
			},
		}
		err := client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)

		// now unset
		opts = &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				PasswordPolicy: sdk.Bool(true),
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
	})

	t.Run("set and unset session policy", func(t *testing.T) {
		sessionPolicyTest, sessionPolicyCleanup := testClientHelper().SessionPolicy.CreateSessionPolicy(t)
		t.Cleanup(sessionPolicyCleanup)

		opts := &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				SessionPolicy: sessionPolicyTest.ID(),
			},
		}
		err := client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)

		// now unset
		opts = &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				SessionPolicy: sdk.Bool(true),
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
	})

	t.Run("set and unset authentication policy", func(t *testing.T) {
		t.Skipf("Skipping the test for now TODO: add ticket number")
		authenticationPolicyTest, authenticationPolicyCleanup := testClientHelper().AuthenticationPolicy.CreateAuthenticationPolicy(t)
		t.Cleanup(authenticationPolicyCleanup)
		opts := &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				AuthenticationPolicy: authenticationPolicyTest.ID(),
			},
		}
		err := client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)

		// now unset
		opts = &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				AuthenticationPolicy: sdk.Bool(true),
			},
		}
		err = client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
	})

	t.Run("set and unset tag", func(t *testing.T) {
		tagTest1, tagCleanup1 := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup1)
		tagTest2, tagCleanup2 := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup2)

		opts := &sdk.AlterAccountOptions{
			SetTag: []sdk.TagAssociation{
				{
					Name:  tagTest1.ID(),
					Value: "abc",
				},
				{
					Name:  tagTest2.ID(),
					Value: "123",
				},
			},
		}
		err := client.Accounts.Alter(ctx, opts)
		require.NoError(t, err)
		tagValue, err := client.SystemFunctions.GetTag(ctx, tagTest1.ID(), testClientHelper().Ids.AccountIdentifierWithLocator(), sdk.ObjectTypeAccount)
		require.NoError(t, err)
		assert.Equal(t, "abc", tagValue)
		tagValue, err = client.SystemFunctions.GetTag(ctx, tagTest2.ID(), testClientHelper().Ids.AccountIdentifierWithLocator(), sdk.ObjectTypeAccount)
		require.NoError(t, err)
		assert.Equal(t, "123", tagValue)
	})
}
