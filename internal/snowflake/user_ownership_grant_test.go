// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/stretchr/testify/require"
)

func TestUserOwnershipGrantQuery(t *testing.T) {
	r := require.New(t)
	copyBuilder := snowflake.NewUserOwnershipGrantBuilder("user1", "COPY")
	revokeBuilder := snowflake.NewUserOwnershipGrantBuilder("user1", "REVOKE")

	g1 := copyBuilder.Role("role1").Grant()
	r.Equal(`GRANT OWNERSHIP ON USER "user1" TO ROLE "role1" COPY CURRENT GRANTS`, g1)

	r1 := copyBuilder.Role("ACCOUNTADMIN").Revoke()
	r.Equal(`GRANT OWNERSHIP ON USER "user1" TO ROLE "ACCOUNTADMIN" COPY CURRENT GRANTS`, r1)

	g2 := revokeBuilder.Role("role1").Grant()
	r.Equal(`GRANT OWNERSHIP ON USER "user1" TO ROLE "role1" REVOKE CURRENT GRANTS`, g2)

	r2 := revokeBuilder.Role("ACCOUNTADMIN").Revoke()
	r.Equal(`GRANT OWNERSHIP ON USER "user1" TO ROLE "ACCOUNTADMIN" REVOKE CURRENT GRANTS`, r2)
}
