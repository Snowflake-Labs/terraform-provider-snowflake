package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// RowAccessPoliciesDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func RowAccessPoliciesDatasourceShowOutput(t *testing.T, name string) *RowAccessPolicyShowOutputAssert {
	t.Helper()

	r := RowAccessPolicyShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "row_access_policies.0."),
	}
	r.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &r
}
