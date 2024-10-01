package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"strconv"
)

func (t *TaskResourceAssert) HasAfterLen(len int) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("after.#", strconv.FormatInt(int64(len), 10)))
	return t
}
