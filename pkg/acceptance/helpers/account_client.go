package helpers

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
	"testing"
)

type AccountClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewAccountClient(context *TestClientContext, idsGenerator *IdsGenerator) *AccountClient {
	return &AccountClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *AccountClient) client() sdk.Accounts {
	return c.context.client.Accounts
}

// GetAccountIdentifier gets the account identifier from Snowflake API, by fetching the account locator
// and by filtering the list of accounts in replication accounts by it (because there is no direct way to get).
func (c *AccountClient) GetAccountIdentifier(t *testing.T) sdk.AccountIdentifier {
	t.Helper()
	ctx := context.Background()

	currentAccountLocator, err := c.context.client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)

	replicationAccounts, err := c.context.client.ReplicationFunctions.ShowReplicationAccounts(ctx)
	require.NoError(t, err)

	for _, replicationAccount := range replicationAccounts {
		if replicationAccount.AccountLocator == currentAccountLocator {
			return sdk.NewAccountIdentifier(replicationAccount.OrganizationName, replicationAccount.AccountName)
		}
	}
	t.Fatal("could not find the account identifier for the locator")
	return sdk.AccountIdentifier{}
}

func (c *AccountClient) Create(t *testing.T) *sdk.Account {
	id := c.ids.RandomAccountObjectIdentifier()
	name := c.ids.Alpha()
	password := random.Password()
	email := random.Email()

	err := c.client().Create(context.Background(), id, &sdk.CreateAccountOptions{
		AdminName:     name,
		AdminPassword: sdk.String(password),
		Email:         email,
		Edition:       sdk.EditionStandard,
	})
	require.NoError(t, err)
	t.Cleanup(c.DropFunc(t, id))

	account, err := c.client().ShowByID(context.Background(), id)
	require.NoError(t, err)

	return account
}

func (c *AccountClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	return func() {
		require.NoError(t, c.Drop(t, id))
	}
}

func (c *AccountClient) Drop(t *testing.T, id sdk.AccountObjectIdentifier) error {
	t.Helper()
	ctx := context.Background()

	if err := c.client().Drop(ctx, id, 3, &sdk.DropAccountOptions{IfExists: sdk.Bool(true)}); err != nil {
		return err
	}
	return nil
}

type Region struct {
	SnowflakeRegion string
	Cloud           string
	Region          string
	DisplayName     string
}

func (c *AccountClient) ShowRegions(t *testing.T) []Region {
	t.Helper()

	results, err := c.context.client.QueryUnsafe(context.Background(), "SHOW REGIONS")
	require.NoError(t, err)

	return collections.Map(results, func(result map[string]*any) Region {
		require.NotNil(t, result["snowflake_region"])
		require.NotEmpty(t, *result["snowflake_region"])

		require.NotNil(t, result["cloud"])
		require.NotEmpty(t, *result["cloud"])

		require.NotNil(t, result["region"])
		require.NotEmpty(t, *result["region"])

		require.NotNil(t, result["display_name"])
		require.NotEmpty(t, *result["display_name"])

		return Region{
			SnowflakeRegion: (*result["snowflake_region"]).(string),
			Cloud:           (*result["cloud"]).(string),
			Region:          (*result["region"]).(string),
			DisplayName:     (*result["display_name"]).(string),
		}
	})
}
