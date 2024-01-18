package testint

import (
	"os"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-SNOW-1011985]: unskip the tests
func TestInt_ManagedAccounts(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_MANAGED_ACCOUNT_INT_TEST"); ok {
		t.Skip("Skipping TestInt_ManagedAccounts due to error: 090337 (23001): Number of managed accounts allowed exceeded the limit. Please contact Snowflake support.")
	}

	client := testClient(t)
	ctx := testContext(t)

	assertManagedAccount := func(t *testing.T, managedAccount *sdk.ManagedAccount, id sdk.AccountObjectIdentifier, comment string) {
		t.Helper()
		assert.Equal(t, id.Name(), managedAccount.Name)
		assert.Equal(t, "aws", managedAccount.Cloud)
		assert.NotEmpty(t, managedAccount.Region)
		assert.NotEmpty(t, managedAccount.Locator)
		assert.NotEmpty(t, managedAccount.CreatedOn)
		assert.NotEmpty(t, managedAccount.URL)
		assert.NotEmpty(t, managedAccount.AccountLocatorURL)
		assert.Equal(t, true, managedAccount.IsReader)
		assert.Equal(t, comment, managedAccount.Comment)
	}

	cleanupMangedAccountProvider := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.ManagedAccounts.Drop(ctx, sdk.NewDropManagedAccountRequest(id))
			require.NoError(t, err)
		}
	}

	createManagedAccountBasicRequest := func(t *testing.T) *sdk.CreateManagedAccountRequest {
		t.Helper()
		// 090348 (42602): Account name or alias is invalid: (...) can only contain capital letters, numbers, and underscores
		name := strings.ToUpper(random.AlphanumericN(10))
		id := sdk.NewAccountObjectIdentifier(name)

		// 090088 (22000): ADMIN_NAME can only contain letters, numbers and underscores.
		// 090089 (22000): ADMIN_NAME must start with a letter.
		adminName := random.AlphaN(1) + random.AlphanumericN(12)
		adminPassword := random.String()
		params := sdk.NewCreateManagedAccountParamsRequest(adminName, adminPassword)

		return sdk.NewCreateManagedAccountRequest(id, *params)
	}

	createManagedAccountWithRequest := func(t *testing.T, request *sdk.CreateManagedAccountRequest) *sdk.ManagedAccount {
		t.Helper()
		id := request.GetName()

		err := client.ManagedAccounts.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupMangedAccountProvider(id))

		managedAccount, err := client.ManagedAccounts.ShowByID(ctx, id)
		require.NoError(t, err)

		return managedAccount
	}

	createManagedAccount := func(t *testing.T) *sdk.ManagedAccount {
		t.Helper()
		return createManagedAccountWithRequest(t, createManagedAccountBasicRequest(t))
	}

	t.Run("create managed account: no optionals", func(t *testing.T) {
		request := createManagedAccountBasicRequest(t)

		managedAccount := createManagedAccountWithRequest(t, request)

		assertManagedAccount(t, managedAccount, request.GetName(), "")
	})

	t.Run("create managed account: full", func(t *testing.T) {
		request := createManagedAccountBasicRequest(t)
		request.CreateManagedAccountParams.Comment = sdk.String("some comment")

		managedAccount := createManagedAccountWithRequest(t, request)

		assertManagedAccount(t, managedAccount, request.GetName(), "some comment")
	})

	t.Run("drop managed account: existing", func(t *testing.T) {
		request := createManagedAccountBasicRequest(t)
		id := request.GetName()

		err := client.ManagedAccounts.Create(ctx, request)
		require.NoError(t, err)

		err = client.ManagedAccounts.Drop(ctx, sdk.NewDropManagedAccountRequest(id))
		require.NoError(t, err)

		_, err = client.ManagedAccounts.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop managed account: non-existing", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier("does_not_exist")

		err := client.ManagedAccounts.Drop(ctx, sdk.NewDropManagedAccountRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show managed account: default", func(t *testing.T) {
		managedAccount1 := createManagedAccount(t)
		managedAccount2 := createManagedAccount(t)

		showRequest := sdk.NewShowManagedAccountRequest()
		returnedManagedAccounts, err := client.ManagedAccounts.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedManagedAccounts))
		assert.Contains(t, returnedManagedAccounts, *managedAccount1)
		assert.Contains(t, returnedManagedAccounts, *managedAccount2)
	})

	t.Run("show managed account: with like", func(t *testing.T) {
		managedAccount1 := createManagedAccount(t)
		managedAccount2 := createManagedAccount(t)

		showRequest := sdk.NewShowManagedAccountRequest().
			WithLike(&sdk.Like{Pattern: &managedAccount1.Name})
		returnedManagedAccounts, err := client.ManagedAccounts.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedManagedAccounts))
		assert.Contains(t, returnedManagedAccounts, *managedAccount1)
		assert.NotContains(t, returnedManagedAccounts, *managedAccount2)
	})
}
