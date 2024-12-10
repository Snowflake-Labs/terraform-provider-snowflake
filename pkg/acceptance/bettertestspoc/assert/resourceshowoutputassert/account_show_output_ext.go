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

func (a *AccountShowOutputAssert) HasAccountUrlNotEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("account_url"))
	return a
}

func (a *AccountShowOutputAssert) HasCreatedOnNotEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return a
}

func (a *AccountShowOutputAssert) HasAccountLocatorNotEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("account_locator"))
	return a
}

func (a *AccountShowOutputAssert) HasAccountLocatorUrlNotEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("account_locator_url"))
	return a
}

func (a *AccountShowOutputAssert) HasConsumptionBillingEntityNameNotEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("consumption_billing_entity_name"))
	return a
}

func (a *AccountShowOutputAssert) HasNoOrganizationOldUrl() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("organization_old_url"))
	return a
}

func (a *AccountShowOutputAssert) HasMarketplaceProviderBillingEntityNameNotEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("marketplace_provider_billing_entity_name"))
	return a
}

func (a *AccountShowOutputAssert) HasNoAccountOldUrlSavedOn() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("account_old_url_saved_on"))
	return a
}

func (a *AccountShowOutputAssert) HasNoAccountOldUrlLastUsed() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("account_old_url_last_used"))
	return a
}

func (a *AccountShowOutputAssert) HasNoOrganizationOldUrlSavedOn() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("organization_old_url_saved_on"))
	return a
}

func (a *AccountShowOutputAssert) HasNoOrganizationOldUrlLastUsed() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("organization_old_url_last_used"))
	return a
}

func (a *AccountShowOutputAssert) HasNoDroppedOn() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("dropped_on"))
	return a
}

func (a *AccountShowOutputAssert) HasNoScheduledDeletionTime() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("scheduled_deletion_time"))
	return a
}

func (a *AccountShowOutputAssert) HasNoRestoredOn() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("restored_on"))
	return a
}

func (a *AccountShowOutputAssert) HasNoMovedToOrganization() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("moved_to_organization"))
	return a
}

func (a *AccountShowOutputAssert) HasMovedOnEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("moved_on"))
	return a
}

func (a *AccountShowOutputAssert) HasNoOrganizationUrlExpirationOn() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("organization_url_expiration_on"))
	return a
}

func (a *AccountShowOutputAssert) HasNoIsEventsAccount() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("is_events_account"))
	return a
}

func (a *AccountShowOutputAssert) HasNoIsOrganizationAccount() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueNotSet("is_organization_account"))
	return a
}
