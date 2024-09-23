package objectassert

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"testing"

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

func (t *TaskAssert) HasPredecessors(ids ...sdk.SchemaObjectIdentifier) *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if len(o.Predecessors) != len(ids) {
			return fmt.Errorf("expected %d (%v) predecessors, got %d (%v)", len(ids), ids, len(o.Predecessors), o.Predecessors)
		}
		var errs []error
		for _, id := range ids {
			if !slices.ContainsFunc(o.Predecessors, func(predecessorId sdk.SchemaObjectIdentifier) bool {
				return predecessorId.FullyQualifiedName() == id.FullyQualifiedName()
			}) {
				errs = append(errs, fmt.Errorf("expected id: %s, to be in the list of predecessors: %v", id.FullyQualifiedName(), o.Predecessors))
			}
		}
		return errors.Join(errs...)
	})
	return t
}

func (t *TaskAssert) HasTaskRelations(expected sdk.TaskRelations) *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if len(o.TaskRelations.Predecessors) != len(expected.Predecessors) {
			return fmt.Errorf("expected %d (%v) predecessors in task relations, got %d (%v)", len(expected.Predecessors), expected.Predecessors, len(o.TaskRelations.Predecessors), o.TaskRelations.Predecessors)
		}
		var errs []error
		for _, id := range expected.Predecessors {
			if !slices.ContainsFunc(o.TaskRelations.Predecessors, func(predecessorId sdk.SchemaObjectIdentifier) bool {
				return predecessorId.FullyQualifiedName() == id.FullyQualifiedName()
			}) {
				errs = append(errs, fmt.Errorf("expected id: %s, to be in the list of predecessors in task relations: %v", id.FullyQualifiedName(), o.TaskRelations.Predecessors))
			}
		}
		if !reflect.DeepEqual(expected.FinalizerTask, o.TaskRelations.FinalizerTask) {
			errs = append(errs, fmt.Errorf("expected finalizer task: %v; got: %v", expected.FinalizerTask, o.TaskRelations.FinalizerTask))
		}
		return errors.Join(errs...)
	})
	return t
}
