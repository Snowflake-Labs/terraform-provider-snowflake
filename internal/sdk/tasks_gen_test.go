// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import (
	"fmt"
	"testing"
)

func TestTasks_Create(t *testing.T) {
	id := RandomSchemaObjectIdentifier()
	sql := "SELECT CURRENT_TIMESTAMP"

	// Minimal valid CreateTaskOptions
	defaultOpts := func() *CreateTaskOptions {
		return &CreateTaskOptions{
			name: id,
			sql:  sql,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateTaskOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: exactly one field from [opts.Warehouse.Warehouse opts.Warehouse.UserTaskManagedInitialWarehouseSize] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Warehouse = &CreateTaskWarehouse{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("Warehouse", "UserTaskManagedInitialWarehouseSize"))
	})

	t.Run("validation: opts.SessionParameters.SessionParameters should be valid", func(t *testing.T) {
		opts := defaultOpts()
		opts.SessionParameters = &SessionParameters{
			JSONIndent: Int(25),
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("JSON_INDENT must be between 0 and 16"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE TASK %s AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("with initial warehouse size", func(t *testing.T) {
		req := NewCreateTaskRequest(id, sql).
			WithWarehouse(NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(&WarehouseSizeXSmall))
		assertOptsValidAndSQLEquals(t, req.toOpts(), "CREATE TASK %s USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = 'XSMALL' AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("all options", func(t *testing.T) {
		warehouseId := RandomAccountObjectIdentifier()
		otherTaskId := RandomSchemaObjectIdentifier()
		tagId := RandomSchemaObjectIdentifier()

		req := NewCreateTaskRequest(id, sql).
			WithOrReplace(Bool(true)).
			WithWarehouse(NewCreateTaskWarehouseRequest().WithWarehouse(&warehouseId)).
			WithSchedule(String("10 MINUTE")).
			WithConfig(String(`$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`)).
			WithAllowOverlappingExecution(Bool(true)).
			WithSessionParameters(&SessionParameters{
				JSONIndent: Int(10),
			}).
			WithUserTaskTimeoutMs(Int(5)).
			WithSuspendTaskAfterNumFailures(Int(6)).
			WithErrorIntegration(String("some_error_integration")).
			WithCopyGrants(Bool(true)).
			WithComment(String("some comment")).
			WithAfter([]SchemaObjectIdentifier{otherTaskId}).
			WithTag([]TagAssociation{{
				Name:  tagId,
				Value: "v1",
			}}).
			WithWhen(String(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`))

		assertOptsValidAndSQLEquals(t, req.toOpts(), "CREATE OR REPLACE TASK %s WAREHOUSE = %s SCHEDULE = '10 MINUTE' CONFIG = $${\"output_dir\": \"/temp/test_directory/\", \"learning_rate\": 0.1}$$ ALLOW_OVERLAPPING_EXECUTION = true JSON_INDENT = 10 USER_TASK_TIMEOUT_MS = 5 SUSPEND_TASK_AFTER_NUM_FAILURES = 6 ERROR_INTEGRATION = some_error_integration COPY GRANTS COMMENT = 'some comment' AFTER %s TAG (%s = 'v1') WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS SELECT CURRENT_TIMESTAMP", id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), otherTaskId.FullyQualifiedName(), tagId.FullyQualifiedName())
	})
}

func TestTasks_Clone(t *testing.T) {
	id := RandomSchemaObjectIdentifier()
	sourceId := RandomSchemaObjectIdentifier()

	// Minimal valid CloneTaskOptions
	defaultOpts := func() *CloneTaskOptions {
		return &CloneTaskOptions{
			name:       id,
			sourceTask: sourceId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CloneTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.sourceTask]", func(t *testing.T) {
		opts := defaultOpts()
		opts.sourceTask = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE TASK %s CLONE %s", id.FullyQualifiedName(), sourceId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.CopyGrants = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE TASK %s CLONE %s COPY GRANTS", id.FullyQualifiedName(), sourceId.FullyQualifiedName())
	})
}

func TestTasks_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()
	otherTaskId := RandomSchemaObjectIdentifier()

	// Minimal valid AlterTaskOptions
	defaultOpts := func() *AlterTaskOptions {
		return &AlterTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Resume opts.Suspend opts.RemoveAfter opts.AddAfter opts.Set opts.Unset opts.SetTags opts.UnsetTags opts.ModifyAs opts.ModifyWhen] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("Resume", "Suspend", "RemoveAfter", "AddAfter", "Set", "Unset", "SetTags", "UnsetTags", "ModifyAs", "ModifyWhen"))
	})

	t.Run("validation: exactly one field from [opts.Resume opts.Suspend opts.RemoveAfter opts.AddAfter opts.Set opts.Unset opts.SetTags opts.UnsetTags opts.ModifyAs opts.ModifyWhen] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = Bool(true)
		opts.Suspend = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("Resume", "Suspend", "RemoveAfter", "AddAfter", "Set", "Unset", "SetTags", "UnsetTags", "ModifyAs", "ModifyWhen"))
	})

	t.Run("validation: at least one of the fields [opts.Set.Warehouse opts.Set.UserTaskManagedInitialWarehouseSize opts.Set.Schedule opts.Set.Config opts.Set.AllowOverlappingExecution opts.Set.UserTaskTimeoutMs opts.Set.SuspendTaskAfterNumFailures opts.Set.ErrorIntegration opts.Set.Comment opts.Set.SessionParameters] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("Warehouse", "UserTaskManagedInitialWarehouseSize", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "ErrorIntegration", "Comment", "SessionParameters"))
	})

	t.Run("validation: conflicting fields for [opts.Set.Warehouse opts.Set.UserTaskManagedInitialWarehouseSize]", func(t *testing.T) {
		warehouseId := RandomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &TaskSet{}
		opts.Set.Warehouse = &warehouseId
		opts.Set.UserTaskManagedInitialWarehouseSize = &WarehouseSizeXSmall
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("Set", "Warehouse", "UserTaskManagedInitialWarehouseSize"))
	})

	t.Run("validation: opts.Set.SessionParameters.SessionParameters should be valid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{}
		opts.Set.SessionParameters = &SessionParameters{
			JSONIndent: Int(25),
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("JSON_INDENT must be between 0 and 16"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.Warehouse opts.Unset.Schedule opts.Unset.Config opts.Unset.AllowOverlappingExecution opts.Unset.UserTaskTimeoutMs opts.Unset.SuspendTaskAfterNumFailures opts.Unset.ErrorIntegration opts.Unset.Comment opts.Unset.SessionParametersUnset] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("Warehouse", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "ErrorIntegration", "Comment", "SessionParametersUnset"))
	})

	t.Run("validation: opts.Unset.SessionParametersUnset.SessionParametersUnset should be valid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{}
		opts.Unset.SessionParametersUnset = &SessionParametersUnset{}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("at least one session parameter must be set"))
	})

	t.Run("alter resume", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s RESUME", id.FullyQualifiedName())
	})

	t.Run("alter suspend", func(t *testing.T) {
		opts := defaultOpts()
		opts.Suspend = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SUSPEND", id.FullyQualifiedName())
	})

	t.Run("alter remove after", func(t *testing.T) {
		opts := defaultOpts()
		opts.RemoveAfter = []SchemaObjectIdentifier{otherTaskId}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s REMOVE AFTER %s", id.FullyQualifiedName(), otherTaskId.FullyQualifiedName())
	})

	t.Run("alter add after", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddAfter = []SchemaObjectIdentifier{otherTaskId}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s ADD AFTER %s", id.FullyQualifiedName(), otherTaskId.FullyQualifiedName())
	})

	t.Run("alter set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{
			Comment: String("some comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("alter set warehouse", func(t *testing.T) {
		warehouseId := RandomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &TaskSet{
			Warehouse: &warehouseId,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET WAREHOUSE = %s", id.FullyQualifiedName(), warehouseId.FullyQualifiedName())
	})

	t.Run("alter set session parameter", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{
			SessionParameters: &SessionParameters{
				JSONIndent: Int(15),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET JSON_INDENT = 15", id.FullyQualifiedName())
	})

	t.Run("alter unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("alter set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
			{
				Name:  NewAccountObjectIdentifier("tag2"),
				Value: "value2",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TASK %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName())
	})

	t.Run("alter unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TASK %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})

	t.Run("alter modify as", func(t *testing.T) {
		opts := defaultOpts()
		opts.ModifyAs = String("new as")
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s MODIFY AS new as", id.FullyQualifiedName())
	})

	t.Run("alter modify when", func(t *testing.T) {
		opts := defaultOpts()
		opts.ModifyWhen = String("new when")
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s MODIFY WHEN new when", id.FullyQualifiedName())
	})
}

func TestTasks_Drop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DropTaskOptions
	defaultOpts := func() *DropTaskOptions {
		return &DropTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP TASK %s", id.FullyQualifiedName())
	})
}

func TestTasks_Show(t *testing.T) {
	// Minimal valid ShowTaskOptions
	defaultOpts := func() *ShowTaskOptions {
		return &ShowTaskOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW TASKS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		opts.Like = &Like{
			Pattern: String("myaccount"),
		}
		opts.In = &In{
			Account: Bool(true),
		}
		opts.StartsWith = String("abc")
		opts.RootOnly = Bool(true)
		opts.Limit = &LimitFrom{Rows: Int(10)}
		assertOptsValidAndSQLEquals(t, opts, "SHOW TERSE TASKS LIKE 'myaccount' IN ACCOUNT STARTS WITH 'abc' ROOT ONLY LIMIT 10")
	})
}

func TestTasks_Describe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DescribeTaskOptions
	defaultOpts := func() *DescribeTaskOptions {
		return &DescribeTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE TASK %s", id.FullyQualifiedName())
	})
}

func TestTasks_Execute(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid ExecuteTaskOptions
	defaultOpts := func() *ExecuteTaskOptions {
		return &ExecuteTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ExecuteTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "EXECUTE TASK %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.RetryLast = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "EXECUTE TASK %s RETRY LAST", id.FullyQualifiedName())
	})
}
