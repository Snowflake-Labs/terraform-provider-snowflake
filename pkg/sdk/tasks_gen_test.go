package sdk

import "testing"

func TestTasks_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	// Minimal valid CreateTaskOptions
	defaultOpts := func() *CreateTaskOptions {
		return &CreateTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateTaskOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: valid identifier for [opts.Warehouse.Warehouse]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Warehouse.Warehouse opts.Warehouse.UserTaskManagedInitialWarehouseSize] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("Warehouse", "UserTaskManagedInitialWarehouseSize"))
	})

	t.Run("validation: opts.SessionParameters.SessionParameters should be valid", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, err)
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
		assertOptsInvalidJoinedErrors(t, opts, err)
	})

	t.Run("validation: at least one of the fields [opts.Unset.Warehouse opts.Unset.Schedule opts.Unset.Config opts.Unset.AllowOverlappingExecution opts.Unset.UserTaskTimeoutMs opts.Unset.SuspendTaskAfterNumFailures opts.Unset.Comment opts.Unset.SessionParametersUnset] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("Warehouse", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "Comment", "SessionParametersUnset"))
	})

	t.Run("validation: opts.Unset.SessionParametersUnset.SessionParametersUnset should be valid", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, err)
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
