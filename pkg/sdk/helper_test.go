package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// there is no direct way to get the account identifier from Snowflake API, but you can get it if you know
// the account locator and by filtering the list of accounts in replication accounts by the account locator
func getAccountIdentifier(t *testing.T, client *Client) AccountIdentifier {
	t.Helper()
	ctx := context.Background()
	currentAccountLocator, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	replicationAccounts, err := client.ReplicationFunctions.ShowReplicationAccounts(ctx)
	require.NoError(t, err)
	for _, replicationAccount := range replicationAccounts {
		if replicationAccount.AccountLocator == currentAccountLocator {
			return AccountIdentifier{
				organizationName: replicationAccount.OrganizationName,
				accountName:      replicationAccount.AccountName,
			}
		}
	}
	return AccountIdentifier{}
}

func getSecondaryAccountIdentifier(t *testing.T) AccountIdentifier {
	t.Helper()
	client := testSecondaryClient(t)
	return getAccountIdentifier(t, client)
}

func testClient(t *testing.T) *Client {
	t.Helper()

	client, err := NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}

const (
	secondaryAccountProfile = "secondary_test_account"
)

func testSecondaryClient(t *testing.T) *Client {
	t.Helper()

	client, err := testClientFromProfile(t, secondaryAccountProfile)
	if err != nil {
		t.Skipf("Snowflake secondary account not configured. Must be set in ~./snowflake/config.yml with profile name: %s", secondaryAccountProfile)
	}

	return client
}

func testClientFromProfile(t *testing.T, profile string) (*Client, error) {
	t.Helper()
	config, err := ProfileConfig(profile)
	if err != nil {
		return nil, err
	}
	return NewClient(config)
}
