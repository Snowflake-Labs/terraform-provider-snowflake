package resourceshowoutputassert

import (
	"fmt"
	"strconv"
)

// HasCategories asserts the whole list. The generator emits helpers for scalars only, and categories is a list
// of strings. Named after the generated objectassert equivalent.
func (o *OpenflowConnectorDefinitionShowOutputAssert) HasCategories(expected ...string) *OpenflowConnectorDefinitionShowOutputAssert {
	o.ValueSet("categories.#", strconv.Itoa(len(expected)))
	for i, category := range expected {
		o.ValueSet(fmt.Sprintf("categories.%d", i), category)
	}
	return o
}
