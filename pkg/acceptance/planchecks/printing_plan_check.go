package planchecks

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = printingPlanCheck{}

type printingPlanCheck struct {
	resourceAddress string
	attributes      []string
}

// TODO: test
// TODO: add traversal
func (e printingPlanCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result []error

	for _, change := range req.Plan.ResourceDrift {
		if e.resourceAddress != change.Address {
			continue
		}
		actions := change.Change.Actions
		var before, after, computed map[string]any
		if change.Change.Before != nil {
			before = change.Change.Before.(map[string]any)
		}
		if change.Change.After != nil {
			after = change.Change.After.(map[string]any)
		}
		if change.Change.AfterUnknown != nil {
			computed = change.Change.AfterUnknown.(map[string]any)
		}
		fmt.Printf("resource drift for [%s]: actions: %v\n", change.Address, actions)
		for _, attr := range e.attributes {
			valueBefore := before[attr]
			valueAfter := after[attr]
			_, isComputed := computed[attr]
			fmt.Printf("\t[%s]: before: %v, after: %v, computed: %t\n", attr, valueBefore, valueAfter, isComputed)
		}
	}

	for _, change := range req.Plan.ResourceChanges {
		if e.resourceAddress != change.Address {
			continue
		}
		actions := change.Change.Actions
		var before, after, computed map[string]any
		if change.Change.Before != nil {
			before = change.Change.Before.(map[string]any)
		}
		if change.Change.After != nil {
			after = change.Change.After.(map[string]any)
		}
		if change.Change.AfterUnknown != nil {
			computed = change.Change.AfterUnknown.(map[string]any)
		}
		fmt.Printf("resource change for [%s]: actions: %v\n", change.Address, actions)
		for _, attr := range e.attributes {
			valueBefore := before[attr]
			valueAfter := after[attr]
			_, isComputed := computed[attr]
			fmt.Printf("\t[%s]: before: %v, after: %v, computed: %t\n", attr, valueBefore, valueAfter, isComputed)
		}
	}

	resp.Error = errors.Join(result...)
}

func PrintPlanDetails(resourceAddress string, attributes ...string) plancheck.PlanCheck {
	return printingPlanCheck{
		resourceAddress,
		attributes,
	}
}
