package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func AccountDatasourceShowOutput(t *testing.T, name string) *AccountShowOutputAssert {
	t.Helper()

	a := AccountShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "accounts.0."),
	}
	a.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &a
}

// TODO: Why those are not duplicates

func (a *AccountShowOutputAssert) HasAccountUrlNotEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("account_url"))
	return a
}

func (a *AccountShowOutputAssert) HasAccountOldUrlSavedOnEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("account_old_url_saved_on", ""))
	return a
}

func (a *AccountShowOutputAssert) HasAccountOldUrlLastUsedEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("account_old_url_last_used", ""))
	return a
}

func (a *AccountShowOutputAssert) HasOrganizationOldUrlSavedOnEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("organization_old_url_saved_on", ""))
	return a
}

func (a *AccountShowOutputAssert) HasOrganizationOldUrlLastUsedEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("organization_old_url_last_used", ""))
	return a
}

func (a *AccountShowOutputAssert) HasDroppedOnEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("dropped_on", ""))
	return a
}

func (a *AccountShowOutputAssert) HasScheduledDeletionTimeEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("scheduled_deletion_time", ""))
	return a
}

func (a *AccountShowOutputAssert) HasRestoredOnEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("restored_on", ""))
	return a
}

func (a *AccountShowOutputAssert) HasOrganizationUrlExpirationOnEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("organization_url_expiration_on", ""))
	return a
}

func (a *AccountShowOutputAssert) HasIsEventsAccountEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("is_events_account", ""))
	return a
}

func (a *AccountShowOutputAssert) HasIsOrganizationAccountEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("is_organization_account", ""))
	return a
}
