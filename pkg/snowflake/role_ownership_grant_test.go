package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestRoleOwnershipGrantQuery(t *testing.T) {
	r := require.New(t)
	copyBuilder := snowflake.NewRoleOwnershipGrantBuilder("role1", "COPY")
	revokeBuilder := snowflake.NewRoleOwnershipGrantBuilder("role1", "REVOKE")

	g1 := copyBuilder.Role("role2").Grant()
	r.Equal(`GRANT OWNERSHIP ON ROLE "role1" TO ROLE "role2" COPY CURRENT GRANTS`, g1)

	r1 := copyBuilder.Role("ACCOUNTADMIN").Revoke()
	r.Equal(`GRANT OWNERSHIP ON ROLE "role1" TO ROLE "ACCOUNTADMIN" COPY CURRENT GRANTS`, r1)

	g2 := revokeBuilder.Role("role2").Grant()
	r.Equal(`GRANT OWNERSHIP ON ROLE "role1" TO ROLE "role2" REVOKE CURRENT GRANTS`, g2)

	r2 := revokeBuilder.Role("ACCOUNTADMIN").Revoke()
	r.Equal(`GRANT OWNERSHIP ON ROLE "role1" TO ROLE "ACCOUNTADMIN" REVOKE CURRENT GRANTS`, r2)
}
