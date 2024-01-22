package sdk

import "testing"

func TestMaterializedViews_Create(t *testing.T) {
	id := RandomSchemaObjectIdentifier()
	sql := "SELECT id FROM t"

	// Minimal valid CreateMaterializedViewOptions
	defaultOpts := func() *CreateMaterializedViewOptions {
		return &CreateMaterializedViewOptions{
			name: id,
			sql:  sql,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateMaterializedViewOptions = nil
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
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateMaterializedViewOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: valid identifier for [opts.RowAccessPolicy.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.RowAccessPolicy = &MaterializedViewRowAccessPolicy{
			RowAccessPolicy: NewSchemaObjectIdentifier("", "", ""),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: [opts.RowAccessPolicy.On] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.RowAccessPolicy = &MaterializedViewRowAccessPolicy{
			RowAccessPolicy: RandomSchemaObjectIdentifier(),
			On:              []string{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateMaterializedViewOptions.RowAccessPolicy", "On"))
	})

	t.Run("validation: [opts.ClusterBy.Expressions] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.ClusterBy = &MaterializedViewClusterBy{
			Expressions: []MaterializedViewClusterByExpression{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateMaterializedViewOptions.ClusterBy", "Expressions"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE MATERIALIZED VIEW %s AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("all options", func(t *testing.T) {
		rowAccessPolicyId := RandomSchemaObjectIdentifier()
		tag1Id := RandomSchemaObjectIdentifier()
		tag2Id := RandomSchemaObjectIdentifier()
		maskingPolicy1Id := RandomSchemaObjectIdentifier()
		maskingPolicy2Id := RandomSchemaObjectIdentifier()

		req := NewCreateMaterializedViewRequest(id, sql).
			WithOrReplace(Bool(true)).
			WithSecure(Bool(true)).
			WithColumns([]MaterializedViewColumnRequest{
				*NewMaterializedViewColumnRequest("column_without_comment"),
				*NewMaterializedViewColumnRequest("column_with_comment").WithComment(String("column 2 comment")),
			}).
			WithColumnsMaskingPolicies([]MaterializedViewColumnMaskingPolicyRequest{
				*NewMaterializedViewColumnMaskingPolicyRequest("column", maskingPolicy1Id).
					WithUsing([]string{"a", "b"}).
					WithTag([]TagAssociation{{
						Name:  tag1Id,
						Value: "v1",
					}}),
				*NewMaterializedViewColumnMaskingPolicyRequest("column 2", maskingPolicy2Id),
			}).
			WithCopyGrants(Bool(true)).
			WithComment(String("comment")).
			WithRowAccessPolicy(NewMaterializedViewRowAccessPolicyRequest(rowAccessPolicyId, []string{"c", "d"})).
			WithTag([]TagAssociation{{
				Name:  tag2Id,
				Value: "v2",
			}}).
			WithClusterBy(NewMaterializedViewClusterByRequest().WithExpressions([]MaterializedViewClusterByExpressionRequest{{"column_without_comment"}, {"column_with_comment"}}))

		assertOptsValidAndSQLEquals(t, req.toOpts(), `CREATE OR REPLACE SECURE MATERIALIZED VIEW %s COPY GRANTS ("column_without_comment", "column_with_comment" COMMENT 'column 2 comment') column MASKING POLICY %s USING (a, b) TAG (%s = 'v1'), column 2 MASKING POLICY %s COMMENT = 'comment' ROW ACCESS POLICY %s ON (c, d) TAG (%s = 'v2') CLUSTER BY ("column_without_comment", "column_with_comment") AS %s`, id.FullyQualifiedName(), maskingPolicy1Id.FullyQualifiedName(), tag1Id.FullyQualifiedName(), maskingPolicy2Id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(), tag2Id.FullyQualifiedName(), sql)
	})
}

func TestMaterializedViews_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterMaterializedViewOptions
	defaultOpts := func() *AlterMaterializedViewOptions {
		return &AlterMaterializedViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterMaterializedViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.ClusterBy opts.DropClusteringKey opts.SuspendRecluster opts.ResumeRecluster opts.Suspend opts.Resume opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterMaterializedViewOptions", "RenameTo", "ClusterBy", "DropClusteringKey", "SuspendRecluster", "ResumeRecluster", "Suspend", "Resume", "Set", "Unset"))
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.ClusterBy opts.DropClusteringKey opts.SuspendRecluster opts.ResumeRecluster opts.Suspend opts.Resume opts.Set opts.Unset] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SuspendRecluster = Bool(true)
		opts.Suspend = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterMaterializedViewOptions", "RenameTo", "ClusterBy", "DropClusteringKey", "SuspendRecluster", "ResumeRecluster", "Suspend", "Resume", "Set", "Unset"))
	})

	t.Run("validation: at least one of the fields [opts.Set.Secure opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &MaterializedViewSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterMaterializedViewOptions.Set", "Secure", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.Secure opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &MaterializedViewUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterMaterializedViewOptions.Unset", "Secure", "Comment"))
	})

	t.Run("rename", func(t *testing.T) {
		newId := RandomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.RenameTo = &newId
		assertOptsValidAndSQLEquals(t, opts, "ALTER MATERIALIZED VIEW %s RENAME TO %s", id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("cluster by", func(t *testing.T) {
		opts := defaultOpts()
		opts.ClusterBy = &MaterializedViewClusterBy{
			Expressions: []MaterializedViewClusterByExpression{{"column"}},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER MATERIALIZED VIEW %s CLUSTER BY ("column")`, id.FullyQualifiedName())
	})
}

func TestMaterializedViews_Drop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DropMaterializedViewOptions
	defaultOpts := func() *DropMaterializedViewOptions {
		return &DropMaterializedViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropMaterializedViewOptions = nil
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

func TestMaterializedViews_Show(t *testing.T) {
	// Minimal valid ShowMaterializedViewOptions
	defaultOpts := func() *ShowMaterializedViewOptions {
		return &ShowMaterializedViewOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowMaterializedViewOptions = nil
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

func TestMaterializedViews_Describe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DescribeMaterializedViewOptions
	defaultOpts := func() *DescribeMaterializedViewOptions {
		return &DescribeMaterializedViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeMaterializedViewOptions = nil
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
