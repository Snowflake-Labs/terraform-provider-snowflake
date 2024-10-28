package resourceshowoutputassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ConnectionShowOutputAssert) HasPrimaryIdentifier(expected sdk.ExternalObjectIdentifier) *ConnectionShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValueSet("primary", expected.FullyQualifiedName()))
	return c
}

func (c *ConnectionShowOutputAssert) HasFailoverAllowedToAccounts(expected []sdk.AccountIdentifier) *ConnectionShowOutputAssert {
	for i, v := range expected {
		c.AddAssertion(assert.ResourceShowOutputValueSet(fmt.Sprintf("failover_allowed_to_accounts.0.to_accounts.%d.account_identifier", i), v.Name()))
	}
	return c
}
