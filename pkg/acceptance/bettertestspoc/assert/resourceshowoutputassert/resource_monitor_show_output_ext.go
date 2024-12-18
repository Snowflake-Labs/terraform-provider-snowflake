package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// ResourceMonitorDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func ResourceMonitorDatasourceShowOutput(t *testing.T, name string) *ResourceMonitorShowOutputAssert {
	t.Helper()

	u := ResourceMonitorShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "resource_monitors.0."),
	}
	u.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &u
}
