package resourceshowoutputassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (t *TaskShowOutputAssert) HasCreatedOnNotEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return t
}

func (t *TaskShowOutputAssert) HasIdNotEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValuePresent("id"))
	return t
}

func (t *TaskShowOutputAssert) HasLastCommittedOnNotEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValuePresent("last_committed_on"))
	return t
}

func (t *TaskShowOutputAssert) HasLastSuspendedOnNotEmpty() *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValuePresent("last_suspended_on"))
	return t
}

func (t *TaskShowOutputAssert) HasPredecessors(predecessors ...sdk.SchemaObjectIdentifier) *TaskShowOutputAssert {
	t.AddAssertion(assert.ResourceShowOutputValueSet("predecessors.#", strconv.Itoa(len(predecessors))))
	for i, predecessor := range predecessors {
		t.AddAssertion(assert.ResourceShowOutputValueSet(fmt.Sprintf("predecessors.%d", i), predecessor.FullyQualifiedName()))
	}
	return t
}
