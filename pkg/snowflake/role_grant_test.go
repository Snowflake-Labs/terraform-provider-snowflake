package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestRoleGrant(t *testing.T) {
	r := require.New(t)
	rg := snowflake.RoleGrant("role1")

	u := rg.User("user1").Grant()
	r.Equal(`GRANT ROLE "role1" TO USER "user1"`, u)

	role := rg.Role("role2").Grant()
	r.Equal(`GRANT ROLE "role1" TO ROLE "role2"`, role)

	u2 := rg.User("user1").Revoke()
	r.Equal(`REVOKE ROLE "role1" FROM USER "user1"`, u2)

	r2 := rg.Role("role2").Revoke()
	r.Equal(`REVOKE ROLE "role1" FROM ROLE "role2"`, r2)

}
