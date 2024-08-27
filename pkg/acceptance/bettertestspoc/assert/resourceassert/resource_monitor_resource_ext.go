package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (r *ResourceMonitorResourceAssert) HasNotifyUsersLen(len int) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet("notify_users.#", strconv.FormatInt(int64(len), 10)))
	return r
}

func (r *ResourceMonitorResourceAssert) HasTriggerLen(len int) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet("trigger.#", strconv.FormatInt(int64(len), 10)))
	return r
}
