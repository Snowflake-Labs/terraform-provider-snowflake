// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/stretchr/testify/require"
)

func TestManagedAccount(t *testing.T) {
	r := require.New(t)
	u := snowflake.NewManagedAccountBuilder("managedaccount1")
	r.NotNil(u)

	q := u.Show()
	r.Equal("SHOW MANAGED ACCOUNTS LIKE 'managedaccount1'", q)

	q = u.Drop()
	r.Equal(`DROP MANAGED ACCOUNT "managedaccount1"`, q)

	c := u.Create()
	c.SetString("foo", "bar")
	c.SetBool("bam", false)
	q = c.Statement()
	r.Equal(`CREATE MANAGED ACCOUNT "managedaccount1" FOO='bar' BAM=false`, q)
}
