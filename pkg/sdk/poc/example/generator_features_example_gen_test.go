package example

import "testing"

func TestFeaturesExample_Alter(t *testing.T) {
	id := RandomDatabaseObjectIdentifier(t)
	// Minimal valid AlterFeaturesExamplesOptions
	defaultOpts := func() *AlterFeaturesExamplesOptions {
		return &AlterFeaturesExamplesOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterFeaturesExamplesOptions = nil
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
