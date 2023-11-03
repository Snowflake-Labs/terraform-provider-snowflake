package sdk

import "testing"

func TestViews_Create(t *testing.T) {
	id := RandomSchemaObjectIdentifier()
	sql := "SELECT id FROM t"

	// Minimal valid CreateViewOptions
	defaultOpts := func() *CreateViewOptions {
		return &CreateViewOptions{
			name: id,
			sql:  sql,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateViewOptions = nil
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
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateViewOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: valid identifier for [opts.RowAccessPolicy.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.RowAccessPolicy = &ViewRowAccessPolicy{
			RowAccessPolicy: NewSchemaObjectIdentifier("", "", ""),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE VIEW %s AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("all options", func(t *testing.T) {
		rowAccessPolicyId := RandomSchemaObjectIdentifier()
		tag1Id := RandomSchemaObjectIdentifier()
		tag2Id := RandomSchemaObjectIdentifier()
		maskingPolicy1Id := RandomSchemaObjectIdentifier()
		maskingPolicy2Id := RandomSchemaObjectIdentifier()

		req := NewCreateViewRequest(id, sql).
			WithOrReplace(Bool(true)).
			WithSecure(Bool(true)).
			WithTemporary(Bool(true)).
			WithRecursive(Bool(true)).
			WithColumns([]ViewColumnRequest{
				*NewViewColumnRequest("column_without_comment"),
				*NewViewColumnRequest("column_with_comment").WithComment(String("column 2 comment")),
			}).
			WithColumnsMaskingPolicies([]ViewColumnMaskingPolicyRequest{
				*NewViewColumnMaskingPolicyRequest("column", maskingPolicy1Id).
					WithUsing([]string{"a", "b"}).
					WithTag([]TagAssociation{{
						Name:  tag1Id,
						Value: "v1",
					}}),
				*NewViewColumnMaskingPolicyRequest("column 2", maskingPolicy2Id),
			}).
			WithCopyGrants(Bool(true)).
			WithComment(String("comment")).
			WithRowAccessPolicy(NewViewRowAccessPolicyRequest(rowAccessPolicyId).WithOn([]string{"c", "d"})).
			WithTag([]TagAssociation{{
				Name:  tag2Id,
				Value: "v2",
			}})

		assertOptsValidAndSQLEquals(t, req.toOpts(), `CREATE OR REPLACE SECURE TEMPORARY RECURSIVE VIEW %s ("column_without_comment", "column_with_comment" COMMENT 'column 2 comment') column MASKING POLICY %s USING (a, b) TAG (%s = 'v1'), column 2 MASKING POLICY %s COPY GRANTS COMMENT = 'comment' ROW ACCESS POLICY %s ON (c, d) TAG (%s = 'v2') AS %s`, id.FullyQualifiedName(), maskingPolicy1Id.FullyQualifiedName(), tag1Id.FullyQualifiedName(), maskingPolicy2Id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(), tag2Id.FullyQualifiedName(), sql)
	})
}

func TestViews_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterViewOptions
	defaultOpts := func() *AlterViewOptions {
		return &AlterViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetComment opts.UnsetComment opts.SetSecure opts.SetChangeTracking opts.UnsetSecure opts.SetTags opts.UnsetTags opts.AddRowAccessPolicy opts.DropRowAccessPolicy opts.DropAndAddRowAccessPolicy opts.DropAllRowAccessPolicies opts.SetMaskingPolicyOnColumn opts.UnsetMaskingPolicyOnColumn opts.SetTagsOnColumn opts.UnsetTagsOnColumn] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterViewOptions", "RenameTo", "SetComment", "UnsetComment", "SetSecure", "SetChangeTracking", "UnsetSecure", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "SetMaskingPolicyOnColumn", "UnsetMaskingPolicyOnColumn", "SetTagsOnColumn", "UnsetTagsOnColumn"))
	})

	t.Run("validation: valid identifier for [opts.AddRowAccessPolicy.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.DropRowAccessPolicy.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.DropAndAddRowAccessPolicy.Drop.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.DropAndAddRowAccessPolicy.Add.RowAccessPolicy]", func(t *testing.T) {
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

func TestViews_Drop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DropViewOptions
	defaultOpts := func() *DropViewOptions {
		return &DropViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropViewOptions = nil
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

func TestViews_Show(t *testing.T) {
	// Minimal valid ShowViewOptions
	defaultOpts := func() *ShowViewOptions {
		return &ShowViewOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowViewOptions = nil
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

func TestViews_Describe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DescribeViewOptions
	defaultOpts := func() *DescribeViewOptions {
		return &DescribeViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeViewOptions = nil
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
