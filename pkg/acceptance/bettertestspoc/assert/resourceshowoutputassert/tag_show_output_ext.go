package resourceshowoutputassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *TagShowOutputAssert) HasCreatedOnNotEmpty() *TagShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return s
}

func (s *TagShowOutputAssert) HasAllowedValues(expected ...string) *TagShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("allowed_values.#", strconv.FormatInt(int64(len(expected)), 10)))
	for i := range expected {
		s.AddAssertion(assert.ResourceShowOutputValueSet(fmt.Sprintf("allowed_values.%d", i), expected[i]))
	}
	return s
}
