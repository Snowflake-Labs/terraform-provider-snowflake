package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (r *RowAccessPolicyShowOutputAssert) HasCreatedOnNotEmpty() *RowAccessPolicyShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return r
}
