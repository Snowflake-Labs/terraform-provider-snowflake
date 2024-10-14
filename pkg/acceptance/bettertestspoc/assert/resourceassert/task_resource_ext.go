package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (t *TaskResourceAssert) HasAfterIds(ids ...sdk.SchemaObjectIdentifier) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("after.#", strconv.FormatInt(int64(len(ids)), 10)))
	for i, id := range ids {
		t.AddAssertion(assert.ValueSet(fmt.Sprintf("after.%d", i), id.FullyQualifiedName()))
	}
	return t
}
