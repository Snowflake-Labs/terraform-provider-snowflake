package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

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

func (a *AccountShowOutputAssert) HasMarketplaceProviderBillingEntityNameNotEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("marketplace_provider_billing_entity_name"))
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

func (a *AccountShowOutputAssert) HasMovedOnEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("moved_on", ""))
	return a
}

func (a *AccountShowOutputAssert) HasOrganizationUrlExpirationOnEmpty() *AccountShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValueSet("organization_url_expiration_on", ""))
	return a
}
