package model

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func (t *TagModel) WithMaskingPoliciesValue(value tfconfig.Variable) *TagModel {
	t.MaskingPolicies = value
	return t
}
