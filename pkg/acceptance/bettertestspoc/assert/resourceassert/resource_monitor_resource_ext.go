package resourceassert

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (r *ResourceMonitorResourceAssert) HasStartTimestampNotEmpty() *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValuePresent("start_timestamp"))
	return r
}

func (r *ResourceMonitorResourceAssert) HasEndTimestampNotEmpty() *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValuePresent("end_timestamp"))
	return r
}

func (r *ResourceMonitorResourceAssert) HasNotifyUsersLen(len int) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet("notify_users.#", strconv.FormatInt(int64(len), 10)))
	return r
}

func (r *ResourceMonitorResourceAssert) HasNotifyUser(index int, userName string) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet(fmt.Sprintf("notify_users.%d", index), userName))
	return r
}

func (r *ResourceMonitorResourceAssert) HasTriggerLen(len int) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet("trigger.#", strconv.FormatInt(int64(len), 10)))
	return r
}

func (r *ResourceMonitorResourceAssert) HasTrigger(index int, threshold int, action sdk.TriggerAction) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet(fmt.Sprintf("trigger.%d.threshold", index), strconv.FormatInt(int64(threshold), 10)))
	r.AddAssertion(assert.ValueSet(fmt.Sprintf("trigger.%d.on_threshold_reached", index), string(action)))
	return r
}
