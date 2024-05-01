package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type AccountClient struct {
	context *TestClientContext
}

func NewAccountClient(context *TestClientContext) *AccountClient {
	return &AccountClient{
		context: context,
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
