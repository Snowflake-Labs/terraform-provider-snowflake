package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (t *TagAssociationResourceAssert) HasObjectIdentifiersLength(len int) *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("object_identifiers.#", fmt.Sprintf("%d", len)))
	return t
}

func (t *TagAssociationResourceAssert) HasIdString(expected string) *TagAssociationResourceAssert {
	t.AddAssertion(assert.ValueSet("id", expected))
	return t
}
