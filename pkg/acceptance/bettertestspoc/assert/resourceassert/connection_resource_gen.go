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

func (c *ConnectionResourceAssert) HasAsReplicaOfString(expected string) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("as_replica_of", expected))
	return c
}

func (c *ConnectionResourceAssert) HasCommentString(expected string) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("comment", expected))
	return c
}

func (c *ConnectionResourceAssert) HasEnableFailoverString(expected string) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("enable_failover", expected))
	return c
}

func (c *ConnectionResourceAssert) HasFullyQualifiedNameString(expected string) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return c
}

func (c *ConnectionResourceAssert) HasIsPrimaryString(expected string) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("is_primary", expected))
	return c
}

func (c *ConnectionResourceAssert) HasNameString(expected string) *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("name", expected))
	return c
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (c *ConnectionResourceAssert) HasNoAsReplicaOf() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueNotSet("as_replica_of"))
	return c
}

func (c *ConnectionResourceAssert) HasNoComment() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueNotSet("comment"))
	return c
}

func (c *ConnectionResourceAssert) HasNoEnableFailover() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueNotSet("enable_failover"))
	return c
}

func (c *ConnectionResourceAssert) HasNoFullyQualifiedName() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return c
}

func (c *ConnectionResourceAssert) HasNoIsPrimary() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueNotSet("is_primary"))
	return c
}

func (c *ConnectionResourceAssert) HasNoName() *ConnectionResourceAssert {
	c.AddAssertion(assert.ValueNotSet("name"))
	return c
}
