package sdk

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type testTasks struct {
	tasks

	stubbedTasks map[string]*Task
}

func (v *testTasks) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Task, error) {
	t, ok := v.stubbedTasks[id.Name()]
	if !ok {
		return nil, errors.New("no task configured, check test config")
	}
	return t, nil
}

func TestTasks_GetRootTasks(t *testing.T) {
	db := "database"
	sc := "schema"

	setUpTasks := func(p map[string][]string) map[string]*Task {
		r := make(map[string]*Task)
		for k, v := range p {
			if k == "initial" || k == "expected" {
				continue
			}
			t := Task{DatabaseName: db, SchemaName: sc, Name: k}
			predecessors := make([]SchemaObjectIdentifier, len(v))
			for i, predecessor := range v {
				predecessors[i] = NewSchemaObjectIdentifier(db, sc, predecessor)
			}
			t.Predecessors = predecessors
			r[k] = &t
		}
		return r
	}

	testParameters := []map[string][]string{
		{"t1": {}, "initial": {"t1"}, "expected": {"t1"}},
		{"t1": {"t2"}, "t2": {"t3"}, "t3": {}, "initial": {"t1"}, "expected": {"t3"}},
		{"t1": {"t2", "t3"}, "t2": {"t3"}, "t3": {}, "initial": {"t1"}, "expected": {"t3"}},
		{"t1": {"t2", "t3"}, "t2": {}, "t3": {}, "initial": {"t1"}, "expected": {"t2", "t3"}},
		{"t1": {}, "t2": {}, "initial": {"t1"}, "expected": {"t1"}},
		{"t1": {"t2", "t3", "t4"}, "t2": {}, "t3": {}, "t4": {}, "initial": {"t1"}, "expected": {"t2", "t3", "t4"}},
		{"t1": {"t2", "t3", "t4"}, "t2": {}, "t3": {"t2"}, "t4": {"t3"}, "initial": {"t1"}, "expected": {"t2"}},
	}
	for i, p := range testParameters {
		t.Run(fmt.Sprintf("test case [%v]", i), func(t *testing.T) {
			ctx := context.Background()
			initial, ok := p["initial"]
			if !ok {
				t.FailNow()
			}
			expected, ok := p["expected"]
			if !ok {
				t.FailNow()
			}
			client := new(testTasks)
			client.stubbedTasks = setUpTasks(p)

			rootTasks, err := GetRootTasks(client, ctx, NewSchemaObjectIdentifier(db, sc, initial[0]))
			require.NoError(t, err)
			for _, v := range rootTasks {
				assert.Contains(t, expected, v.Name)
			}
			require.Len(t, rootTasks, len(expected))
		})
	}
}
