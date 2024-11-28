package testint

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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

// TODO: Change AccountObjectIdentifier to AccountIdentifier
// TODO: See if wait time is needed after create before calling next commands (so far it's not needed)

func TestInt_Account(t *testing.T) {
	if !testClientHelper().Context.IsRoleInSession(t, snowflakeroles.Orgadmin) {
		t.Skip("ORGADMIN role is not in current session")
	}
	client := testClient(t)
	ctx := testContext(t)
	currentAccountName := testClientHelper().Context.CurrentAccountName(t)

	assertAccount := func(t *testing.T, account sdk.Account, accountName string) {
		t.Helper()
		assert.NotEmpty(t, account.OrganizationName)
		assert.Equal(t, accountName, account.AccountName)
		assert.Nil(t, account.RegionGroup) // TODO: Removed?
		assert.NotEmpty(t, account.SnowflakeRegion)
		assert.Equal(t, sdk.EditionEnterprise, account.Edition)
		assert.NotEmpty(t, account.AccountURL)
		assert.NotEmpty(t, account.CreatedOn)
		assert.Equal(t, "SNOWFLAKE", account.Comment)
		assert.NotEmpty(t, account.AccountLocator)
		assert.NotEmpty(t, account.AccountLocatorURL)
		assert.Zero(t, account.ManagedAccounts)
		assert.NotEmpty(t, account.ConsumptionBillingEntityName)
		assert.Nil(t, account.MarketplaceConsumerBillingEntityName)
		assert.NotNil(t, account.MarketplaceProviderBillingEntityName)
		assert.Empty(t, account.OldAccountURL)
		//assert.False(t, account.IsOrgAdmin)
		assert.True(t, account.IsOrgAdmin)
		assert.Nil(t, account.AccountOldUrlSavedOn)
		assert.Nil(t, account.AccountOldUrlLastUsed)
		assert.Empty(t, account.OrganizationOldUrl)
		assert.Nil(t, account.OrganizationOldUrlSavedOn)
		assert.Nil(t, account.OrganizationOldUrlLastUsed)
		assert.False(t, account.IsEventsAccount)
		assert.False(t, account.IsOrganizationAccount)
	}

	assertHistoryAccount := func(t *testing.T, account sdk.Account, accountName string) {
		t.Helper()
		assertAccount(t, account, currentAccountName)
		assert.Nil(t, account.DroppedOn)
		assert.Nil(t, account.ScheduledDeletionTime)
		assert.Nil(t, account.RestoredOn)
		assert.Empty(t, account.MovedToOrganization)
		assert.Nil(t, account.MovedOn)
		assert.Nil(t, account.OrganizationUrlExpirationOn)
	}

	// TODO: Uncomment and use if needed; otherwise remove
	// awaitAndAssertAccountCreation := func(t *testing.T, id sdk.AccountObjectIdentifier) {
	//	require.Eventually(t, func() bool {
	//		account, err := client.Accounts.ShowByID(ctx, id)
	//		return err == nil && account.ID() == id
	//	}, 45*time.Second, time.Second)
	// }

	t.Run("create: minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		name := testClientHelper().Ids.Alpha()
		password := random.Password()
		email := random.Email()

		err := client.Accounts.Create(ctx, id, &sdk.CreateAccountOptions{
			AdminName:     name,
			AdminPassword: sdk.String(password),
			Email:         email,
			Edition:       sdk.EditionStandard,
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
	})

	t.Run("create: user type service", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		name := testClientHelper().Ids.Alpha()
		key, _ := random.GenerateRSAPublicKey(t)
		email := random.Email()

		err := client.Accounts.Create(ctx, id, &sdk.CreateAccountOptions{
			AdminName:         name,
			AdminRSAPublicKey: sdk.String(key),
			AdminUserType:     sdk.Pointer(sdk.UserTypeService),
			Email:             email,
			Edition:           sdk.EditionStandard,
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
	})

	t.Run("create: user type legacy service", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		name := testClientHelper().Ids.Alpha()
		password := random.Password()
		email := random.Email()

		err := client.Accounts.Create(ctx, id, &sdk.CreateAccountOptions{
			AdminName:     name,
			AdminPassword: sdk.String(password),
			AdminUserType: sdk.Pointer(sdk.UserTypeLegacyService),
			Email:         email,
			Edition:       sdk.EditionStandard,
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		name := testClientHelper().Ids.Alpha()
		password := random.Password()
		email := random.Email()
		region := testClientHelper().Context.CurrentRegion(t)
		regions := testClientHelper().Account.ShowRegions(t)
		currentRegion, err := collections.FindFirst(regions, func(r helpers.Region) bool {
			return r.SnowflakeRegion == region
		})
		require.NoError(t, err)
		comment := random.Comment()

		err = client.Accounts.Create(ctx, id, &sdk.CreateAccountOptions{
			AdminName:          name,
			AdminPassword:      sdk.String(password),
			FirstName:          sdk.String("firstName"),
			LastName:           sdk.String("lastName"),
			Email:              email,
			MustChangePassword: sdk.Bool(true),
			Edition:            sdk.EditionStandard,
			// TODO: passed region group cannot be recognised by Snowflake
			//RegionGroup: sdk.String(currentRegion.SnowflakeRegion),
			Region:  sdk.String(currentRegion.SnowflakeRegion),
			Comment: sdk.String(comment),
			// TODO(TODO: ticket): with polaris Snowflake returns an error saying: "invalid property polaris for account"
			// Polaris: sdk.Bool(true),
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
	})

	// create a connection to the newly created account and test self-alter?
	// t.Run("self alter: set / unset", func(t *testing.T) {})

	t.Run("alter: set / unset is org admin", func(t *testing.T) {
		account := testClientHelper().Account.Create(t)
		t.Cleanup(testClientHelper().Account.DropFunc(t, account.ID()))

		require.Equal(t, false, account.IsOrgAdmin)

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
				Name:     account.ID(),
				OrgAdmin: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, true, acc.IsOrgAdmin)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
				Name:     account.ID(),
				OrgAdmin: sdk.Bool(false),
			},
		})
		require.NoError(t, err)

		acc, err = client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, false, acc.IsOrgAdmin)
	})

	t.Run("alter: rename", func(t *testing.T) {
		account := testClientHelper().Account.Create(t)
		t.Cleanup(testClientHelper().Account.DropFunc(t, account.ID()))

		newName := testClientHelper().Ids.RandomAccountObjectIdentifier()
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Rename: &sdk.AccountRename{
				Name:    account.ID(),
				NewName: newName,
			},
		})
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, account.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		acc, err := client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.NotNil(t, acc)
		require.NotEqual(t, account.AccountURL, acc.AccountURL)
		require.Equal(t, account.AccountURL, acc.OldAccountURL)
	})

	t.Run("alter: rename with new url", func(t *testing.T) {
		account := testClientHelper().Account.Create(t)
		t.Cleanup(testClientHelper().Account.DropFunc(t, account.ID()))

		newName := testClientHelper().Ids.RandomAccountObjectIdentifier()
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Rename: &sdk.AccountRename{
				Name:       account.ID(),
				NewName:    newName,
				SaveOldURL: sdk.Bool(false),
			},
		})
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, account.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		acc, err := client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.NotEqual(t, account.AccountURL, acc.AccountURL)
		require.Empty(t, acc.OldAccountURL)
	})

	t.Run("alter: drop url when there's no old url", func(t *testing.T) {
		account := testClientHelper().Account.Create(t)
		t.Cleanup(testClientHelper().Account.DropFunc(t, account.ID()))

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Drop: &sdk.AccountDrop{
				Name:   account.ID(),
				OldUrl: sdk.Bool(true),
			},
		})
		require.ErrorContains(t, err, "The account has no old url")
	})

	t.Run("alter: drop url after rename", func(t *testing.T) {
		account := testClientHelper().Account.Create(t)
		t.Cleanup(testClientHelper().Account.DropFunc(t, account.ID()))

		newName := testClientHelper().Ids.RandomAccountObjectIdentifier()
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Rename: &sdk.AccountRename{
				Name:    account.ID(),
				NewName: newName,
			},
		})
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.NotEmpty(t, acc.OldAccountURL)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Drop: &sdk.AccountDrop{
				Name:   newName,
				OldUrl: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		acc, err = client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.Empty(t, acc.OldAccountURL)
	})

	// TODO(TODO: Ticket): This cannot be tested as it requires capabilities of moving accounts between organizations.
	// From the documentation: https://docs.snowflake.com/en/sql-reference/sql/show-accounts#output
	// `
	// If the account’s organization was changed in a way that created a new account URL and the original account URL was saved,
	// provides the original account URL. If the original account URL was dropped, the value is NULL even if the organization changed.
	// `
	// t.Run("alter: drop organization url after rename", func(t *testing.T) {
	//	account := testClientHelper().Account.Create(t)
	//	t.Cleanup(testClientHelper().Account.DropFunc(t, account.ID()))
	//
	//	newName := testClientHelper().Ids.RandomAccountObjectIdentifier()
	//	t.Cleanup(testClientHelper().Account.DropFunc(t, newName))
	//
	//  // move the account to another organization
	//
	//	acc, err := client.Accounts.ShowByID(ctx, newName)
	//	require.NoError(t, err)
	//	require.NotEmpty(t, acc.OrganizationOldUrl)
	//
	//	err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
	//		Drop: &sdk.AccountDrop{
	//			Name:               newName,
	//			OldOrganizationUrl: sdk.Bool(true),
	//		},
	//	})
	//	require.NoError(t, err)
	//
	//	acc, err = client.Accounts.ShowByID(ctx, newName)
	//	require.NoError(t, err)
	//	require.Empty(t, acc.OrganizationOldUrl)
	// })

	t.Run("drop: without options", func(t *testing.T) {
		err := client.Accounts.Drop(ctx, sdk.NewAccountObjectIdentifier("non-existing-account"), 3, &sdk.DropAccountOptions{})
		require.Error(t, err)

		account := testClientHelper().Account.Create(t)

		err = client.Accounts.Drop(ctx, account.ID(), 3, &sdk.DropAccountOptions{})
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, account.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: with if exists", func(t *testing.T) {
		err := client.Accounts.Drop(ctx, sdk.NewAccountObjectIdentifier("non-existing-account"), 3, &sdk.DropAccountOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)

		account := testClientHelper().Account.Create(t)

		err = client.Accounts.Drop(ctx, account.ID(), 3, &sdk.DropAccountOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, account.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("undrop", func(t *testing.T) {
		account := testClientHelper().Account.Create(t)
		t.Cleanup(testClientHelper().Account.DropFunc(t, account.ID()))

		require.NoError(t, testClientHelper().Account.Drop(t, account.ID()))

		err := client.Accounts.Undrop(ctx, account.ID())
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, account.ID(), acc.ID())
	})

	t.Run("show: with like", func(t *testing.T) {
		currentAccount := testClientHelper().Context.CurrentAccount(t)
		opts := &sdk.ShowAccountOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(currentAccount),
			},
		}

		accounts, err := client.Accounts.Show(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, 1, len(accounts))
		assertAccount(t, accounts[0], currentAccountName)
	})

	t.Run("show: with history", func(t *testing.T) {
		currentAccount := testClientHelper().Context.CurrentAccount(t)
		opts := &sdk.ShowAccountOptions{
			History: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(currentAccount),
			},
		}

		accounts, err := client.Accounts.Show(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, 1, len(accounts))
		assertHistoryAccount(t, accounts[0], currentAccountName)
	})
}

func TestInt_AccountCreate(t *testing.T) {
	if !testClientHelper().Context.IsRoleInSession(t, snowflakeroles.Orgadmin) {
		t.Skip("ORGADMIN role is not in current session")
	}
	client := testClient(t)
	ctx := testContext(t)

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
				OldUrl: sdk.Bool(true),
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
	if !testClientHelper().Context.IsRoleInSession(t, snowflakeroles.Orgadmin) {
		t.Skip("ORGADMIN role is not in current session")
	}
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
		authenticationPolicyTest, authenticationPolicyCleanup := testClientHelper().AuthenticationPolicy.Create(t)
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
}
