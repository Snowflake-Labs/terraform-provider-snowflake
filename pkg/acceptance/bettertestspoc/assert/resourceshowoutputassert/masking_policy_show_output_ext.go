package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// MaskingPoliciesDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func MaskingPoliciesDatasourceShowOutput(t *testing.T, name string) *MaskingPolicyShowOutputAssert {
	t.Helper()

	m := MaskingPolicyShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "masking_policies.0."),
	}
	m.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &m
}
