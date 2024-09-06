package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func (r *ResourceMonitorShowOutputAssert) HasStartTimeNotEmpty() *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("start_time"))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasEndTimeNotEmpty() *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("end_time"))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasCreatedOnNotEmpty() *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasOwnerNotEmpty() *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("owner"))
	return r
}
