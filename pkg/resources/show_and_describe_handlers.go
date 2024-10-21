package resources

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ShowOutputAttributeName        = "show_output"
	DescribeOutputAttributeName    = "describe_output"
	ParametersAttributeName        = "parameters"
	RelatedParametersAttributeName = "related_parameters"
)

// handleExternalChangesToObjectInShow assumes that show output is kept in ShowOutputAttributeName attribute
func handleExternalChangesToObjectInShow(d *schema.ResourceData, mappings ...showMapping) error {
	if showOutput, ok := d.GetOk(ShowOutputAttributeName); ok {
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

// handleExternalValueChangesToObjectInDescribe assumes that describe output is kept in DescribeOutputAttributeName attribute
func handleExternalValueChangesToObjectInDescribe(d *schema.ResourceData, mappings ...describeMapping) error {
	if descOutput, ok := d.GetOk(DescribeOutputAttributeName); ok {
		descOutputList := descOutput.([]any)
		if len(descOutputList) == 1 {
			result := descOutputList[0].(map[string]any)
			for _, mapping := range mappings {
				valueToCompareFrom := result[mapping.nameInDescribe]
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

// handleExternalChangesToObjectInDescribe assumes that show output is kept in DescribeOutputAttributeName attribute
func handleExternalChangesToObjectInDescribe(d *schema.ResourceData, mappings ...describeMapping) error {
	if describeOutput, ok := d.GetOk(DescribeOutputAttributeName); ok {
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

// setStateToValuesFromConfig currently handles only int, float, and string types.
// It's needed for the case where:
// - previous config was empty (therefore Snowflake defaults had been used)
// - new config have the same values that are already in SF
func setStateToValuesFromConfig(d *schema.ResourceData, resourceSchema map[string]*schema.Schema, fields []string) error {
	if !d.GetRawConfig().IsNull() {
		vMap := d.GetRawConfig().AsValueMap()
		for _, field := range fields {
			if v, ok := vMap[field]; ok && !v.IsNull() {
				if schemaField, ok := resourceSchema[field]; ok {
					switch schemaField.Type {
					case schema.TypeInt:
						intVal, _ := v.AsBigFloat().Int64()
						if err := d.Set(field, intVal); err != nil {
							return err
						}
					case schema.TypeFloat:
						if err := d.Set(field, v.AsBigFloat()); err != nil {
							return err
						}
					case schema.TypeString:
						if err := d.Set(field, v.AsString()); err != nil {
							return err
						}
					case schema.TypeSet:
						if err := d.Set(field, ctyValToSliceString(v.AsValueSlice())); err != nil {
							return err
						}
					default:
						log.Printf("[DEBUG] field %s has unsupported schema type %v not found", field, schemaField.Type)
					}
				} else {
					log.Printf("[DEBUG] schema field %s not found", field)
				}
			}
		}
	}
	return nil
}
