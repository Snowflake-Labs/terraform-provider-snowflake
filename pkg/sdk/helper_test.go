package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
)

func defaultTestClient(t *testing.T) *Client {
	t.Helper()

	client, err := NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func testSecondaryClient(t *testing.T) *Client {
	t.Helper()
	return testClient(t, testprofiles.Secondary)
}

func testThirdClient(t *testing.T) *Client {
	t.Helper()
	return testClient(t, testprofiles.Third)
}

func testFourthClient(t *testing.T) *Client {
	t.Helper()
	return testClient(t, testprofiles.Fourth)
}

func testClient(t *testing.T, profile string) *Client {
	t.Helper()

	config, err := ProfileConfig(profile)
	if err != nil {
		t.Skipf("Snowflake %s profile not configured. Must be set in ~./snowflake/config.yml", profile)
	}
	client, err := NewClient(config)
	if err != nil {
		t.Skipf("Snowflake %s profile not configured. Must be set in ~./snowflake/config.yml", profile)
	}

	return client
}
