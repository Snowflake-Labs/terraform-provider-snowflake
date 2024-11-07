package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (t *TaskResourceAssert) HasAfterIdsInOrder(ids ...sdk.SchemaObjectIdentifier) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("after.#", strconv.FormatInt(int64(len(ids)), 10)))
	for i, id := range ids {
		t.AddAssertion(assert.ValueSet(fmt.Sprintf("after.%d", i), id.FullyQualifiedName()))
	}
	return t
}

func (t *TaskResourceAssert) HasScheduleMinutes(minutes int) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("schedule.#", "1"))
	t.AddAssertion(assert.ValueSet("schedule.0.minutes", strconv.Itoa(minutes)))
	return t
}

func (t *TaskResourceAssert) HasScheduleCron(cron string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("schedule.#", "1"))
	t.AddAssertion(assert.ValueSet("schedule.0.using_cron", cron))
	return t
}

func (t *TaskResourceAssert) HasNoScheduleSet() *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("schedule.#", "0"))
	return t
}
