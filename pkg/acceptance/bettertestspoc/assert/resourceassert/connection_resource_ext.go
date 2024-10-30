package resourceassert

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ConnectionResourceAssert) HasAsReplicaOfIdentifier(expected sdk.ExternalObjectIdentifier) *ConnectionResourceAssert {
	expectedString := strings.ReplaceAll(expected.FullyQualifiedName(), `"`, "")
	c.AddAssertion(assert.ValueSet("as_replica_of", expectedString))
	return c
}

func (c *ConnectionResourceAssert) HasEnableFailoverToAccounts(expected ...sdk.AccountIdentifier) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("enable_failover_to_accounts.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		c.AddAssertion(assert.ValueSet(fmt.Sprintf("enable_failover_to_accounts.%d", i), v.Name()))
	}
	return c
}

func (c *ConnectionResourceAssert) HasNoEnableFailoverToAccounts() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("enable_failover_to_accounts.#", "0"))
	return c
}
