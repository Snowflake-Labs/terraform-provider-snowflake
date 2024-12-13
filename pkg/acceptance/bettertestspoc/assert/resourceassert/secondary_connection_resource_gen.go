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

func (s *SecondaryConnectionResourceAssert) HasIsPrimaryString(expected string) *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("is_primary", expected))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasNameString(expected string) *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("name", expected))
	return s
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

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

func (s *SecondaryConnectionResourceAssert) HasNoIsPrimary() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueNotSet("is_primary"))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasNoName() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueNotSet("name"))
	return s
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (s *SecondaryConnectionResourceAssert) HasAsReplicaOfEmpty() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("as_replica_of", ""))
	return s
}
func (s *SecondaryConnectionResourceAssert) HasCommentEmpty() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("comment", ""))
	return s
}
func (s *SecondaryConnectionResourceAssert) HasFullyQualifiedNameEmpty() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("fully_qualified_name", ""))
	return s
}
func (s *SecondaryConnectionResourceAssert) HasNameEmpty() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("name", ""))
	return s
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (s *SecondaryConnectionResourceAssert) HasAsReplicaOfNotEmpty() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValuePresent("as_replica_of"))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasCommentNotEmpty() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValuePresent("comment"))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasFullyQualifiedNameNotEmpty() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValuePresent("fully_qualified_name"))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasIsPrimaryNotEmpty() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValuePresent("is_primary"))
	return s
}

func (s *SecondaryConnectionResourceAssert) HasNameNotEmpty() *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValuePresent("name"))
	return s
}
