package testint

import (
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: Generate assertions
// TODO: Add newly added fields to all options tests

func TestInt_Tasks(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	sql := "SELECT CURRENT_TIMESTAMP"

	assertTask := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, warehouseName string) {
		t.Helper()
		assert.Equal(t, id, task.ID())
		assert.NotEmpty(t, task.CreatedOn)
		assert.Equal(t, id.Name(), task.Name)
		assert.NotEmpty(t, task.Id)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), task.DatabaseName)
		assert.Equal(t, testClientHelper().Ids.SchemaId().Name(), task.SchemaName)
		assert.Equal(t, "ACCOUNTADMIN", task.Owner)
		assert.Equal(t, "", task.Comment)
		assert.Equal(t, warehouseName, task.Warehouse)
		assert.Equal(t, "", task.Schedule)
		assert.Empty(t, task.Predecessors)
		assert.Equal(t, sdk.TaskStateSuspended, task.State)
		assert.Equal(t, sql, task.Definition)
		assert.Equal(t, "", task.Condition)
		assert.Equal(t, false, task.AllowOverlappingExecution)
		assert.Empty(t, task.ErrorIntegration)
		assert.Empty(t, task.LastCommittedOn)
		assert.Empty(t, task.LastSuspendedOn)
		assert.Equal(t, "ROLE", task.OwnerRoleType)
		assert.Empty(t, task.Config)
		assert.Empty(t, task.Budget)
	}

	assertTaskWithOptions := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, comment string, warehouse string, schedule string, condition string, allowOverlappingExecution bool, config string, predecessor *sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.Equal(t, id, task.ID())
		assert.NotEmpty(t, task.CreatedOn)
		assert.Equal(t, id.Name(), task.Name)
		assert.NotEmpty(t, task.Id)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), task.DatabaseName)
		assert.Equal(t, testClientHelper().Ids.SchemaId().Name(), task.SchemaName)
		assert.Equal(t, "ACCOUNTADMIN", task.Owner)
		assert.Equal(t, comment, task.Comment)
		assert.Equal(t, warehouse, task.Warehouse)
		assert.Equal(t, schedule, task.Schedule)
		assert.Equal(t, sdk.TaskStateSuspended, task.State)
		assert.Equal(t, sql, task.Definition)
		assert.Equal(t, condition, task.Condition)
		assert.Equal(t, allowOverlappingExecution, task.AllowOverlappingExecution)
		assert.Empty(t, task.ErrorIntegration)
		assert.Empty(t, task.LastCommittedOn)
		assert.Empty(t, task.LastSuspendedOn)
		assert.Equal(t, "ROLE", task.OwnerRoleType)
		assert.Equal(t, config, task.Config)
		assert.Empty(t, task.Budget)
		if predecessor != nil {
			assert.Len(t, task.Predecessors, 1)
			assert.Contains(t, task.Predecessors, *predecessor)
		} else {
			assert.Empty(t, task.Predecessors)
		}
	}

	assertTaskTerse := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, schedule string) {
		t.Helper()
		assert.Equal(t, id, task.ID())
		assert.NotEmpty(t, task.CreatedOn)
		assert.Equal(t, id.Name(), task.Name)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), task.DatabaseName)
		assert.Equal(t, testClientHelper().Ids.SchemaId().Name(), task.SchemaName)
		assert.Equal(t, schedule, task.Schedule)

		// all below are not contained in the terse response, that's why all of them we expect to be empty
		assert.Empty(t, task.Id)
		assert.Empty(t, task.Owner)
		assert.Empty(t, task.Comment)
		assert.Empty(t, task.Warehouse)
		assert.Empty(t, task.Predecessors)
		assert.Empty(t, task.State)
		assert.Empty(t, task.Definition)
		assert.Empty(t, task.Condition)
		assert.Empty(t, task.AllowOverlappingExecution)
		assert.Empty(t, task.ErrorIntegration)
		assert.Empty(t, task.LastCommittedOn)
		assert.Empty(t, task.LastSuspendedOn)
		assert.Empty(t, task.OwnerRoleType)
		assert.Empty(t, task.Config)
		assert.Empty(t, task.Budget)
	}

	t.Run("create task: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropTaskFunc(t, id))

		task, err := testClientHelper().Task.Show(t, id)
		require.NoError(t, err)

		assertTask(t, task, id, "")
	})

	t.Run("create task: with initial warehouse", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql).WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall)))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropTaskFunc(t, id))

		task, err := testClientHelper().Task.Show(t, id)
		require.NoError(t, err)

		assertTask(t, task, id, "")
	})

	t.Run("create task: almost complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql).
			WithOrReplace(true).
			WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(testClientHelper().Ids.WarehouseId())).
			WithSchedule("10 MINUTE").
			WithConfig(`$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`).
			WithAllowOverlappingExecution(true).
			WithSessionParameters(sdk.SessionParameters{
				JSONIndent: sdk.Int(4),
			}).
			WithUserTaskTimeoutMs(500).
			WithSuspendTaskAfterNumFailures(3).
			WithComment("some comment").
			WithWhen(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropTaskFunc(t, id))

		task, err := testClientHelper().Task.Show(t, id)
		require.NoError(t, err)

		assertTaskWithOptions(t, task, id, "some comment", testClientHelper().Ids.WarehouseId().Name(), "10 MINUTE", `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`, true, `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`, nil)
	})

	t.Run("create task: with after", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(rootTaskId, sql))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropTaskFunc(t, rootTaskId))

		err = testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql).WithAfter([]sdk.SchemaObjectIdentifier{rootTaskId}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropTaskFunc(t, id))

		task, err := testClientHelper().Task.Show(t, id)
		require.NoError(t, err)

		assertTaskWithOptions(t, task, id, "", "", "", "", false, "", &rootTaskId)
	})

	t.Run("create dag of tasks", func(t *testing.T) {
		rootId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		root, rootCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootId, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(rootCleanup)
		require.Empty(t, root.Predecessors)

		t1, t1Cleanup := testClientHelper().Task.CreateWithAfter(t, rootId)
		t.Cleanup(t1Cleanup)
		require.Equal(t, []sdk.SchemaObjectIdentifier{rootId}, t1.Predecessors)

		t2, t2Cleanup := testClientHelper().Task.CreateWithAfter(t, t1.ID(), rootId)
		t.Cleanup(t2Cleanup)

		require.Contains(t, t2.Predecessors, rootId)
		require.Contains(t, t2.Predecessors, t1.ID())
		require.Len(t, t2.Predecessors, 2)

		t3, t3Cleanup := testClientHelper().Task.CreateWithAfter(t, t2.ID(), t1.ID())
		t.Cleanup(t3Cleanup)

		require.Contains(t, t3.Predecessors, t2.ID())
		require.Contains(t, t3.Predecessors, t1.ID())
		require.Len(t, t3.Predecessors, 2)

		rootTasks, err := sdk.GetRootTasks(client.Tasks, ctx, rootId)
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t1.ID())
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t2.ID())
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t3.ID())
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		// cannot set ALLOW_OVERLAPPING_EXECUTION on child task
		alterRequest := sdk.NewAlterTaskRequest(t1.ID()).WithSet(*sdk.NewTaskSetRequest().WithAllowOverlappingExecution(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.ErrorContains(t, err, "Cannot set allow_overlapping_execution on non-root task")

		// can set ALLOW_OVERLAPPING_EXECUTION on root task
		alterRequest = sdk.NewAlterTaskRequest(rootId).WithSet(*sdk.NewTaskSetRequest().WithAllowOverlappingExecution(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		// can create cycle, because DAG is suspended
		alterRequest = sdk.NewAlterTaskRequest(t1.ID()).WithAddAfter([]sdk.SchemaObjectIdentifier{t3.ID()})
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		// can get the root task even with cycle
		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t3.ID())
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		// we get an error when trying to start
		alterRequest = sdk.NewAlterTaskRequest(rootId).WithResume(true)
		err = client.Tasks.Alter(ctx, alterRequest)
		require.ErrorContains(t, err, "Graph has at least one cycle containing task")
	})

	t.Run("create dag of tasks - multiple roots", func(t *testing.T) {
		root1Id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		root1, root1Cleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(root1Id, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(root1Cleanup)
		require.Empty(t, root1.Predecessors)

		root2Id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		root2, root2Cleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(root2Id, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(root2Cleanup)
		require.Empty(t, root2.Predecessors)

		t1, t1Cleanup := testClientHelper().Task.CreateWithAfter(t, root1.ID(), root2.ID())
		t.Cleanup(t1Cleanup)

		require.Contains(t, t1.Predecessors, root1Id)
		require.Contains(t, t1.Predecessors, root2Id)
		require.Len(t, t1.Predecessors, 2)

		rootTasks, err := sdk.GetRootTasks(client.Tasks, ctx, t1.ID())
		require.NoError(t, err)
		require.Len(t, rootTasks, 2)
		require.Contains(t, []sdk.SchemaObjectIdentifier{root1Id, root2Id}, rootTasks[0].ID())
		require.Contains(t, []sdk.SchemaObjectIdentifier{root1Id, root2Id}, rootTasks[1].ID())

		// we get an error when trying to start
		alterRequest := sdk.NewAlterTaskRequest(root1Id).WithResume(true)
		err = client.Tasks.Alter(ctx, alterRequest)
		require.ErrorContains(t, err, "The graph has more than one root task (one without predecessors)")
	})

	t.Run("create task: with tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		task, taskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id, sql).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: "v1",
				},
			}),
		)
		t.Cleanup(taskCleanup)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), task.ID(), sdk.ObjectTypeTask)
		require.NoError(t, err)

		assert.Equal(t, "v1", returnedTagValue)
	})

	t.Run("clone task: default", func(t *testing.T) {
		sourceTask, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tasks.Clone(ctx, sdk.NewCloneTaskRequest(id, sourceTask.ID()))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropTaskFunc(t, id))

		task, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTask(t, task, id, testClientHelper().Ids.WarehouseId().Name())
	})

	t.Run("drop task: existing", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		err := client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(task.ID()))
		require.NoError(t, err)

		_, err = client.Tasks.ShowByID(ctx, task.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop task: non-existing", func(t *testing.T) {
		err := client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter task: set value and unset value", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		alterRequest := sdk.NewAlterTaskRequest(task.ID()).WithSet(*sdk.NewTaskSetRequest().WithComment("new comment").WithUserTaskTimeoutMs(1000))
		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, task.ID())
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredTask.Comment)

		alterRequest = sdk.NewAlterTaskRequest(task.ID()).WithUnset(*sdk.NewTaskUnsetRequest().WithComment(true).WithUserTaskTimeoutMs(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, task.ID())
		require.NoError(t, err)

		assert.Equal(t, "", alteredTask.Comment)
	})

	t.Run("alter task: set and unset tag", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		tagValue := "abc"
		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSetTags([]sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}))
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), task.ID(), sdk.ObjectTypeTask)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithUnsetTags([]sdk.ObjectIdentifier{
			tag.ID(),
		}))
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), task.ID(), sdk.ObjectTypeTask)
		require.Error(t, err)
	})

	t.Run("alter task: resume and suspend", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		task, taskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(taskCleanup)

		assert.Equal(t, sdk.TaskStateSuspended, task.State)

		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithResume(true))
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, sdk.TaskStateStarted, alteredTask.State)

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithSuspend(true))
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, sdk.TaskStateSuspended, alteredTask.State)
	})

	t.Run("alter task: remove after and add after", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		rootTask, rootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(rootTaskCleanup)

		task, taskCleanup := testClientHelper().Task.CreateWithAfter(t, rootTask.ID())
		t.Cleanup(taskCleanup)

		assert.Contains(t, task.Predecessors, rootTask.ID())

		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithRemoveAfter([]sdk.SchemaObjectIdentifier{rootTask.ID()}))
		require.NoError(t, err)

		task, err = client.Tasks.ShowByID(ctx, task.ID())

		require.NoError(t, err)
		assert.Empty(t, task.Predecessors)

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithAddAfter([]sdk.SchemaObjectIdentifier{rootTask.ID()}))
		require.NoError(t, err)

		task, err = client.Tasks.ShowByID(ctx, task.ID())

		require.NoError(t, err)
		assert.Contains(t, task.Predecessors, rootTask.ID())
	})

	t.Run("alter task: modify when and as", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		newSql := "SELECT CURRENT_DATE"
		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithModifyAs(newSql))
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, task.ID())
		require.NoError(t, err)

		assert.Equal(t, newSql, alteredTask.Definition)

		newWhen := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithModifyWhen(newWhen))
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, task.ID())
		require.NoError(t, err)

		assert.Equal(t, newWhen, alteredTask.Condition)
	})

	// TODO: Change this test (the search is too broad)
	t.Run("show task: default", func(t *testing.T) {
		task1, task1Cleanup := testClientHelper().Task.Create(t)
		t.Cleanup(task1Cleanup)

		task2, task2Cleanup := testClientHelper().Task.Create(t)
		t.Cleanup(task2Cleanup)

		returnedTasks, err := client.Tasks.Show(ctx, sdk.NewShowTaskRequest())
		require.NoError(t, err)

		assert.LessOrEqual(t, 2, len(returnedTasks))
		assert.Contains(t, returnedTasks, *task1)
		assert.Contains(t, returnedTasks, *task2)
	})

	// TODO: Change this test (the search is too broad)
	t.Run("show task: terse", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		task, taskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(taskCleanup)

		returnedTasks, err := client.Tasks.Show(ctx, sdk.NewShowTaskRequest().WithTerse(true))
		require.NoError(t, err)

		assert.LessOrEqual(t, 1, len(returnedTasks))
		assertTaskTerse(t, &returnedTasks[0], task.ID(), "10 MINUTE")
	})

	t.Run("show task: with options", func(t *testing.T) {
		task1, task1Cleanup := testClientHelper().Task.Create(t)
		t.Cleanup(task1Cleanup)

		task2, task2Cleanup := testClientHelper().Task.Create(t)
		t.Cleanup(task2Cleanup)

		showRequest := sdk.NewShowTaskRequest().
			WithLike(sdk.Like{Pattern: &task1.Name}).
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(5)})
		returnedTasks, err := client.Tasks.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedTasks))
		assert.Contains(t, returnedTasks, *task1)
		assert.NotContains(t, returnedTasks, *task2)
	})

	t.Run("describe task: default", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		returnedTask, err := client.Tasks.Describe(ctx, task.ID())
		require.NoError(t, err)

		assertTask(t, returnedTask, task.ID(), testClientHelper().Ids.WarehouseId().Name())
	})

	t.Run("execute task: default", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		executeRequest := sdk.NewExecuteTaskRequest(task.ID())
		err := client.Tasks.Execute(ctx, executeRequest)
		require.NoError(t, err)
	})

	t.Run("temporarily suspend root tasks", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		rootTask, rootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, sql).WithSchedule("60 MINUTES"))
		t.Cleanup(rootTaskCleanup)

		task, taskCleanup := testClientHelper().Task.CreateWithAfter(t, rootTask.ID())
		t.Cleanup(taskCleanup)

		require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithResume(true)))
		t.Cleanup(func() {
			require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithSuspend(true)))
		})

		tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, task.ID(), task.ID())
		require.NoError(t, err)
		require.NotEmpty(t, tasksToResume)

		rootTaskStatus, err := client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateSuspended, rootTaskStatus.State)

		require.NoError(t, client.Tasks.ResumeTasks(ctx, tasksToResume))

		rootTaskStatus, err = client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, rootTaskStatus.State)
	})

	t.Run("resume root tasks within a graph containing more than one root task", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		rootTask, rootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, sql).WithSchedule("60 MINUTES"))
		t.Cleanup(rootTaskCleanup)

		secondRootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secondRootTask, secondRootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(secondRootTaskId, sql).WithSchedule("60 MINUTES"))
		t.Cleanup(secondRootTaskCleanup)

		_, cleanupTask := testClientHelper().Task.CreateWithAfter(t, rootTask.ID(), secondRootTask.ID())
		t.Cleanup(cleanupTask)

		require.ErrorContains(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithResume(true)), "The graph has more than one root task (one without predecessors)")
		require.ErrorContains(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(secondRootTask.ID()).WithResume(true)), "The graph has more than one root task (one without predecessors)")
	})

	t.Run("suspend root tasks temporarily with three sequentially connected tasks - last in DAG", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		rootTask, rootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, sql).WithSchedule("60 MINUTES"))
		t.Cleanup(rootTaskCleanup)

		middleTask, middleTaskCleanup := testClientHelper().Task.CreateWithAfter(t, rootTask.ID())
		t.Cleanup(middleTaskCleanup)

		task, taskCleanup := testClientHelper().Task.CreateWithAfter(t, middleTask.ID())
		t.Cleanup(taskCleanup)

		require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(middleTask.ID()).WithResume(true)))
		t.Cleanup(func() {
			require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(middleTask.ID()).WithSuspend(true)))
		})

		require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithResume(true)))
		t.Cleanup(func() {
			require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithSuspend(true)))
		})

		tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, task.ID(), task.ID())
		require.NoError(t, err)
		require.NotEmpty(t, tasksToResume)
		require.Contains(t, tasksToResume, rootTask.ID())

		rootTaskStatus, err := client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateSuspended, rootTaskStatus.State)

		middleTaskStatus, err := client.Tasks.ShowByID(ctx, middleTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, middleTaskStatus.State)

		require.NoError(t, client.Tasks.ResumeTasks(ctx, tasksToResume))

		rootTaskStatus, err = client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, rootTaskStatus.State)

		middleTaskStatus, err = client.Tasks.ShowByID(ctx, middleTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, middleTaskStatus.State)
	})

	t.Run("suspend root tasks temporarily with three sequentially connected tasks - middle in DAG", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		rootTask, rootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, sql).WithSchedule("60 MINUTES"))
		t.Cleanup(rootTaskCleanup)

		middleTask, middleTaskCleanup := testClientHelper().Task.CreateWithAfter(t, rootTask.ID())
		t.Cleanup(middleTaskCleanup)

		childTask, childTaskCleanup := testClientHelper().Task.CreateWithAfter(t, middleTask.ID())
		t.Cleanup(childTaskCleanup)

		require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(childTask.ID()).WithResume(true)))
		t.Cleanup(func() {
			require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(childTask.ID()).WithSuspend(true)))
		})

		require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithResume(true)))
		t.Cleanup(func() {
			require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithSuspend(true)))
		})

		tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, middleTask.ID(), middleTask.ID())
		require.NoError(t, err)
		require.NotEmpty(t, tasksToResume)

		rootTaskStatus, err := client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateSuspended, rootTaskStatus.State)

		childTaskStatus, err := client.Tasks.ShowByID(ctx, childTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, childTaskStatus.State)

		require.NoError(t, client.Tasks.ResumeTasks(ctx, tasksToResume))

		rootTaskStatus, err = client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, rootTaskStatus.State)

		childTaskStatus, err = client.Tasks.ShowByID(ctx, childTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, childTaskStatus.State)
	})

	// TODO(SNOW-1277135): Create more tests with different sets of roots/children and see if the current implementation
	// acts correctly in certain situations/edge cases.
}

func TestInt_TasksShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupTaskHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createTaskHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		err := client.Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, "SELECT CURRENT_TIMESTAMP"))
		require.NoError(t, err)
		t.Cleanup(cleanupTaskHandle(id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createTaskHandle(t, id1)
		createTaskHandle(t, id2)

		e1, err := client.Tasks.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Tasks.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
