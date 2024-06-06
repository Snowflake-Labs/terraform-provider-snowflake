package importchecks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ComposeImportStateCheck is based on unexported composeImportStateCheck from teststep_providers_test.go
func ComposeImportStateCheck(fs ...resource.ImportStateCheckFunc) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		for i, f := range fs {
			if err := f(s); err != nil {
				return fmt.Errorf("check %d/%d error: %s", i+1, len(fs), err)
			}
		}
		return nil
	}
}

// TestCheckResourceAttrInstanceState is based on unexported testCheckResourceAttrInstanceState from teststep_providers_test.go
func TestCheckResourceAttrInstanceState(id string, attributeName, attributeValue string) resource.ImportStateCheckFunc {
	return func(is []*terraform.InstanceState) error {
		for _, v := range is {
			if v.ID != id {
				continue
			}

			if attrVal, ok := v.Attributes[attributeName]; ok {
				if attrVal != attributeValue {
					return fmt.Errorf("expected: %s got: %s", attributeValue, attrVal)
				}

				return nil
			}
		}

		return fmt.Errorf("attribute %s not found in instance state", attributeName)
	}
}
