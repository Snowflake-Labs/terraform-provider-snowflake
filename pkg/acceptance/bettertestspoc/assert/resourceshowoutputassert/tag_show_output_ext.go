package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// TagsDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func TagsDatasourceShowOutput(t *testing.T, name string) *TagShowOutputAssert {
	t.Helper()

	s := TagShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "tags.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
