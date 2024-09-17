package objectassert

import (
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"reflect"
	"slices"
	"testing"
)

func (w *TaskAssert) HasNotEmptyCreatedOn() *TaskAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created on not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return w
}

func (w *TaskAssert) HasNotEmptyId() *TaskAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if o.Id == "" {
			return fmt.Errorf("expected id not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return w
}

func (w *TaskAssert) HasPredecessors(ids ...sdk.SchemaObjectIdentifier) *TaskAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Task) error {
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
	return w
}

func (t *TaskAssert) HasTaskRelations(expected sdk.TaskRelations) *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if slices.EqualFunc(o.TaskRelations.Predecessors, expected.Predecessors, func(id sdk.SchemaObjectIdentifier, id2 sdk.SchemaObjectIdentifier) bool {
			return id.FullyQualifiedName() == id2.FullyQualifiedName()
		}) {
			return fmt.Errorf("expected task predecessors: %v; got: %v", expected.Predecessors, o.TaskRelations.Predecessors)
		}

		if !reflect.DeepEqual(expected.FinalizerTask, o.TaskRelations.FinalizerTask) {
			return fmt.Errorf("expected finalizer task: %v; got: %v", expected.FinalizerTask, o.TaskRelations.FinalizerTask)
		}
		return nil
	})
	return t
}
