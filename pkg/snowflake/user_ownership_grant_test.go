package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestUserOwnershipGrantQuery(t *testing.T) {
	r := require.New(t)
	copy := snowflake.UserOwnershipGrant("user1", "COPY")
	revoke := snowflake.UserOwnershipGrant("user1", "REVOKE")

	g1 := copy.Role("role1").Grant()
	r.Equal(`GRANT OWNERSHIP ON USER "user1" TO ROLE "role1" COPY CURRENT GRANTS`, g1)

	r1 := copy.Role("ACCOUNTADMIN").Revoke()
	r.Equal(`GRANT OWNERSHIP ON USER "user1" TO ROLE "ACCOUNTADMIN" COPY CURRENT GRANTS`, r1)

	g2 := revoke.Role("role1").Grant()
	r.Equal(`GRANT OWNERSHIP ON USER "user1" TO ROLE "role1" REVOKE CURRENT GRANTS`, g2)

	r2 := revoke.Role("ACCOUNTADMIN").Revoke()
	r.Equal(`GRANT OWNERSHIP ON USER "user1" TO ROLE "ACCOUNTADMIN" REVOKE CURRENT GRANTS`, r2)
}
