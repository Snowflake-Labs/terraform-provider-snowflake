package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
)

func testClient(t *testing.T) *Client {
	t.Helper()

	client, err := NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func testSecondaryClient(t *testing.T) *Client {
	t.Helper()

	client, err := testClientFromProfile(t, testprofiles.Secondary)
	if err != nil {
		t.Skipf("Snowflake secondary account not configured. Must be set in ~./snowflake/config.yml with profile name: %s", testprofiles.Secondary)
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
