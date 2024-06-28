package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	showOutputAttributeName     = "show_output"
	describeOutputAttributeName = "describe_output"
)

// handleExternalChangesToObjectInShow assumes that show output is kept in showOutputAttributeName attribute
func handleExternalChangesToObjectInShow(d *schema.ResourceData, mappings ...showMapping) error {
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

// handleExternalChangesToObjectInDescribe assumes that show output is kept in describeOutputAttributeName attribute
func handleExternalChangesToObjectInDescribe(d *schema.ResourceData, mappings ...describeMapping) error {
	if describeOutput, ok := d.GetOk(describeOutputAttributeName); ok {
		describeOutputList := describeOutput.([]any)
		if len(describeOutputList) == 1 {
			result := describeOutputList[0].(map[string]any)

			for _, mapping := range mappings {
				if result[mapping.nameInDescribe] == nil {
					continue
				}

				valueToCompareFromList := result[mapping.nameInDescribe].([]any)
				if len(valueToCompareFromList) != 1 {
					continue
				}

				valueToCompareFrom := valueToCompareFromList[0].(map[string]any)["value"]
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

type describeMapping struct {
	nameInDescribe string
	nameInConfig   string
	valueToCompare any
	valueToSet     any
	normalizeFunc  func(any) any
}
