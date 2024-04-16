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

	config, err := ProfileConfig(testprofiles.Secondary)
	if err != nil {
		t.Skipf("Snowflake secondary account not configured. Must be set in ~./snowflake/config.yml with profile name: %s", testprofiles.Secondary)
	}
	client, err := NewClient(config)
	if err != nil {
		t.Skipf("Snowflake secondary account not configured. Must be set in ~./snowflake/config.yml with profile name: %s", testprofiles.Secondary)
	}

	return client
}
