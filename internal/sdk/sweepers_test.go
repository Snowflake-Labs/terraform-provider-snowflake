// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSweepAll(t *testing.T) {
	enableSweep := os.Getenv("SNOWFLAKE_ENABLE_SWEEP")
	if enableSweep != "1" {
		t.Skip("SNOWFLAKE_ENABLE_SWEEP not enabled, skipping sweep tests")
	}

	t.Run("all sweepers in secondary account", func(t *testing.T) {
		client := testSecondaryClient(t)
		err := SweepAll(client)
		require.NoError(t, err)
	})

	t.Run("all sweepers in primary account", func(t *testing.T) {
		client := testClient(t)
		err := SweepAll(client)
		require.NoError(t, err)
	})
}

func TestSweep(t *testing.T) {
	enableSweep := os.Getenv("SNOWFLAKE_ENABLE_SWEEP")
	if enableSweep != "1" {
		t.Skip("SNOWFLAKE_ENABLE_SWEEP not enabled, skipping sweep tests")
	}
	t.Run("sweepers", func(t *testing.T) {
		client := testClient(t)
		err := Sweep(client, "TEST_")
		require.NoError(t, err)
	})
}
