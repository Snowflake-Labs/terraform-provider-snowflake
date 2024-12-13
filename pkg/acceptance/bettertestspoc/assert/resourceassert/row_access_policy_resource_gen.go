// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type RowAccessPolicyResourceAssert struct {
	*assert.ResourceAssert
}

func RowAccessPolicyResource(t *testing.T, name string) *RowAccessPolicyResourceAssert {
	t.Helper()

	return &RowAccessPolicyResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedRowAccessPolicyResource(t *testing.T, id string) *RowAccessPolicyResourceAssert {
	t.Helper()

	return &RowAccessPolicyResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (r *RowAccessPolicyResourceAssert) HasArgumentString(expected string) *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("argument", expected))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasBodyString(expected string) *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("body", expected))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasCommentString(expected string) *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("comment", expected))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasDatabaseString(expected string) *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("database", expected))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasFullyQualifiedNameString(expected string) *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasNameString(expected string) *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("name", expected))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasSchemaString(expected string) *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("schema", expected))
	return r
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (r *RowAccessPolicyResourceAssert) HasNoArgument() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("argument.#", "0"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasNoBody() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueNotSet("body"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasNoComment() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueNotSet("comment"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasNoDatabase() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueNotSet("database"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasNoFullyQualifiedName() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasNoName() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueNotSet("name"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasNoSchema() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueNotSet("schema"))
	return r
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (r *RowAccessPolicyResourceAssert) HasBodyEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("body", ""))
	return r
}
func (r *RowAccessPolicyResourceAssert) HasCommentEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("comment", ""))
	return r
}
func (r *RowAccessPolicyResourceAssert) HasDatabaseEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("database", ""))
	return r
}
func (r *RowAccessPolicyResourceAssert) HasFullyQualifiedNameEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("fully_qualified_name", ""))
	return r
}
func (r *RowAccessPolicyResourceAssert) HasNameEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("name", ""))
	return r
}
func (r *RowAccessPolicyResourceAssert) HasSchemaEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("schema", ""))
	return r
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (r *RowAccessPolicyResourceAssert) HasArgumentNotEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValuePresent("argument"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasBodyNotEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValuePresent("body"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasCommentNotEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValuePresent("comment"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasDatabaseNotEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValuePresent("database"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasFullyQualifiedNameNotEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValuePresent("fully_qualified_name"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasNameNotEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValuePresent("name"))
	return r
}

func (r *RowAccessPolicyResourceAssert) HasSchemaNotEmpty() *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValuePresent("schema"))
	return r
}
