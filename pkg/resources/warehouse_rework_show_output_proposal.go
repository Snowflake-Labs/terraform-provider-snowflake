package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	showOutputAttributeName     = "show_output"
	describeOutputAttributeName = "describe_output"
)

// handleExternalChangesToObject assumes that show output is kept in showOutputAttributeName attribute
func handleExternalChangesToObject(d *schema.ResourceData, mappings ...showMapping) error {
	if showOutput, ok := d.GetOk(showOutputAttributeName); ok {
		showOutputList := showOutput.([]any)
		if len(showOutputList) == 1 {
			result := showOutputList[0].(map[string]any)
			for _, mapping := range mappings {
				valueToCompareFrom := result[mapping.nameInShow]
				if mapping.normalizeFunc != nil {
					valueToCompareFrom = mapping.normalizeFunc(valueToCompareFrom)
				}
				if valueToCompareFrom != mapping.valueToCompare {
					if err := d.Set(mapping.nameInConfig, mapping.valueToSet); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// TODO: merge with above
func handleExternalChangesToObjectDescribe(d *schema.ResourceData, mappings ...showMapping) error {
	if showOutput, ok := d.GetOk(describeOutputAttributeName); ok {
		showOutputList := showOutput.([]any)
		if len(showOutputList) == 1 {
			result := showOutputList[0].(map[string]any)
			for _, mapping := range mappings {
				valueToCompareFrom := result[mapping.nameInShow]
				if mapping.normalizeFunc != nil {
					valueToCompareFrom = mapping.normalizeFunc(valueToCompareFrom)
				}
				if valueToCompareFrom != mapping.valueToCompare {
					if err := d.Set(mapping.nameInConfig, mapping.valueToSet); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

type showMapping struct {
	nameInShow     string
	nameInConfig   string
	valueToCompare any
	valueToSet     any
	normalizeFunc  func(any) any
}
