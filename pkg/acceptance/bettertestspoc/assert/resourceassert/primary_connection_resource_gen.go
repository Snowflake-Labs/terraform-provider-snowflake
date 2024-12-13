// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type PrimaryConnectionResourceAssert struct {
	*assert.ResourceAssert
}

func PrimaryConnectionResource(t *testing.T, name string) *PrimaryConnectionResourceAssert {
	t.Helper()

	return &PrimaryConnectionResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedPrimaryConnectionResource(t *testing.T, id string) *PrimaryConnectionResourceAssert {
	t.Helper()

	return &PrimaryConnectionResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (p *PrimaryConnectionResourceAssert) HasCommentString(expected string) *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueSet("comment", expected))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasEnableFailoverToAccountsString(expected string) *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueSet("enable_failover_to_accounts", expected))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasFullyQualifiedNameString(expected string) *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasIsPrimaryString(expected string) *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueSet("is_primary", expected))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasNameString(expected string) *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueSet("name", expected))
	return p
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (p *PrimaryConnectionResourceAssert) HasNoComment() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueNotSet("comment"))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasNoEnableFailoverToAccounts() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueSet("enable_failover_to_accounts.#", "0"))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasNoFullyQualifiedName() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasNoIsPrimary() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueNotSet("is_primary"))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasNoName() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueNotSet("name"))
	return p
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (p *PrimaryConnectionResourceAssert) HasCommentEmpty() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueSet("comment", ""))
	return p
}
func (p *PrimaryConnectionResourceAssert) HasFullyQualifiedNameEmpty() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueSet("fully_qualified_name", ""))
	return p
}
func (p *PrimaryConnectionResourceAssert) HasNameEmpty() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValueSet("name", ""))
	return p
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (p *PrimaryConnectionResourceAssert) HasCommentNotEmpty() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValuePresent("comment"))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasEnableFailoverToAccountsNotEmpty() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValuePresent("enable_failover_to_accounts"))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasFullyQualifiedNameNotEmpty() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValuePresent("fully_qualified_name"))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasIsPrimaryNotEmpty() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValuePresent("is_primary"))
	return p
}

func (p *PrimaryConnectionResourceAssert) HasNameNotEmpty() *PrimaryConnectionResourceAssert {
	p.AddAssertion(assert.ValuePresent("name"))
	return p
}
