// Code generated by assertions generator; DO NOT EDIT.

package resourceshowoutputassert

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// to ensure sdk package is used
var _ = sdk.Object{}

type TagShowOutputAssert struct {
	*assert.ResourceAssert
}

func TagShowOutput(t *testing.T, name string) *TagShowOutputAssert {
	t.Helper()

	tt := TagShowOutputAssert{
		ResourceAssert: assert.NewResourceAssert(name, "show_output"),
	}
	tt.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &tt
}

func ImportedTagShowOutput(t *testing.T, id string) *TagShowOutputAssert {
	t.Helper()

	tt := TagShowOutputAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "show_output"),
	}
	tt.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &tt
}

////////////////////////////
// Attribute value checks //
////////////////////////////

func (t *TagShowOutputAssert) HasCreatedOn(expected time.Time) *TagShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("created_on", expected.String()))
	return t
}

func (t *TagShowOutputAssert) HasName(expected string) *TagShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("name", expected))
	return t
}

func (t *TagShowOutputAssert) HasDatabaseName(expected string) *TagShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("database_name", expected))
	return t
}

func (t *TagShowOutputAssert) HasSchemaName(expected string) *TagShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("schema_name", expected))
	return t
}

func (t *TagShowOutputAssert) HasOwner(expected string) *TagShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("owner", expected))
	return t
}

func (t *TagShowOutputAssert) HasComment(expected string) *TagShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("comment", expected))
	return t
}

func (t *TagShowOutputAssert) HasOwnerRoleType(expected string) *TagShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("owner_role_type", expected))
	return t
}
