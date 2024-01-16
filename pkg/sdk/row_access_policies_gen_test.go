package sdk

import "testing"

func TestRowAccessPolicies_Create(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid CreateRowAccessPolicyOptions
	defaultOpts := func() *CreateRowAccessPolicyOptions {
		return &CreateRowAccessPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateRowAccessPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateRowAccessPolicyOptions", "OrReplace", "IfNotExists"))
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

func TestRowAccessPolicies_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterRowAccessPolicyOptions
	defaultOpts := func() *AlterRowAccessPolicyOptions {
		return &AlterRowAccessPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterRowAccessPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetBody opts.SetTags opts.UnsetTags opts.SetComment opts.UnsetComment] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterRowAccessPolicyOptions", "RenameTo", "SetBody", "SetTags", "UnsetTags", "SetComment", "UnsetComment"))
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

func TestRowAccessPolicies_Drop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DropRowAccessPolicyOptions
	defaultOpts := func() *DropRowAccessPolicyOptions {
		return &DropRowAccessPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropRowAccessPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
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

func TestRowAccessPolicies_Show(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid ShowRowAccessPolicyOptions
	defaultOpts := func() *ShowRowAccessPolicyOptions {
		return &ShowRowAccessPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowRowAccessPolicyOptions = nil
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

func TestRowAccessPolicies_Describe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DescribeRowAccessPolicyOptions
	defaultOpts := func() *DescribeRowAccessPolicyOptions {
		return &DescribeRowAccessPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeRowAccessPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
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
