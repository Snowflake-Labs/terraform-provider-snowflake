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
