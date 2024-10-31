// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type SecondaryConnectionResourceAssert struct {
	*assert.ResourceAssert
}

func SecondaryConnectionResource(t *testing.T, name string) *SecondaryConnectionResourceAssert {
	t.Helper()

	return &SecondaryConnectionResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedSecondaryConnectionResource(t *testing.T, id string) *SecondaryConnectionResourceAssert {
	t.Helper()

	return &SecondaryConnectionResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (s *SecondaryConnectionResourceAssert) HasAsReplicaOfString(expected string) *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("as_replica_of", expected))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasCommentString(expected string) *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("comment", expected))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasFullyQualifiedNameString(expected string) *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasNameString(expected string) *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("name", expected))
	return s
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (s *SecondaryConnectionResourceAssert) HasNoAsReplicaOf() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueNotSet("as_replica_of"))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasNoComment() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueNotSet("comment"))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasNoFullyQualifiedName() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasNoName() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueNotSet("name"))
	return s
}