package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (u *UserResourceAssert) HasDisabled(expected bool) *UserResourceAssert {
	u.AddAssertion(assert.ValueSet("disabled", strconv.FormatBool(expected)))
	return u
}

func (u *UserResourceAssert) HasEmptyPassword() *UserResourceAssert {
	u.AddAssertion(assert.ValueSet("password", ""))
	return u
}

func (u *UserResourceAssert) HasDefaultSecondaryRoles(roles ...string) *UserResourceAssert {
	for idx, role := range roles {
		u.AddAssertion(assert.ValueSet(fmt.Sprintf("default_secondary_roles.%d", idx), role))
	}
	return u
}

func (u *UserResourceAssert) HasDefaultSecondaryRolesEmpty() *UserResourceAssert {
	u.AddAssertion(assert.ValueSet("default_secondary_roles.#", "0"))
	return u
}

func (u *UserResourceAssert) HasMustChangePassword(expected bool) *UserResourceAssert {
	u.AddAssertion(assert.ValueSet("must_change_password", strconv.FormatBool(expected)))
	return u
}
