package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestRoleOwnershipGrantQuery(t *testing.T) {
	r := require.New(t)
	copy := snowflake.RoleOwnershipGrant("role1", "COPY")
	revoke := snowflake.RoleOwnershipGrant("role1", "REVOKE")

	g1 := copy.Role("role2").Grant()
	r.Equal(`GRANT OWNERSHIP ON ROLE "role1" TO ROLE "role2" COPY CURRENT GRANTS`, g1)

	r1 := copy.Role("ACCOUNTADMIN").Revoke()
	r.Equal(`GRANT OWNERSHIP ON ROLE "role1" TO ROLE "ACCOUNTADMIN" COPY CURRENT GRANTS`, r1)

	g2 := revoke.Role("role2").Grant()
	r.Equal(`GRANT OWNERSHIP ON ROLE "role1" TO ROLE "role2" REVOKE CURRENT GRANTS`, g2)

	r2 := revoke.Role("ACCOUNTADMIN").Revoke()
	r.Equal(`GRANT OWNERSHIP ON ROLE "role1" TO ROLE "ACCOUNTADMIN" REVOKE CURRENT GRANTS`, r2)
}
