package resourceshowoutputassert

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ConnectionShowOutputAssert) HasPrimaryIdentifier(expected sdk.ExternalObjectIdentifier) *ConnectionShowOutputAssert {
	expectedString := strings.ReplaceAll(expected.FullyQualifiedName(), `"`, "")
	c.AddAssertion(assert.ResourceShowOutputValueSet("primary", expectedString))
	return c
}

func (c *ConnectionShowOutputAssert) HasFailoverAllowedToAccounts(expected ...sdk.AccountIdentifier) *ConnectionShowOutputAssert {
	for i, v := range expected {
		c.AddAssertion(assert.ResourceShowOutputValueSet(fmt.Sprintf("failover_allowed_to_accounts.%d", i), v.Name()))
	}
	return c
}

func (c *ConnectionShowOutputAssert) HasNoFailoverAllowedToAccounts() *ConnectionShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValueSet("failover_allowed_to_accounts.#", "0"))
	return c
}
