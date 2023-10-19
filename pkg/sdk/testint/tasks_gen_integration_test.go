package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Tasks(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	sql := "SELECT CURRENT_TIMESTAMP"

	assertTask := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.Equal(t, id, task.ID())
		assert.NotEmpty(t, task.CreatedOn)
		assert.Equal(t, id.Name(), task.Name)
		assert.NotEmpty(t, task.Id)
		assert.Equal(t, testDb(t).Name, task.DatabaseName)
		assert.Equal(t, testSchema(t).Name, task.SchemaName)
		assert.Equal(t, "ACCOUNTADMIN", task.Owner)
		assert.Equal(t, "", task.Comment)
		assert.Equal(t, "", task.Warehouse)
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
		assert.Equal(t, testDb(t).Name, task.DatabaseName)
		assert.Equal(t, testSchema(t).Name, task.SchemaName)
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
		assert.Equal(t, testDb(t).Name, task.DatabaseName)
		assert.Equal(t, testSchema(t).Name, task.SchemaName)
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

	cleanupTaskProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(id))
			require.NoError(t, err)
		}
	}

	createTaskBasicRequest := func(t *testing.T) *sdk.CreateTaskRequest {
		t.Helper()
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		return sdk.NewCreateTaskRequest(id, sql)
	}

	createTaskWithRequest := func(t *testing.T, request *sdk.CreateTaskRequest) *sdk.Task {
		t.Helper()
		id := request.GetName()

		err := client.Tasks.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		task, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		return task
	}

	createTask := func(t *testing.T) *sdk.Task {
		t.Helper()
		return createTaskWithRequest(t, createTaskBasicRequest(t))
	}

	t.Run("create task: no optionals", func(t *testing.T) {
		request := createTaskBasicRequest(t)

		task := createTaskWithRequest(t, request)

		assertTask(t, task, request.GetName())
	})

	t.Run("create task: with initial warehouse", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithWarehouse(sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(&sdk.WarehouseSizeXSmall))

		task := createTaskWithRequest(t, request)

		assertTask(t, task, request.GetName())
	})

	t.Run("create task: almost complete case", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		request := createTaskBasicRequest(t).
			WithOrReplace(sdk.Bool(true)).
			WithWarehouse(sdk.NewCreateTaskWarehouseRequest().WithWarehouse(sdk.Pointer(warehouse.ID()))).
			WithSchedule(sdk.String("10 MINUTE")).
			WithConfig(sdk.String(`$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`)).
			WithAllowOverlappingExecution(sdk.Bool(true)).
			WithSessionParameters(&sdk.SessionParameters{
				JSONIndent: sdk.Int(4),
			}).
			WithUserTaskTimeoutMs(sdk.Int(500)).
			WithSuspendTaskAfterNumFailures(sdk.Int(3)).
			WithComment(sdk.String("some comment")).
			WithWhen(sdk.String(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`))
		id := request.GetName()

		task := createTaskWithRequest(t, request)

		assertTaskWithOptions(t, task, id, "some comment", warehouse.Name, "10 MINUTE", `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`, true, `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`, nil)
	})

	t.Run("create task: with after", func(t *testing.T) {
		otherName := random.String()
		otherId := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, otherName)

		request := sdk.NewCreateTaskRequest(otherId, sql).WithSchedule(sdk.String("10 MINUTE"))

		createTaskWithRequest(t, request)

		request = createTaskBasicRequest(t).
			WithAfter([]sdk.SchemaObjectIdentifier{otherId})

		task := createTaskWithRequest(t, request)

		assertTaskWithOptions(t, task, request.GetName(), "", "", "", "", false, "", &otherId)
	})

	t.Run("create dag of tasks", func(t *testing.T) {
		rootName := random.String()
		rootId := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, rootName)

		request := sdk.NewCreateTaskRequest(rootId, sql).WithSchedule(sdk.String("10 MINUTE"))
		root := createTaskWithRequest(t, request)

		require.Equal(t, []sdk.SchemaObjectIdentifier{}, root.Predecessors)

		t1Name := random.String()
		t1Id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, t1Name)

		request = sdk.NewCreateTaskRequest(t1Id, sql).WithAfter([]sdk.SchemaObjectIdentifier{rootId})
		t1 := createTaskWithRequest(t, request)

		require.Equal(t, []sdk.SchemaObjectIdentifier{rootId}, t1.Predecessors)

		t2Name := random.String()
		t2Id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, t2Name)

		request = sdk.NewCreateTaskRequest(t2Id, sql).WithAfter([]sdk.SchemaObjectIdentifier{t1Id, rootId})
		t2 := createTaskWithRequest(t, request)

		require.Contains(t, t2.Predecessors, rootId)
		require.Contains(t, t2.Predecessors, t1Id)
		require.Len(t, t2.Predecessors, 2)

		t3Name := random.String()
		t3Id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, t3Name)

		request = sdk.NewCreateTaskRequest(t3Id, sql).WithAfter([]sdk.SchemaObjectIdentifier{t2Id, t1Id})
		t3 := createTaskWithRequest(t, request)

		require.Contains(t, t3.Predecessors, t2Id)
		require.Contains(t, t3.Predecessors, t1Id)
		require.Len(t, t3.Predecessors, 2)

		rootTasks, err := sdk.GetRootTasks(client.Tasks, ctx, rootId)
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t1Id)
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t2Id)
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t3Id)
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		// cannot set ALLOW_OVERLAPPING_EXECUTION on child task
		alterRequest := sdk.NewAlterTaskRequest(t1Id).WithSet(sdk.NewTaskSetRequest().WithAllowOverlappingExecution(sdk.Bool(true)))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.ErrorContains(t, err, "Cannot set allow_overlapping_execution on non-root task")

		// can set ALLOW_OVERLAPPING_EXECUTION on root task
		alterRequest = sdk.NewAlterTaskRequest(rootId).WithSet(sdk.NewTaskSetRequest().WithAllowOverlappingExecution(sdk.Bool(true)))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		// can create cycle, because DAG is suspended
		alterRequest = sdk.NewAlterTaskRequest(t1Id).WithAddAfter([]sdk.SchemaObjectIdentifier{t3Id})
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		// we get an error when trying to start
		alterRequest = sdk.NewAlterTaskRequest(rootId).WithResume(sdk.Bool(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.ErrorContains(t, err, "Graph has at least one cycle containing task")
	})

	t.Run("create dag of tasks - multiple roots", func(t *testing.T) {
		root1Name := random.String()
		root1Id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, root1Name)

		request := sdk.NewCreateTaskRequest(root1Id, sql).WithSchedule(sdk.String("10 MINUTE"))
		root1 := createTaskWithRequest(t, request)

		require.Equal(t, []sdk.SchemaObjectIdentifier{}, root1.Predecessors)

		root2Name := random.String()
		root2Id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, root2Name)

		request = sdk.NewCreateTaskRequest(root2Id, sql).WithSchedule(sdk.String("10 MINUTE"))
		root2 := createTaskWithRequest(t, request)

		require.Equal(t, []sdk.SchemaObjectIdentifier{}, root2.Predecessors)

		t1Name := random.String()
		t1Id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, t1Name)

		request = sdk.NewCreateTaskRequest(t1Id, sql).WithAfter([]sdk.SchemaObjectIdentifier{root1Id, root2Id})
		t1 := createTaskWithRequest(t, request)

		require.Contains(t, t1.Predecessors, root1Id)
		require.Contains(t, t1.Predecessors, root2Id)
		require.Len(t, t1.Predecessors, 2)

		rootTasks, err := sdk.GetRootTasks(client.Tasks, ctx, t1Id)
		require.NoError(t, err)
		require.Len(t, rootTasks, 2)
		require.Contains(t, []sdk.SchemaObjectIdentifier{root1Id, root2Id}, rootTasks[0].ID())
		require.Contains(t, []sdk.SchemaObjectIdentifier{root1Id, root2Id}, rootTasks[1].ID())

		// we get an error when trying to start
		alterRequest := sdk.NewAlterTaskRequest(root1Id).WithResume(sdk.Bool(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.ErrorContains(t, err, "The graph has more than one root task (one without predecessors)")
	})

	// TODO: this fails with `syntax error line 1 at position 89 unexpected 'GRANTS'`.
	// The reason is that in the documentation there is a note: "This parameter is not supported currently.".
	// t.Run("create task: with grants", func(t *testing.T) {
	//	name := randomString(t)
	//	id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
	//
	//	request := NewCreateTaskRequest(id, sql).
	//		WithOrReplace(Bool(true)).
	//		WithCopyGrants(Bool(true))
	//
	//	err := client.Tasks.Create(ctx, request)
	//	require.NoError(t, err)
	//	t.Cleanup(cleanupTaskProvider(id))
	//
	//	task, err := client.Tasks.ShowByID(ctx, id)
	//
	//	require.NoError(t, err)
	//	assertTaskWithOptions(t, task, id, name, "", "", "", "", false, "", nil)
	// })

	t.Run("create task: with tags", func(t *testing.T) {
		tag, tagCleanup := createTag(t, client, testDb(t), testSchema(t))
		t.Cleanup(tagCleanup)

		request := createTaskBasicRequest(t).
			WithTag([]sdk.TagAssociation{{
				Name:  tag.ID(),
				Value: "v1",
			}})

		task := createTaskWithRequest(t, request)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), task.ID(), sdk.ObjectTypeTask)
		require.NoError(t, err)

		assert.Equal(t, "v1", returnedTagValue)
	})

	t.Run("clone task: default", func(t *testing.T) {
		sourceTask := createTask(t)

		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		request := sdk.NewCloneTaskRequest(id, sourceTask.ID())

		err := client.Tasks.Clone(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		task, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTask(t, task, request.GetName())
	})

	t.Run("drop task: existing", func(t *testing.T) {
		request := createTaskBasicRequest(t)
		id := request.GetName()

		err := client.Tasks.Create(ctx, request)
		require.NoError(t, err)

		err = client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(id))
		require.NoError(t, err)

		_, err = client.Tasks.ShowByID(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop task: non-existing", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, "does_not_exist")

		err := client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter task: set value and unset value", func(t *testing.T) {
		task := createTask(t)
		id := task.ID()

		alterRequest := sdk.NewAlterTaskRequest(id).WithSet(sdk.NewTaskSetRequest().WithComment(sdk.String("new comment")))
		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredTask.Comment)

		alterRequest = sdk.NewAlterTaskRequest(id).WithUnset(sdk.NewTaskUnsetRequest().WithComment(sdk.Bool(true)))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredTask.Comment)
	})

	t.Run("alter task: set and unset tag", func(t *testing.T) {
		tag, tagCleanup := createTag(t, client, testDb(t), testSchema(t))
		t.Cleanup(tagCleanup)

		task := createTask(t)
		id := task.ID()

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := sdk.NewAlterTaskRequest(id).WithSetTags(tags)

		err := client.Tasks.Alter(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeTask)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterTaskRequest(id).WithUnsetTags(unsetTags)

		err = client.Tasks.Alter(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeTask)
		require.Error(t, err)
	})

	t.Run("alter task: resume and suspend", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithSchedule(sdk.String("10 MINUTE"))

		task := createTaskWithRequest(t, request)
		id := task.ID()

		assert.Equal(t, sdk.TaskStateSuspended, task.State)

		alterRequest := sdk.NewAlterTaskRequest(id).WithResume(sdk.Bool(true))
		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, sdk.TaskStateStarted, alteredTask.State)

		alterRequest = sdk.NewAlterTaskRequest(id).WithSuspend(sdk.Bool(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, sdk.TaskStateSuspended, alteredTask.State)
	})

	t.Run("alter task: remove after and add after", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithSchedule(sdk.String("10 MINUTE"))

		otherTask := createTaskWithRequest(t, request)
		otherId := otherTask.ID()

		request = createTaskBasicRequest(t).
			WithAfter([]sdk.SchemaObjectIdentifier{otherId})

		task := createTaskWithRequest(t, request)
		id := task.ID()

		assert.Contains(t, task.Predecessors, otherId)

		alterRequest := sdk.NewAlterTaskRequest(id).WithRemoveAfter([]sdk.SchemaObjectIdentifier{otherId})

		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		task, err = client.Tasks.ShowByID(ctx, id)

		require.NoError(t, err)
		assert.Empty(t, task.Predecessors)

		alterRequest = sdk.NewAlterTaskRequest(id).WithAddAfter([]sdk.SchemaObjectIdentifier{otherId})

		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		task, err = client.Tasks.ShowByID(ctx, id)

		require.NoError(t, err)
		assert.Contains(t, task.Predecessors, otherId)
	})

	t.Run("alter task: modify when and as", func(t *testing.T) {
		task := createTask(t)
		id := task.ID()

		newSql := "SELECT CURRENT_DATE"
		alterRequest := sdk.NewAlterTaskRequest(id).WithModifyAs(sdk.String(newSql))
		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, newSql, alteredTask.Definition)

		newWhen := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
		alterRequest = sdk.NewAlterTaskRequest(id).WithModifyWhen(sdk.String(newWhen))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, newWhen, alteredTask.Condition)
	})

	t.Run("show task: default", func(t *testing.T) {
		task1 := createTask(t)
		task2 := createTask(t)

		showRequest := sdk.NewShowTaskRequest()
		returnedTasks, err := client.Tasks.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedTasks))
		assert.Contains(t, returnedTasks, *task1)
		assert.Contains(t, returnedTasks, *task2)
	})

	t.Run("show task: terse", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithSchedule(sdk.String("10 MINUTE"))

		task := createTaskWithRequest(t, request)

		showRequest := sdk.NewShowTaskRequest().WithTerse(sdk.Bool(true))
		returnedTasks, err := client.Tasks.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 1, len(returnedTasks))
		assertTaskTerse(t, &returnedTasks[0], task.ID(), "10 MINUTE")
	})

	t.Run("show task: with options", func(t *testing.T) {
		task1 := createTask(t)
		task2 := createTask(t)

		showRequest := sdk.NewShowTaskRequest().
			WithLike(&sdk.Like{Pattern: &task1.Name}).
			WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(testDb(t).Name, testSchema(t).Name)}).
			WithLimit(&sdk.LimitFrom{Rows: sdk.Int(5)})
		returnedTasks, err := client.Tasks.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedTasks))
		assert.Contains(t, returnedTasks, *task1)
		assert.NotContains(t, returnedTasks, *task2)
	})

	t.Run("describe task: default", func(t *testing.T) {
		task := createTask(t)

		returnedTask, err := client.Tasks.Describe(ctx, task.ID())
		require.NoError(t, err)

		assertTask(t, returnedTask, task.ID())
	})

	t.Run("execute task: default", func(t *testing.T) {
		task := createTask(t)

		executeRequest := sdk.NewExecuteTaskRequest(task.ID())
		err := client.Tasks.Execute(ctx, executeRequest)
		require.NoError(t, err)
	})
}
