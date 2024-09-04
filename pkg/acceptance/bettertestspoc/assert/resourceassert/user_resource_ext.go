package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (u *UserResourceAssert) HasDisabled(expected bool) *UserResourceAssert {
	u.AddAssertion(assert.ValueSet("disabled", strconv.FormatBool(expected)))
	return u
}

func (u *UserResourceAssert) HasEmptyPassword() *UserResourceAssert {
	u.AddAssertion(assert.ValueSet("password", ""))
	return u
}

func (u *UserResourceAssert) HasMustChangePassword(expected bool) *UserResourceAssert {
	u.AddAssertion(assert.ValueSet("must_change_password", strconv.FormatBool(expected)))
	return u
}

func (u *UserResourceAssert) HasDefaultSecondaryRolesOption(expected sdk.SecondaryRolesOption) *UserResourceAssert {
	return u.HasDefaultSecondaryRolesOptionString(string(expected))
}
