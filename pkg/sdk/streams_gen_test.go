package sdk

import "testing"

func TestStreams_CreateOnTable(t *testing.T) {
	id := randomAccountObjectIdentifier(t)
	tableId := randomAccountObjectIdentifier(t)

	// Minimal valid CreateOnTableStreamOptions
	defaultOpts := func() *CreateOnTableStreamOptions {
		return &CreateOnTableStreamOptions{
			name:    id,
			TableId: tableId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnTableStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.TableId]", func(t *testing.T) {
		opts := defaultOpts()
		opts.TableId = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.CloneStream] if set", func(t *testing.T) {
		opts := defaultOpts()
		invalidId := NewAccountObjectIdentifier("")
		opts.CloneStream = &invalidId
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.CloneStream]", func(t *testing.T) {
		opts := defaultOpts()
		streamId := randomAccountObjectIdentifier(t)
		opts.IfNotExists = Bool(true)
		opts.CloneStream = &streamId
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnTableStreamOptions", "IfNotExists", "CloneStream"))
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnTableStreamOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON TABLE %s", id.FullyQualifiedName(), tableId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.On = &OnStream{
			At:     Bool(true),
			Stream: String("123"),
		}
		opts.AppendOnly = true
		opts.ShowInitialRows = true
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE STREAM %s ON TABLE %s AT (STREAM => '123') APPEND_ONLY = TRUE SHOW_INITIAL_ROWS = TRUE COMMENT = 'some comment'", id.FullyQualifiedName(), tableId.FullyQualifiedName())
	})
}

func TestStreams_CreateOnExternalTable(t *testing.T) {
	id := randomAccountObjectIdentifier(t)
	externalTableId := randomAccountObjectIdentifier(t)

	// Minimal valid CreateOnExternalTableStreamOptions
	defaultOpts := func() *CreateOnExternalTableStreamOptions {
		return &CreateOnExternalTableStreamOptions{
			name:            id,
			ExternalTableId: externalTableId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnExternalTableStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.ExternalTableId]", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalTableId = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.CloneStream] if set", func(t *testing.T) {
		opts := defaultOpts()
		invalidId := NewAccountObjectIdentifier("")
		opts.CloneStream = &invalidId
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.CloneStream]", func(t *testing.T) {
		opts := defaultOpts()
		streamId := randomAccountObjectIdentifier(t)
		opts.IfNotExists = Bool(true)
		opts.CloneStream = &streamId
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnExternalTableStreamOptions", "IfNotExists", "CloneStream"))
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnExternalTableStreamOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON EXTERNAL TABLE %s", id.FullyQualifiedName(), externalTableId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.copyGrants = true
		opts.On = &OnStream{
			At:        Bool(true),
			Statement: String("123"),
		}
		opts.InsertOnly = true
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE STREAM IF NOT EXISTS %s ON EXTERNAL TABLE %s AT (STATEMENT => "123") INSERT_ONLY = TRUE COMMENT = 'some comment'`, id.FullyQualifiedName(), externalTableId.FullyQualifiedName())
	})
}

func TestStreams_CreateOnStage(t *testing.T) {
	id := randomAccountObjectIdentifier(t)
	stageId := randomAccountObjectIdentifier(t)

	// Minimal valid CreateOnStageStreamOptions
	defaultOpts := func() *CreateOnStageStreamOptions {
		return &CreateOnStageStreamOptions{
			name:    id,
			StageId: stageId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnStageStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.StageId]", func(t *testing.T) {
		opts := defaultOpts()
		opts.StageId = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.CloneStream] if set", func(t *testing.T) {
		opts := defaultOpts()
		invalidId := NewAccountObjectIdentifier("")
		opts.CloneStream = &invalidId
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.CloneStream]", func(t *testing.T) {
		opts := defaultOpts()
		streamId := randomAccountObjectIdentifier(t)
		opts.IfNotExists = Bool(true)
		opts.CloneStream = &streamId
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnStageStreamOptions", "IfNotExists", "CloneStream"))
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnStageStreamOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON EXTERNAL TABLE %s", id.FullyQualifiedName(), stageId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.copyGrants = true
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE STREAM IF NOT EXISTS %s ON STAGE %s AT COMMENT = 'some comment'`, id.FullyQualifiedName(), stageId.FullyQualifiedName())
	})
}

func TestStreams_CreateOnView(t *testing.T) {
	id := randomAccountObjectIdentifier(t)
	viewId := randomAccountObjectIdentifier(t)

	// Minimal valid CreateOnViewStreamOptions
	defaultOpts := func() *CreateOnViewStreamOptions {
		return &CreateOnViewStreamOptions{
			name:   id,
			ViewId: viewId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnViewStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.viewId]", func(t *testing.T) {
		opts := defaultOpts()
		opts.ViewId = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.CloneStream] if set", func(t *testing.T) {
		opts := defaultOpts()
		invalidId := NewAccountObjectIdentifier("")
		opts.CloneStream = &invalidId
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.CloneStream]", func(t *testing.T) {
		opts := defaultOpts()
		streamId := randomAccountObjectIdentifier(t)
		opts.IfNotExists = Bool(true)
		opts.CloneStream = &streamId
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnExternalTableStreamOptions", "IfNotExists", "CloneStream"))
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnExternalTableStreamOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON VIEW %s", id.FullyQualifiedName(), viewId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.On = &OnStream{
			At:     Bool(true),
			Stream: String("123"),
		}
		opts.AppendOnly = true
		opts.ShowInitialRows = true
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE STREAM %s ON VIEW %s BEFORE (TIMESTAMP => "") APPEND_ONLY = TRUE SHOW_INITIAL_ROWS = TRUE COMMENT = 'some comment'`, id.FullyQualifiedName(), viewId.FullyQualifiedName())
	})
}

func TestStreams_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid AlterStreamOptions
	defaultOpts := func() *AlterStreamOptions {
		return &AlterStreamOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.UnsetTags]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.UnsetTags = []ObjectIdentifier{randomAccountObjectIdentifier(t)}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterStreamOptions", "IfExists", "UnsetTags"))
	})

	t.Run("validation: exactly one field from [opts.SetComment opts.UnsetComment opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("SetComment", "UnsetComment", "SetTags", "UnsetTags"))
	})

	t.Run("set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `ALTER STREAM IF EXISTS %s SET COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetComment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER STREAM IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER STREAM IF EXISTS %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER STREAM IF EXISTS %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})
}

func TestStreams_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid DropStreamOptions
	defaultOpts := func() *DropStreamOptions {
		return &DropStreamOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP STREAM IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestStreams_Show(t *testing.T) {
	// Minimal valid ShowStreamOptions
	defaultOpts := func() *ShowStreamOptions {
		return &ShowStreamOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW STREAMS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		opts.Like = &Like{Pattern: String("pattern")}
		schemaId := randomDatabaseObjectIdentifier(t)
		opts.In = &In{Schema: schemaId}
		opts.StartsWith = String("starts with pattern")
		opts.Limit = &LimitFrom{Rows: Int(123), From: String("from pattern")}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE STREAMS LIKE 'pattern' IN SCHEMA %s STARTS WITH 'starts with pattern' LIMIT 123 FROM 'from pattern'`, schemaId.FullyQualifiedName())
	})
}

func TestStreams_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid DescribeStreamOptions
	defaultOpts := func() *DescribeStreamOptions {
		return &DescribeStreamOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("valid sql", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE STREAM %s`, id.FullyQualifiedName())
	})
}
