// Code generated by assertions generator; DO NOT EDIT.

package resourceshowoutputassert

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// to ensure sdk package is used
var _ = sdk.Object{}

type ResourceMonitorShowOutputAssert struct {
	*assert.ResourceAssert
}

func ResourceMonitorShowOutput(t *testing.T, name string) *ResourceMonitorShowOutputAssert {
	t.Helper()

	r := ResourceMonitorShowOutputAssert{
		ResourceAssert: assert.NewResourceAssert(name, "show_output"),
	}
	r.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &r
}

func ImportedResourceMonitorShowOutput(t *testing.T, id string) *ResourceMonitorShowOutputAssert {
	t.Helper()

	r := ResourceMonitorShowOutputAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "show_output"),
	}
	r.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &r
}

////////////////////////////
// Attribute value checks //
////////////////////////////

func (r *ResourceMonitorShowOutputAssert) HasName(expected string) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValueSet("name", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasCreditQuota(expected float64) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputFloatValueSet("credit_quota", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasUsedCredits(expected float64) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputFloatValueSet("used_credits", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasRemainingCredits(expected float64) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputFloatValueSet("remaining_credits", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasLevel(expected sdk.ResourceMonitorLevel) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("level", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasFrequency(expected sdk.Frequency) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("frequency", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasStartTime(expected string) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValueSet("start_time", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasEndTime(expected string) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValueSet("end_time", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasSuspendAt(expected int) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputIntValueSet("suspend_at", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasSuspendImmediateAt(expected int) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputIntValueSet("suspend_immediate_at", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasCreatedOn(expected time.Time) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValueSet("created_on", expected.String()))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasOwner(expected string) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValueSet("owner", expected))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasComment(expected string) *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValueSet("comment", expected))
	return r
}
