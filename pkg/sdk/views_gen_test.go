package sdk

import (
	"testing"
)

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
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetComment opts.UnsetComment opts.SetSecure opts.SetChangeTracking opts.UnsetSecure opts.SetTags opts.UnsetTags opts.AddRowAccessPolicy opts.DropRowAccessPolicy opts.DropAndAddRowAccessPolicy opts.DropAllRowAccessPolicies opts.SetMaskingPolicyOnColumn opts.UnsetMaskingPolicyOnColumn opts.SetTagsOnColumn opts.UnsetTagsOnColumn] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterViewOptions", "RenameTo", "SetComment", "UnsetComment", "SetSecure", "SetChangeTracking", "UnsetSecure", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "SetMaskingPolicyOnColumn", "UnsetMaskingPolicyOnColumn", "SetTagsOnColumn", "UnsetTagsOnColumn"))
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetComment opts.UnsetComment opts.SetSecure opts.SetChangeTracking opts.UnsetSecure opts.SetTags opts.UnsetTags opts.AddRowAccessPolicy opts.DropRowAccessPolicy opts.DropAndAddRowAccessPolicy opts.DropAllRowAccessPolicies opts.SetMaskingPolicyOnColumn opts.UnsetMaskingPolicyOnColumn opts.SetTagsOnColumn opts.UnsetTagsOnColumn] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetChangeTracking = Bool(true)
		opts.DropAllRowAccessPolicies = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterViewOptions", "RenameTo", "SetComment", "UnsetComment", "SetSecure", "SetChangeTracking", "UnsetSecure", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "SetMaskingPolicyOnColumn", "UnsetMaskingPolicyOnColumn", "SetTagsOnColumn", "UnsetTagsOnColumn"))
	})

	t.Run("validation: valid identifier for [opts.DropRowAccessPolicy.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropRowAccessPolicy = &ViewDropRowAccessPolicy{
			RowAccessPolicy: NewSchemaObjectIdentifier("", "", ""),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.DropAndAddRowAccessPolicy.Drop.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropAndAddRowAccessPolicy = &ViewDropAndAddRowAccessPolicy{
			Drop: ViewDropRowAccessPolicy{
				RowAccessPolicy: NewSchemaObjectIdentifier("", "", ""),
			},
			Add: ViewAddRowAccessPolicy{
				RowAccessPolicy: RandomSchemaObjectIdentifier(),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.DropAndAddRowAccessPolicy.Add.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropAndAddRowAccessPolicy = &ViewDropAndAddRowAccessPolicy{
			Drop: ViewDropRowAccessPolicy{
				RowAccessPolicy: RandomSchemaObjectIdentifier(),
			},
			Add: ViewAddRowAccessPolicy{
				RowAccessPolicy: NewSchemaObjectIdentifier("", "", ""),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("rename", func(t *testing.T) {
		newId := RandomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.RenameTo = &newId
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s RENAME TO %s", id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s SET COMMENT = 'comment'", id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetComment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("set secure", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetSecure = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s SET SECURE", id.FullyQualifiedName())
	})

	t.Run("set change tracking: true", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetChangeTracking = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s SET CHANGE TRACKING = true", id.FullyQualifiedName())
	})

	t.Run("set change tracking: false", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetChangeTracking = Bool(false)
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s SET CHANGE TRACKING = false", id.FullyQualifiedName())
	})

	t.Run("unset secure", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetSecure = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s UNSET SECURE", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER VIEW %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER VIEW %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})

	t.Run("add row access policy", func(t *testing.T) {
		rowAccessPolicyId := RandomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.AddRowAccessPolicy = &ViewAddRowAccessPolicy{
			RowAccessPolicy: rowAccessPolicyId,
			On:              []string{"a", "b"},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s ADD ROW ACCESS POLICY %s ON (a, b)", id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("drop row access policy", func(t *testing.T) {
		rowAccessPolicyId := RandomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.DropRowAccessPolicy = &ViewDropRowAccessPolicy{
			RowAccessPolicy: rowAccessPolicyId,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s DROP ROW ACCESS POLICY %s", id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("drop and add row access policy", func(t *testing.T) {
		rowAccessPolicy1Id := RandomSchemaObjectIdentifier()
		rowAccessPolicy2Id := RandomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.DropAndAddRowAccessPolicy = &ViewDropAndAddRowAccessPolicy{
			Drop: ViewDropRowAccessPolicy{
				RowAccessPolicy: rowAccessPolicy1Id,
			},
			Add: ViewAddRowAccessPolicy{
				RowAccessPolicy: rowAccessPolicy2Id,
				On:              []string{"a", "b"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s DROP ROW ACCESS POLICY %s, ADD ROW ACCESS POLICY %s ON (a, b)", id.FullyQualifiedName(), rowAccessPolicy1Id.FullyQualifiedName(), rowAccessPolicy2Id.FullyQualifiedName())
	})

	t.Run("drop all row access policies", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropAllRowAccessPolicies = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s DROP ALL ROW ACCESS POLICIES", id.FullyQualifiedName())
	})

	t.Run("set masking policy on column", func(t *testing.T) {
		maskingPolicyId := RandomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.SetMaskingPolicyOnColumn = &ViewSetColumnMaskingPolicy{
			Name:          "column",
			MaskingPolicy: maskingPolicyId,
			Using:         []string{"a", "b"},
			Force:         Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s ALTER COLUMN column SET MASKING POLICY %s USING (a, b) FORCE", id.FullyQualifiedName(), maskingPolicyId.FullyQualifiedName())
	})

	t.Run("unset masking policy on column", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetMaskingPolicyOnColumn = &ViewUnsetColumnMaskingPolicy{
			Name: "column",
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s ALTER COLUMN column UNSET MASKING POLICY", id.FullyQualifiedName())
	})

	t.Run("set tags on column", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTagsOnColumn = &ViewSetColumnTags{
			Name: "column",
			SetTags: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag1"),
					Value: "value1",
				},
				{
					Name:  NewAccountObjectIdentifier("tag2"),
					Value: "value2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER VIEW %s ALTER COLUMN column SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName())
	})

	t.Run("unset tags on column", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTagsOnColumn = &ViewUnsetColumnTags{
			Name: "column",
			UnsetTags: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag1"),
				NewAccountObjectIdentifier("tag2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER VIEW %s ALTER COLUMN column UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
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
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP VIEW %s", id.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, "SHOW VIEWS")
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
		opts.Limit = &LimitFrom{Rows: Int(10)}
		assertOptsValidAndSQLEquals(t, opts, "SHOW TERSE VIEWS LIKE 'myaccount' IN ACCOUNT STARTS WITH 'abc' LIMIT 10")
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
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE VIEW %s", id.FullyQualifiedName())
	})
}
