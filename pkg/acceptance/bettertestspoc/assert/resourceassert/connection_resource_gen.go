// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type ConnectionResourceAssert struct {
	*assert.ResourceAssert
}

func ConnectionResource(t *testing.T, name string) *ConnectionResourceAssert {
	t.Helper()

	return &ConnectionResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedConnectionResource(t *testing.T, id string) *ConnectionResourceAssert {
	t.Helper()

	return &ConnectionResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (c *ConnectionResourceAssert) HasCommentString(expected string) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("comment", expected))
	return c
}

func (c *ConnectionResourceAssert) HasEnableFailoverToAccountsString(expected string) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("enable_failover_to_accounts", expected))
	return c
}

func (c *ConnectionResourceAssert) HasFullyQualifiedNameString(expected string) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return c
}

func (c *ConnectionResourceAssert) HasNameString(expected string) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("name", expected))
	return c
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (c *ConnectionResourceAssert) HasNoComment() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueNotSet("comment"))
	return c
}

/*
func (c *ConnectionResourceAssert) HasNoEnableFailoverToAccounts() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueNotSet("enable_failover_to_accounts"))
	return c
}
*/

func (c *ConnectionResourceAssert) HasNoFullyQualifiedName() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return c
}

func (c *ConnectionResourceAssert) HasNoName() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueNotSet("name"))
	return c
}
