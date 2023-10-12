package sdk

import "testing"

func TestSessionPolicies_Create(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid CreateSessionPolicyOptions
	defaultOpts := func() *CreateSessionPolicyOptions {
		return &CreateSessionPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSessionPolicyOptions = nil
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
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSessionPolicyOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE SESSION POLICY %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.SessionIdleTimeoutMins = Int(5)
		opts.SessionUiIdleTimeoutMins = Int(34)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SESSION POLICY %s SESSION_IDLE_TIMEOUT_MINS = 5 SESSION_UI_IDLE_TIMEOUT_MINS = 34 COMMENT = 'some comment'", id.FullyQualifiedName())
	})
}

func TestSessionPolicies_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterSessionPolicyOptions
	defaultOpts := func() *AlterSessionPolicyOptions {
		return &AlterSessionPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterSessionPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.Set opts.SetTags opts.UnsetTags opts.Unset] should be present - none present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RenameTo", "Set", "SetTags", "UnsetTags", "Unset"))
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.Set opts.SetTags opts.UnsetTags opts.Unset] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SessionPolicySet{
			Comment: String("some comment"),
		}
		opts.Unset = &SessionPolicyUnset{
			Comment: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RenameTo", "Set", "SetTags", "UnsetTags", "Unset"))
	})

	t.Run("validation: at least one of the fields [opts.Set.SessionIdleTimeoutMins opts.Set.SessionUiIdleTimeoutMins opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SessionPolicySet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("SessionIdleTimeoutMins", "SessionUiIdleTimeoutMins", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.SessionIdleTimeoutMins opts.Unset.SessionUiIdleTimeoutMins opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &SessionPolicyUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("SessionIdleTimeoutMins", "SessionUiIdleTimeoutMins", "Comment"))
	})

	t.Run("alter set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SessionPolicySet{
			Comment: String("some comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SESSION POLICY %s SET COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("alter unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &SessionPolicyUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SESSION POLICY %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("alter rename", func(t *testing.T) {
		opts := defaultOpts()
		newId := RandomSchemaObjectIdentifier()
		opts.RenameTo = &newId
		assertOptsValidAndSQLEquals(t, opts, "ALTER SESSION POLICY %s RENAME TO %s", id.FullyQualifiedName(), newId.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER SESSION POLICY %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName())
	})

	t.Run("alter unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SESSION POLICY %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})
}

func TestSessionPolicies_Drop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DropSessionPolicyOptions
	defaultOpts := func() *DropSessionPolicyOptions {
		return &DropSessionPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropSessionPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP SESSION POLICY %s", id.FullyQualifiedName())
	})
}

func TestSessionPolicies_Show(t *testing.T) {
	// Minimal valid ShowSessionPolicyOptions
	defaultOpts := func() *ShowSessionPolicyOptions {
		return &ShowSessionPolicyOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowSessionPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW SESSION POLICIES")
	})
}

func TestSessionPolicies_Describe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DescribeSessionPolicyOptions
	defaultOpts := func() *DescribeSessionPolicyOptions {
		return &DescribeSessionPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeSessionPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE SESSION POLICY %s", id.FullyQualifiedName())
	})
}
