package example

import "testing"

func TestToOptsOptionalExamples_Create(t *testing.T) {
	id := RandomDatabaseObjectIdentifier(t)
	// Minimal valid CreateToOptsOptionalExampleOptions
	defaultOpts := func() *CreateToOptsOptionalExampleOptions {
		return &CreateToOptsOptionalExampleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateToOptsOptionalExampleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
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

func TestToOptsOptionalExamples_Alter(t *testing.T) {
	id := RandomDatabaseObjectIdentifier(t)
	// Minimal valid AlterToOptsOptionalExampleOptions
	defaultOpts := func() *AlterToOptsOptionalExampleOptions {
		return &AlterToOptsOptionalExampleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterToOptsOptionalExampleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
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
