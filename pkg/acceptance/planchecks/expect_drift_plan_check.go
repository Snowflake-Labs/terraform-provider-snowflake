package planchecks

import (
	"context"
	"errors"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"slices"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = expectDriftPlanCheck{}

type expectDriftPlanCheck struct {
	resourceAddress string
	attribute       string
	oldValue        *string
	newValue        *string
}

// TODO: test
// TODO: extract common logic with expectChangePlanCheck
func (e expectDriftPlanCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result []error
	var resourceFound bool

	for _, change := range req.Plan.ResourceDrift {
		if e.resourceAddress != change.Address {
			continue
		}
		resourceFound = true

		var before, after map[string]any
		if change.Change.Before != nil {
			before = change.Change.Before.(map[string]any)
		}
		if change.Change.After != nil {
			after = change.Change.After.(map[string]any)
		}
		valueBefore, valueBeforeOk := before[e.attribute]
		valueAfter, valueAfterOk := after[e.attribute]

		if e.oldValue == nil && !(valueBefore == nil || valueBefore == "") {
			result = append(result, fmt.Errorf("expect drift: attribute %s before expected to be empty, got: %s", e.attribute, valueBefore))
		}
		if e.newValue == nil && !(valueAfter == nil || valueAfter == "") {
			result = append(result, fmt.Errorf("expect drift: attribute %s after expected to be empty, got: %s", e.attribute, valueAfter))
		}

		if e.oldValue != nil {
			if !valueBeforeOk {
				result = append(result, fmt.Errorf("expect drift: attribute %s before expected to be %s, got empty", e.attribute, *e.oldValue))
			} else if *e.oldValue != valueBefore {
				result = append(result, fmt.Errorf("expect drift: attribute %s before expected to be %s, got %s", e.attribute, *e.oldValue, valueBefore))
			}
		}
		if e.newValue != nil {
			if !valueAfterOk {
				result = append(result, fmt.Errorf("expect drift: attribute %s after expected to be %s, got empty", e.attribute, *e.newValue))
			} else if *e.newValue != valueAfter {
				result = append(result, fmt.Errorf("expect drift: attribute %s after expected to be %s, got %s", e.attribute, *e.newValue, valueAfter))
			}
		}

		if !slices.Contains(change.Change.Actions, tfjson.ActionUpdate) {
			result = append(result, fmt.Errorf("expect drift: expected action %s for %s, got: %v", tfjson.ActionUpdate, e.resourceAddress, change.Change.Actions))
		}
	}

	if !resourceFound {
		result = append(result, fmt.Errorf("expect drift: no resource drift found for %s", e.resourceAddress))
	}

	resp.Error = errors.Join(result...)
}

// TODO: describe
func ExpectDrift(resourceAddress string, attribute string, oldValue *string, newValue *string) plancheck.PlanCheck {
	return expectDriftPlanCheck{
		resourceAddress,
		attribute,
		oldValue,
		newValue,
	}
}
