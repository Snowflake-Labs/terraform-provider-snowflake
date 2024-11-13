package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *TagResourceAssert) HasMaskingPolicies(expected ...sdk.SchemaObjectIdentifier) *TagResourceAssert {
	s.AddAssertion(assert.ValueSet("masking_policies.#", strconv.FormatInt(int64(len(expected)), 10)))
	for i := range expected {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("masking_policies.%d", i), expected[i].FullyQualifiedName()))
	}
	return s
}
