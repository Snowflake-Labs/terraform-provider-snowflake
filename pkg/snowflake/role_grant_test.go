package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestRoleGrant(t *testing.T) {
	a := require.New(t)
	rg := snowflake.RoleGrant("role1")

	u := rg.User("user1").Grant()
	a.Equal(`GRANT ROLE "role1" TO USER "user1"`, u)

	r := rg.Role("role2").Grant()
	a.Equal(`GRANT ROLE "role1" TO ROLE "role2"`, r)

	u2 := rg.User("user1").Revoke()
	a.Equal(`REVOKE ROLE "role1" FROM USER "user1"`, u2)

	r2 := rg.Role("role2").Revoke()
	a.Equal(`REVOKE ROLE "role1" FROM ROLE "role2"`, r2)

}
