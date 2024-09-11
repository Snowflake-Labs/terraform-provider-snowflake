package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (r *ResourceMonitorAssert) HasNonEmptyStartTime() *ResourceMonitorAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.ResourceMonitor) error {
		t.Helper()
		if o.StartTime == "" {
			return fmt.Errorf("expected start time to be non empty")
		}
		return nil
	})
	return r
}

func (r *ResourceMonitorAssert) HasNonEmptyEndTime() *ResourceMonitorAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.ResourceMonitor) error {
		t.Helper()
		if o.StartTime == "" {
			return fmt.Errorf("expected end time to be non empty")
		}
		return nil
	})
	return r
}
