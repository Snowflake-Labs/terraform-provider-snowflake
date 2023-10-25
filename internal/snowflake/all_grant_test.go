// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/stretchr/testify/require"
)

func TestExistingSchemaGrant(t *testing.T) {
	r := require.New(t)
	builder := snowflake.AllSchemaGrant("test_db")

	r.Equal("test_db", builder.Name())

	r.Equal(string(snowflake.AllGrantTypeSchema), builder.GrantType())

	eb := builder.Role("bob")
	r.Equal("SHOW GRANTS ON DATABASE \"test_db\"", eb.Show())
	r.Equal(`GRANT USAGE ON ALL SCHEMAS IN DATABASE "test_db" TO ROLE "bob"`, eb.Grant("USAGE", false))
	r.Equal([]string{`REVOKE USAGE ON ALL SCHEMAS IN DATABASE "test_db" FROM ROLE "bob"`}, eb.Revoke("USAGE"))
}
