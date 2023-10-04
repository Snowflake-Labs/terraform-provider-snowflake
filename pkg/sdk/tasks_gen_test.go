package sdk

import (
	"fmt"
	"testing"
)

func TestTasks_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)
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
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
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

	t.Run("all options", func(t *testing.T) {
		warehouseId := randomAccountObjectIdentifier(t)
		otherTaskId := randomSchemaObjectIdentifier(t)
		tagId := randomSchemaObjectIdentifier(t)

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

		assertOptsValidAndSQLEquals(t, req.toOpts(), "CREATE OR REPLACE TASK %s WAREHOUSE %s SCHEDULE = '10 MINUTE' CONFIG = $${\"output_dir\": \"/temp/test_directory/\", \"learning_rate\": 0.1}$$ ALLOW_OVERLAPPING_EXECUTION = true JSON_INDENT = 10 USER_TASK_TIMEOUT_MS = 5 SUSPEND_TASK_AFTER_NUM_FAILURES = 6 ERROR_INTEGRATION = some_error_integration COPY GRANTS COMMENT = 'some comment' AFTER = %s TAG (%s = 'v1') WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS SELECT CURRENT_TIMESTAMP", id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), otherTaskId.FullyQualifiedName(), tagId.FullyQualifiedName())
	})
}

func TestTasks_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	// Minimal valid AlterTaskOptions
	defaultOpts := func() *AlterTaskOptions {
		return &AlterTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Resume opts.Suspend opts.RemoveAfter opts.AddAfter opts.Set opts.Unset opts.SetTags opts.UnsetTags opts.ModifyAs opts.ModifyWhen] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("Resume", "Suspend", "RemoveAfter", "AddAfter", "Set", "Unset", "SetTags", "UnsetTags", "ModifyAs", "ModifyWhen"))
	})

	t.Run("validation: at least one of the fields [opts.Set.Warehouse opts.Set.Schedule opts.Set.Config opts.Set.AllowOverlappingExecution opts.Set.UserTaskTimeoutMs opts.Set.SuspendTaskAfterNumFailures opts.Set.Comment opts.Set.SessionParameters] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("Warehouse", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "Comment", "SessionParameters"))
	})

	t.Run("validation: valid identifier for [opts.Set.Warehouse] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: opts.Set.SessionParameters.SessionParameters should be valid", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf(""))
	})

	t.Run("validation: at least one of the fields [opts.Unset.Warehouse opts.Unset.Schedule opts.Unset.Config opts.Unset.AllowOverlappingExecution opts.Unset.UserTaskTimeoutMs opts.Unset.SuspendTaskAfterNumFailures opts.Unset.Comment opts.Unset.SessionParametersUnset] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("Warehouse", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "Comment", "SessionParametersUnset"))
	})

	t.Run("validation: opts.Unset.SessionParametersUnset.SessionParametersUnset should be valid", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf(""))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestTasks_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	// Minimal valid DropTaskOptions
	defaultOpts := func() *DropTaskOptions {
		return &DropTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestTasks_Show(t *testing.T) {
	//id := randomSchemaObjectIdentifier(t)

	// Minimal valid ShowTaskOptions
	defaultOpts := func() *ShowTaskOptions {
		return &ShowTaskOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestTasks_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	// Minimal valid DescribeTaskOptions
	defaultOpts := func() *DescribeTaskOptions {
		return &DescribeTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestTasks_Execute(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	// Minimal valid ExecuteTaskOptions
	defaultOpts := func() *ExecuteTaskOptions {
		return &ExecuteTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ExecuteTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}
