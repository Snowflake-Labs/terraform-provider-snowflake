package resourceassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func (t *TagResourceAssert) HasNoMaskingPolicies() *TagResourceAssert {
	t.AddAssertion(assert.ValueSet("masking_policies.#", "0"))
	return t
}

func (t *TagResourceAssert) HasNoAllowedValues() *TagResourceAssert {
	t.AddAssertion(assert.ValueSet("allowed_values.#", "0"))
	return t
}
