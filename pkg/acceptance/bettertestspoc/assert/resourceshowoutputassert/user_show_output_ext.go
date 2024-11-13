package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// TaskDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func TaskDatasourceShowOutput(t *testing.T, name string) *TaskShowOutputAssert {
	t.Helper()

	taskAssert := TaskShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "tasks.0."),
	}
	taskAssert.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &taskAssert
}

func (t *TaskShowOutputAssert) HasErrorIntegrationEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("error_integration", ""))
	return t
}

func (u *UserShowOutputAssert) HasCreatedOnNotEmpty() *UserShowOutputAssert {
	u.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return u
}

func (u *UserShowOutputAssert) HasDaysToExpiryNotEmpty() *UserShowOutputAssert {
	u.AddAssertion(assert.ResourceShowOutputValuePresent("days_to_expiry"))
	return u
}

func (u *UserShowOutputAssert) HasMinsToUnlockNotEmpty() *UserShowOutputAssert {
	u.AddAssertion(assert.ResourceShowOutputValuePresent("mins_to_unlock"))
	return u
}

func (u *UserShowOutputAssert) HasMinsToBypassMfaNotEmpty() *UserShowOutputAssert {
	u.AddAssertion(assert.ResourceShowOutputValuePresent("mins_to_bypass_mfa"))
	return u
}
