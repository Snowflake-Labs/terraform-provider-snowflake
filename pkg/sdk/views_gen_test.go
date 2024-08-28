package sdk

import (
	"fmt"
	"testing"
)

func TestViews_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
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
		opts.name = emptySchemaObjectIdentifier
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
			RowAccessPolicy: emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: empty columns for row access policy", func(t *testing.T) {
		opts := defaultOpts()
		opts.RowAccessPolicy = &ViewRowAccessPolicy{
			RowAccessPolicy: randomSchemaObjectIdentifier(),
			On:              []Column{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateViewOptions.RowAccessPolicy", "On"))
	})

	t.Run("validation: valid identifier for [opts.MaskingPolicy.MaskingPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Columns = []ViewColumn{
			{
				Name: "foo",
				MaskingPolicy: &ViewColumnMaskingPolicy{
					MaskingPolicy: emptySchemaObjectIdentifier,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.ProjectionPolicy.ProjectionPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Columns = []ViewColumn{
			{
				Name: "foo",
				ProjectionPolicy: &ViewColumnProjectionPolicy{
					ProjectionPolicy: emptySchemaObjectIdentifier,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE VIEW %s AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("all options", func(t *testing.T) {
		rowAccessPolicyId := randomSchemaObjectIdentifier()
		aggregationPolicyId := randomSchemaObjectIdentifier()
		tag1Id := randomSchemaObjectIdentifier()
		tag2Id := randomSchemaObjectIdentifier()
		maskingPolicy1Id := randomSchemaObjectIdentifier()
		maskingPolicy2Id := randomSchemaObjectIdentifier()

		req := NewCreateViewRequest(id, sql).
			WithOrReplace(true).
			WithSecure(true).
			WithTemporary(true).
			WithRecursive(true).
			WithColumns([]ViewColumnRequest{
				*NewViewColumnRequest("column_without_comment"),
				*NewViewColumnRequest("column_with_comment").WithComment("column 2 comment"),
				*NewViewColumnRequest("column").WithMaskingPolicy(
					*NewViewColumnMaskingPolicyRequest(maskingPolicy1Id).
						WithUsing([]Column{{"a"}, {"b"}}),
				).WithTag([]TagAssociation{{Name: tag1Id, Value: "v1"}}),
				*NewViewColumnRequest("column 2").WithProjectionPolicy(
					*NewViewColumnProjectionPolicyRequest(maskingPolicy2Id),
				),
			}).
			WithCopyGrants(true).
			WithComment("comment").
			WithRowAccessPolicy(*NewViewRowAccessPolicyRequest(rowAccessPolicyId, []Column{{"c"}, {"d"}})).
			WithAggregationPolicy(*NewViewAggregationPolicyRequest(aggregationPolicyId).WithEntityKey([]Column{{"column_with_comment"}})).
			WithTag([]TagAssociation{{
				Name:  tag2Id,
				Value: "v2",
			}})

		assertOptsValidAndSQLEquals(t, req.toOpts(), `CREATE OR REPLACE SECURE TEMPORARY RECURSIVE VIEW %s `+
			`("column_without_comment", "column_with_comment" COMMENT 'column 2 comment', "column" MASKING POLICY %s USING ("a", "b") TAG (%s = 'v1'), "column 2" PROJECTION POLICY %s) COPY GRANTS COMMENT = 'comment' ROW ACCESS POLICY %s ON ("c", "d") AGGREGATION POLICY %s ENTITY KEY ("column_with_comment") TAG (%s = 'v2') AS %s`, id.FullyQualifiedName(), maskingPolicy1Id.FullyQualifiedName(), tag1Id.FullyQualifiedName(), maskingPolicy2Id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(), aggregationPolicyId.FullyQualifiedName(), tag2Id.FullyQualifiedName(), sql)
	})
}

func TestViews_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

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
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetComment opts.UnsetComment opts.SetSecure opts.SetChangeTracking opts.UnsetSecure opts.SetTags opts.UnsetTags opts.AddDataMetricFunction opts.DropDataMetricFunction opts.AddRowAccessPolicy opts.DropRowAccessPolicy opts.DropAndAddRowAccessPolicy opts.DropAllRowAccessPolicies opts.SetMaskingPolicyOnColumn opts.UnsetMaskingPolicyOnColumn opts.SetTagsOnColumn opts.UnsetTagsOnColumn] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterViewOptions", "RenameTo", "SetComment", "UnsetComment", "SetSecure", "SetChangeTracking", "UnsetSecure", "SetTags", "UnsetTags", "AddDataMetricFunction", "DropDataMetricFunction", "SetDataMetricSchedule", "UnsetDataMetricSchedule", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "SetAggregationPolicy", "UnsetAggregationPolicy", "SetMaskingPolicyOnColumn", "UnsetMaskingPolicyOnColumn", "SetProjectionPolicyOnColumn", "UnsetProjectionPolicyOnColumn", "SetTagsOnColumn", "UnsetTagsOnColumn"))
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetComment opts.UnsetComment opts.SetSecure opts.SetChangeTracking opts.UnsetSecure opts.SetTags opts.UnsetTags opts.AddDataMetricFunction opts.DropDataMetricFunction opts.AddRowAccessPolicy opts.DropRowAccessPolicy opts.DropAndAddRowAccessPolicy opts.DropAllRowAccessPolicies opts.SetMaskingPolicyOnColumn opts.UnsetMaskingPolicyOnColumn opts.SetTagsOnColumn opts.UnsetTagsOnColumn] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetChangeTracking = Bool(true)
		opts.DropAllRowAccessPolicies = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterViewOptions", "RenameTo", "SetComment", "UnsetComment", "SetSecure", "SetChangeTracking", "UnsetSecure", "SetTags", "UnsetTags", "AddDataMetricFunction", "DropDataMetricFunction", "SetDataMetricSchedule", "UnsetDataMetricSchedule", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "SetAggregationPolicy", "UnsetAggregationPolicy", "SetMaskingPolicyOnColumn", "UnsetMaskingPolicyOnColumn", "SetProjectionPolicyOnColumn", "UnsetProjectionPolicyOnColumn", "SetTagsOnColumn", "UnsetTagsOnColumn"))
	})

	t.Run("validation: exactly one field from [opts.SetDataMetricSchedule.UsingCron opts.SetDataMetricSchedule.TriggerOnChanges opts.SetDataMetricSchedule.Minutes] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetDataMetricSchedule = &ViewSetDataMetricSchedule{
			UsingCron: &ViewUsingCron{
				Cron: "5 * * * * UTC",
			},
			TriggerOnChanges: Pointer(true),
		}

		opts.DropAllRowAccessPolicies = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterViewOptions.SetDataMetricSchedule", "Minutes", "UsingCron", "TriggerOnChanges"))
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.SetSecure]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.SetSecure = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterViewOptions", "IfExists", "SetSecure"))
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.UnsetSecure]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.UnsetSecure = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterViewOptions", "IfExists", "UnsetSecure"))
	})

	t.Run("validation: valid identifier for [opts.DropRowAccessPolicy.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropRowAccessPolicy = &ViewDropRowAccessPolicy{
			RowAccessPolicy: emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.AddRowAccessPolicy.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddRowAccessPolicy = &ViewAddRowAccessPolicy{
			RowAccessPolicy: emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: empty columns for row access policy (add)", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddRowAccessPolicy = &ViewAddRowAccessPolicy{
			RowAccessPolicy: randomSchemaObjectIdentifier(),
			On:              []Column{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("AlterViewOptions.AddRowAccessPolicy", "On"))
	})

	t.Run("validation: valid identifier for [opts.DropAndAddRowAccessPolicy.Drop.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropAndAddRowAccessPolicy = &ViewDropAndAddRowAccessPolicy{
			Drop: ViewDropRowAccessPolicy{
				RowAccessPolicy: emptySchemaObjectIdentifier,
			},
			Add: ViewAddRowAccessPolicy{
				RowAccessPolicy: randomSchemaObjectIdentifier(),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.DropAndAddRowAccessPolicy.Add.RowAccessPolicy]", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropAndAddRowAccessPolicy = &ViewDropAndAddRowAccessPolicy{
			Drop: ViewDropRowAccessPolicy{
				RowAccessPolicy: randomSchemaObjectIdentifier(),
			},
			Add: ViewAddRowAccessPolicy{
				RowAccessPolicy: emptySchemaObjectIdentifier,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: empty columns for row access policy (drop and add)", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropAndAddRowAccessPolicy = &ViewDropAndAddRowAccessPolicy{
			Drop: ViewDropRowAccessPolicy{
				RowAccessPolicy: randomSchemaObjectIdentifier(),
			},
			Add: ViewAddRowAccessPolicy{
				RowAccessPolicy: randomSchemaObjectIdentifier(),
				On:              []Column{},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("AlterViewOptions.DropAndAddRowAccessPolicy.Add", "On"))
	})

	t.Run("rename", func(t *testing.T) {
		newId := randomSchemaObjectIdentifier()

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
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s SET CHANGE_TRACKING = true", id.FullyQualifiedName())
	})

	t.Run("set change tracking: false", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetChangeTracking = Bool(false)
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s SET CHANGE_TRACKING = false", id.FullyQualifiedName())
	})

	t.Run("unset secure", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetSecure = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s UNSET SECURE", id.FullyQualifiedName())
	})

	t.Run("add data metric function", func(t *testing.T) {
		dmfId := randomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.AddDataMetricFunction = &ViewAddDataMetricFunction{
			DataMetricFunction: []ViewDataMetricFunction{{DataMetricFunction: dmfId, On: []Column{{"foo"}}}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s ADD DATA METRIC FUNCTION %s ON (\"foo\")", id.FullyQualifiedName(), dmfId.FullyQualifiedName())
	})

	t.Run("drop data metric function", func(t *testing.T) {
		dmfId := randomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.DropDataMetricFunction = &ViewDropDataMetricFunction{
			DataMetricFunction: []ViewDataMetricFunction{{DataMetricFunction: dmfId, On: []Column{{"foo"}}}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s DROP DATA METRIC FUNCTION %s ON (\"foo\")", id.FullyQualifiedName(), dmfId.FullyQualifiedName())
	})

	t.Run("set data metric schedule", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetDataMetricSchedule = &ViewSetDataMetricSchedule{
			Minutes: &ViewMinute{
				Minutes: 5,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s SET DATA_METRIC_SCHEDULE = ' 5 MINUTE'", id.FullyQualifiedName())

		opts = defaultOpts()
		opts.SetDataMetricSchedule = &ViewSetDataMetricSchedule{
			UsingCron: &ViewUsingCron{
				Cron: "5 * * * * UTC",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s SET DATA_METRIC_SCHEDULE = 'USING CRON 5 * * * * UTC '", id.FullyQualifiedName())

		opts = defaultOpts()
		opts.SetDataMetricSchedule = &ViewSetDataMetricSchedule{
			TriggerOnChanges: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s SET DATA_METRIC_SCHEDULE = 'TRIGGER_ON_CHANGES'", id.FullyQualifiedName())
	})

	t.Run("unset data metric schedule", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetDataMetricSchedule = &ViewUnsetDataMetricSchedule{}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s UNSET DATA_METRIC_SCHEDULE", id.FullyQualifiedName())
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
		rowAccessPolicyId := randomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.AddRowAccessPolicy = &ViewAddRowAccessPolicy{
			RowAccessPolicy: rowAccessPolicyId,
			On:              []Column{{"a"}, {"b"}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s ADD ROW ACCESS POLICY %s ON (\"a\", \"b\")", id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("drop row access policy", func(t *testing.T) {
		rowAccessPolicyId := randomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.DropRowAccessPolicy = &ViewDropRowAccessPolicy{
			RowAccessPolicy: rowAccessPolicyId,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s DROP ROW ACCESS POLICY %s", id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("drop and add row access policy", func(t *testing.T) {
		rowAccessPolicy1Id := randomSchemaObjectIdentifier()
		rowAccessPolicy2Id := randomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.DropAndAddRowAccessPolicy = &ViewDropAndAddRowAccessPolicy{
			Drop: ViewDropRowAccessPolicy{
				RowAccessPolicy: rowAccessPolicy1Id,
			},
			Add: ViewAddRowAccessPolicy{
				RowAccessPolicy: rowAccessPolicy2Id,
				On:              []Column{{"a"}, {"b"}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s DROP ROW ACCESS POLICY %s, ADD ROW ACCESS POLICY %s ON (\"a\", \"b\")", id.FullyQualifiedName(), rowAccessPolicy1Id.FullyQualifiedName(), rowAccessPolicy2Id.FullyQualifiedName())
	})

	t.Run("drop all row access policies", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropAllRowAccessPolicies = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s DROP ALL ROW ACCESS POLICIES", id.FullyQualifiedName())
	})

	t.Run("set aggregation policy", func(t *testing.T) {
		aggregationPolicyId := randomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.SetAggregationPolicy = &ViewSetAggregationPolicy{
			AggregationPolicy: aggregationPolicyId,
			EntityKey:         []Column{{"a"}, {"b"}},
			Force:             Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s SET AGGREGATION POLICY %s ENTITY KEY (\"a\", \"b\") FORCE", id.FullyQualifiedName(), aggregationPolicyId.FullyQualifiedName())
	})

	t.Run("unset aggregation policy", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetAggregationPolicy = &ViewUnsetAggregationPolicy{}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s UNSET AGGREGATION POLICY", id.FullyQualifiedName())
	})

	t.Run("set masking policy on column", func(t *testing.T) {
		maskingPolicyId := randomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.SetMaskingPolicyOnColumn = &ViewSetColumnMaskingPolicy{
			Name:          "column",
			MaskingPolicy: maskingPolicyId,
			Using:         []Column{{"a"}, {"b"}},
			Force:         Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s ALTER COLUMN \"column\" SET MASKING POLICY %s USING (\"a\", \"b\") FORCE", id.FullyQualifiedName(), maskingPolicyId.FullyQualifiedName())
	})

	t.Run("unset masking policy on column", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetMaskingPolicyOnColumn = &ViewUnsetColumnMaskingPolicy{
			Name: "column",
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s ALTER COLUMN \"column\" UNSET MASKING POLICY", id.FullyQualifiedName())
	})

	t.Run("set projection policy on column", func(t *testing.T) {
		projectionPolicyId := randomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.SetProjectionPolicyOnColumn = &ViewSetProjectionPolicy{
			Name:             "column",
			ProjectionPolicy: projectionPolicyId,
			Force:            Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s ALTER COLUMN \"column\" SET PROJECTION POLICY %s FORCE", id.FullyQualifiedName(), projectionPolicyId.FullyQualifiedName())
	})

	t.Run("unset projection policy on column", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetProjectionPolicyOnColumn = &ViewUnsetProjectionPolicy{
			Name: "column",
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER VIEW %s ALTER COLUMN \"column\" UNSET PROJECTION POLICY", id.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER VIEW %s ALTER COLUMN "column" SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER VIEW %s ALTER COLUMN "column" UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})
}

func TestViews_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

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
		opts.name = emptySchemaObjectIdentifier
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

	t.Run("in database", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.In = &ExtendedIn{
			In: In{
				Database: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, fmt.Sprintf("SHOW VIEWS IN DATABASE %s", id.FullyQualifiedName()))
	})

	t.Run("in schema", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := defaultOpts()
		opts.In = &ExtendedIn{
			In: In{
				Schema: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, fmt.Sprintf("SHOW VIEWS IN SCHEMA %s", id.FullyQualifiedName()))
	})

	t.Run("in application", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.In = &ExtendedIn{
			Application: id,
		}
		assertOptsValidAndSQLEquals(t, opts, fmt.Sprintf("SHOW VIEWS IN APPLICATION %s", id.FullyQualifiedName()))
	})

	t.Run("in application package", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.In = &ExtendedIn{
			ApplicationPackage: id,
		}
		assertOptsValidAndSQLEquals(t, opts, fmt.Sprintf("SHOW VIEWS IN APPLICATION PACKAGE %s", id.FullyQualifiedName()))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		opts.Like = &Like{
			Pattern: String("myaccount"),
		}
		opts.In = &ExtendedIn{
			In: In{
				Account: Bool(true),
			},
		}
		opts.StartsWith = String("abc")
		opts.Limit = &LimitFrom{Rows: Int(10)}
		assertOptsValidAndSQLEquals(t, opts, "SHOW TERSE VIEWS LIKE 'myaccount' IN ACCOUNT STARTS WITH 'abc' LIMIT 10")
	})
}

func TestViews_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

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
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE VIEW %s", id.FullyQualifiedName())
	})
}
