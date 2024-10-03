package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (t *TaskResourceAssert) HasAfterLen(len int) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("after.#", strconv.FormatInt(int64(len), 10)))
	return t
}
