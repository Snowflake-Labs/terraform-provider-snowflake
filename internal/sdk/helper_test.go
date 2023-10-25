// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import (
	"testing"
)

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
