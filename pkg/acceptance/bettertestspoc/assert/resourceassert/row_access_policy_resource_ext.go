package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (r *RowAccessPolicyResourceAssert) HasArguments(args []sdk.RowAccessPolicyArgument) *RowAccessPolicyResourceAssert {
	r.AddAssertion(assert.ValueSet("argument.#", strconv.FormatInt(int64(len(args)), 10)))
	for i, v := range args {
		r.AddAssertion(assert.ValueSet(fmt.Sprintf("argument.%d.name", i), v.Name))
		r.AddAssertion(assert.ValueSet(fmt.Sprintf("argument.%d.type", i), v.Type))
	}
	return r
}
