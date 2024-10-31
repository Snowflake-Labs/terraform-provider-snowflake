package resourceshowoutputassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ConnectionShowOutputAssert) HasPrimaryIdentifier(expected sdk.ExternalObjectIdentifier) *ConnectionShowOutputAssert {
	// expectedString := strings.ReplaceAll(expected.FullyQualifiedName(), `"`, "")
	c.AddAssertion(assert.ResourceShowOutputValueSet("primary", expected.Name()))
	return c
}

func (c *ConnectionShowOutputAssert) HasFailoverAllowedToAccounts(expected ...sdk.AccountIdentifier) *ConnectionShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValueSet("failover_allowed_to_accounts.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		c.AddAssertion(assert.ResourceShowOutputValueSet(fmt.Sprintf("failover_allowed_to_accounts.%d", i), v.Name()))
	}
	return c
}
