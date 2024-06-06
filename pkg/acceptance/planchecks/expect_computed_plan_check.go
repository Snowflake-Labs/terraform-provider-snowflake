package planchecks

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = expectComputedPlanCheck{}

type expectComputedPlanCheck struct {
	resourceAddress string
	attribute       string
	expectComputed  bool
}

// TODO: test
func (e expectComputedPlanCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result []error
	var resourceFound bool

	for _, change := range req.Plan.ResourceChanges {
		if e.resourceAddress != change.Address {
			continue
		}
		resourceFound = true

		var computed map[string]any
		if change.Change.AfterUnknown != nil {
			computed = change.Change.AfterUnknown.(map[string]any)
		}
		_, isComputed := computed[e.attribute]

		if e.expectComputed && !isComputed {
			result = append(result, fmt.Errorf("expect computed: attribute %s expected to be computed", e.attribute))
		}
		if !e.expectComputed && isComputed {
			result = append(result, fmt.Errorf("expect computed: attribute %s expected not to be computed", e.attribute))
		}
	}

	if !resourceFound {
		result = append(result, fmt.Errorf("expect computed: no changes found for %s", e.resourceAddress))
	}

	resp.Error = errors.Join(result...)
}

// TODO: describe
func ExpectComputed(resourceAddress string, attribute string, expectComputed bool) plancheck.PlanCheck {
	return expectComputedPlanCheck{
		resourceAddress,
		attribute,
		expectComputed,
	}
}
