package objectassert

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (t *TaskAssert) HasNotEmptyCreatedOn() *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created on not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return t
}

func (t *TaskAssert) HasNotEmptyId() *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if o.Id == "" {
			return fmt.Errorf("expected id not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return t
}

func (t *TaskAssert) HasPredecessorsInAnyOrder(ids ...sdk.SchemaObjectIdentifier) *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if !assert.ElementsMatch(t, ids, o.Predecessors) {
			return fmt.Errorf("expected %v predecessors in task relations, got %v", ids, o.TaskRelations.Predecessors)
		}
		return nil
	})
	return t
}

func (t *TaskAssert) HasTaskRelations(expected sdk.TaskRelations) *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		errs := make([]error, 0)
		if !assert.ElementsMatch(t, o.TaskRelations.Predecessors, expected.Predecessors) {
			errs = append(errs, fmt.Errorf("expected %v predecessors in task relations, got %v", expected.Predecessors, o.TaskRelations.Predecessors))
		}
		if !reflect.DeepEqual(expected.FinalizerTask, o.TaskRelations.FinalizerTask) {
			errs = append(errs, fmt.Errorf("expected finalizer task: %v; got: %v", expected.FinalizerTask, o.TaskRelations.FinalizerTask))
		}
		if expected.FinalizedRootTask != nil {
			// This is not supported because we would have to traverse the task graph to find the root task.
			errs = append(errs, fmt.Errorf("asserting FinalizedRootTask is not supported"))
		}
		return errors.Join(errs...)
	})
	return t
}

func (t *TaskAssert) HasWarehouse(expected *sdk.AccountObjectIdentifier) *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if o.Warehouse == nil && expected != nil {
			return fmt.Errorf("expected warehouse to have value; got: nil")
		}
		if o.Warehouse != nil && expected == nil {
			return fmt.Errorf("expected warehouse to no have value; got: %s", o.Warehouse.Name())
		}
		if o.Warehouse != nil && expected != nil && o.Warehouse.Name() != expected.Name() {
			return fmt.Errorf("expected warehouse: %v; got: %v", expected.Name(), o.Warehouse.Name())
		}
		return nil
	})
	return t
}

func (t *TaskAssert) HasErrorIntegration(expected *sdk.AccountObjectIdentifier) *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if o.ErrorIntegration == nil && expected != nil {
			return fmt.Errorf("expected error integration to have value; got: nil")
		}
		if o.ErrorIntegration != nil && expected == nil {
			return fmt.Errorf("expected error integration to have no value; got: %s", o.ErrorIntegration.Name())
		}
		if o.ErrorIntegration != nil && expected != nil && o.ErrorIntegration.Name() != expected.Name() {
			return fmt.Errorf("expected error integration: %v; got: %v", expected.Name(), o.ErrorIntegration.Name())
		}
		return nil
	})
	return t
}
