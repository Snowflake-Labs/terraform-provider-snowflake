// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/stretchr/testify/require"
)

func TestNetworkPolicySetOnAccount(t *testing.T) {
	r := require.New(t)
	s := snowflake.NetworkPolicy("test_network_policy")
	r.NotNil(s)

	q := s.SetOnAccount()
	r.Equal(`ALTER ACCOUNT SET NETWORK_POLICY = "test_network_policy"`, q)
}

func TestNetworkPolicyUnsetOnAccount(t *testing.T) {
	r := require.New(t)
	s := snowflake.NetworkPolicy("test_network_policy")
	r.NotNil(s)

	q := s.UnsetOnAccount()
	r.Equal(`ALTER ACCOUNT UNSET NETWORK_POLICY`, q)
}

func TestNetworkPolicySetOnUser(t *testing.T) {
	r := require.New(t)
	s := snowflake.NetworkPolicy("test_network_policy")
	r.NotNil(s)

	q := s.SetOnUser("testuser")
	r.Equal(`ALTER USER "testuser" SET NETWORK_POLICY = "test_network_policy"`, q)
}

func TestNetworkPolicyUnsetOnUser(t *testing.T) {
	r := require.New(t)
	s := snowflake.NetworkPolicy("test_network_policy")
	r.NotNil(s)

	q := s.UnsetOnUser("testuser")
	r.Equal(`ALTER USER "testuser" UNSET NETWORK_POLICY`, q)
}

func TestNetworkPolicyShowOnUser(t *testing.T) {
	r := require.New(t)
	s := snowflake.NetworkPolicy("test_network_policy")
	r.NotNil(s)

	q := s.ShowOnUser("testuser")
	r.Equal(`SHOW PARAMETERS LIKE 'network_policy' IN USER "testuser"`, q)
}

func TestNetworkPolicyShowOnAccount(t *testing.T) {
	r := require.New(t)
	s := snowflake.NetworkPolicy("test_network_policy")
	r.NotNil(s)

	q := s.ShowOnAccount()
	r.Equal(`SHOW PARAMETERS LIKE 'network_policy' IN ACCOUNT`, q)
}
