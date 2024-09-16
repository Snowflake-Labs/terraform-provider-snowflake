package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (p *MaskingPolicyShowOutputAssert) HasCreatedOnNotEmpty() *MaskingPolicyShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return p
}
