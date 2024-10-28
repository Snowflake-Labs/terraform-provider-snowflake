package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ConnectionResourceAssert) HasEnableFailover(expected []sdk.AccountIdentifier) *ConnectionResourceAssert {
	for i, v := range expected {
		c.AddAssertion(assert.ValueSet(fmt.Sprintf("failover_allowed_to_accounts.0.to_accounts.%d.account_identifier", i), v.Name()))
	}
	return c
}
