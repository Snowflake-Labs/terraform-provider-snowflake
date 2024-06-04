package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTestCheckResourceAttrNumberAtLeast(t *testing.T) {
	state := &terraform.State{
		Modules: []*terraform.ModuleState{
			{
				Path: []string{"root"},
				Resources: map[string]*terraform.ResourceState{
					"test": {
						Primary: &terraform.InstanceState{
							Attributes: map[string]string{
								"smaller": "10",
								"equal":   "20",
								"greater": "30",
								"not_int": "string_value",
							},
						},
					},
				},
			},
		},
	}

	assert.ErrorContains(t, TestCheckResourceAttrNumberAtLeast("test", "smaller", 20)(state), "expected attribute smaller to be at least 20, but was 10")
	assert.ErrorContains(t, TestCheckResourceAttrNumberAtLeast("test", "not_int", 20)(state), "failed to parse attribute not_int, err: strconv.Atoi: parsing \"string_value\": invalid syntax")
	assert.Nil(t, TestCheckResourceAttrNumberAtLeast("test", "equal", 20)(state))
	assert.Nil(t, TestCheckResourceAttrNumberAtLeast("test", "greater", 20)(state))
}
