package bettertestspoc

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	_ TestCheckFuncProvider        = (*ResourceAssert)(nil)
	_ ImportStateCheckFuncProvider = (*ResourceAssert)(nil)
)

type ResourceAssert struct {
	name       string
	id         string
	prefix     string
	assertions []resourceAssertion
	TestCheckFuncProvider
}

func NewResourceAssert(name string, prefix string) *ResourceAssert {
	return &ResourceAssert{
		name:       name,
		prefix:     prefix,
		assertions: make([]resourceAssertion, 0),
	}
}

func NewImportedResourceAssert(id string, prefix string) *ResourceAssert {
	return &ResourceAssert{
		id:         id,
		prefix:     prefix,
		assertions: make([]resourceAssertion, 0),
	}
}

type resourceAssertionType string

const (
	resourceAssertionTypeValueSet    = "VALUE_SET"
	resourceAssertionTypeValueNotSet = "VALUE_NOT_SET"
)

type resourceAssertion struct {
	fieldName             string
	expectedValue         string
	resourceAssertionType resourceAssertionType
}

func valueSet(fieldName string, expected string) resourceAssertion {
	return resourceAssertion{fieldName: fieldName, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

func valueNotSet(fieldName string) resourceAssertion {
	return resourceAssertion{fieldName: fieldName, resourceAssertionType: resourceAssertionTypeValueNotSet}
}

const showOutputPrefix = "show_output.0."

func showOutputValueSet(fieldName string, expected string) resourceAssertion {
	return resourceAssertion{fieldName: showOutputPrefix + fieldName, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

const (
	parametersPrefix      = "parameters.0."
	parametersValueSuffix = ".0.value"
	parametersLevelSuffix = ".0.level"
)

func parameterValueSet(fieldName string, expected string) resourceAssertion {
	return resourceAssertion{fieldName: parametersPrefix + fieldName + parametersValueSuffix, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

func parameterLevelSet(fieldName string, expected string) resourceAssertion {
	return resourceAssertion{fieldName: parametersPrefix + fieldName + parametersLevelSuffix, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

func (r *ResourceAssert) ToTerraformTestCheckFunc(t *testing.T) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		var result []error

		for i, a := range r.assertions {
			switch a.resourceAssertionType {
			case resourceAssertionTypeValueSet:
				if err := resource.TestCheckResourceAttr(r.name, a.fieldName, a.expectedValue)(s); err != nil {
					errCut, _ := strings.CutPrefix(err.Error(), fmt.Sprintf("%s: ", r.name))
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %s", r.name, r.prefix, i+1, len(r.assertions), errCut))
				}
			case resourceAssertionTypeValueNotSet:
				if err := resource.TestCheckNoResourceAttr(r.name, a.fieldName)(s); err != nil {
					errCut, _ := strings.CutPrefix(err.Error(), fmt.Sprintf("%s: ", r.name))
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %s", r.name, r.prefix, i+1, len(r.assertions), errCut))
				}
			}
		}

		return errors.Join(result...)
	}
}

func (r *ResourceAssert) ToTerraformImportStateCheckFunc(t *testing.T) resource.ImportStateCheckFunc {
	t.Helper()
	return func(s []*terraform.InstanceState) error {
		var result []error

		for i, a := range r.assertions {
			switch a.resourceAssertionType {
			case resourceAssertionTypeValueSet:
				if err := importchecks.TestCheckResourceAttrInstanceState(r.id, a.fieldName, a.expectedValue)(s); err != nil {
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %w", r.id, r.prefix, i+1, len(r.assertions), err))
				}
			case resourceAssertionTypeValueNotSet:
				panic("implement")
			}
		}

		return errors.Join(result...)
	}
}
