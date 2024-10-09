package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreams_CreateOnTable(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	tableId := randomSchemaObjectIdentifier()

	// Minimal valid CreateOnTableStreamOptions
	defaultOpts := func() *CreateOnTableStreamOptions {
		return &CreateOnTableStreamOptions{
			name:    id,
			TableId: tableId,
			On: &OnStream{
				At: Bool(true),
				Statement: OnStreamStatement{
					Stream: String("123"),
				},
			},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnTableStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.TableId]", func(t *testing.T) {
		opts := defaultOpts()
		opts.TableId = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnTableStreamOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: exactly one field from [opts.On.At opts.On.Before] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.On.At = Bool(true)
		opts.On.Before = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateOnTableStreamOptions.On", "At", "Before"))
	})

	t.Run("validation: exactly one field from [opts.On.Statement.Timestamp opts.On.Statement.Offset opts.On.Statement.Statement opts.On.Statement.Stream] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.On.At = Bool(true)
		opts.On.Statement = OnStreamStatement{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateOnTableStreamOptions.On.Statement", "Timestamp", "Offset", "Statement", "Stream"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = nil
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON TABLE %s", id.FullyQualifiedName(), tableId.FullyQualifiedName())
	})

	t.Run("at timestamp", func(t *testing.T) {
		timestamp := "2024-09-25 06:16:10.359 -0700"
		opts := defaultOpts()
		opts.On = &OnStream{
			At: Bool(true),
			Statement: OnStreamStatement{
				Timestamp: String(timestamp),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON TABLE %s AT (TIMESTAMP => '%s')", id.FullyQualifiedName(), tableId.FullyQualifiedName(), timestamp)
	})

	t.Run("at offset", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = &OnStream{
			At: Bool(true),
			Statement: OnStreamStatement{
				Offset: String("-10"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON TABLE %s AT (OFFSET => -10)", id.FullyQualifiedName(), tableId.FullyQualifiedName())
	})

	t.Run("at statement", func(t *testing.T) {
		queryId := "0111447d-0905-8a5c-0062-f3820281547a"
		opts := defaultOpts()
		opts.On = &OnStream{
			At: Bool(true),
			Statement: OnStreamStatement{
				Statement: String(queryId),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON TABLE %s AT (STATEMENT => '%s')", id.FullyQualifiedName(), tableId.FullyQualifiedName(), queryId)
	})

	t.Run("at stream", func(t *testing.T) {
		streamId := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.On = &OnStream{
			At: Bool(true),
			Statement: OnStreamStatement{
				Stream: String(streamId.FullyQualifiedName()),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON TABLE %s AT (STREAM => '%s')", id.FullyQualifiedName(), tableId.FullyQualifiedName(), temporaryReplace(streamId))
	})

	t.Run("before timestamp", func(t *testing.T) {
		timestamp := "2024-09-25 06:16:10.359 -0700"
		opts := defaultOpts()
		opts.On = &OnStream{
			Before: Bool(true),
			Statement: OnStreamStatement{
				Timestamp: String(timestamp),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON TABLE %s BEFORE (TIMESTAMP => '%s')", id.FullyQualifiedName(), tableId.FullyQualifiedName(), timestamp)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.On = &OnStream{
			At: Bool(true),
			Statement: OnStreamStatement{
				Stream: String("123"),
			},
		}
		opts.AppendOnly = Bool(true)
		opts.ShowInitialRows = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE STREAM %s ON TABLE %s AT (STREAM => '123') APPEND_ONLY = true SHOW_INITIAL_ROWS = true COMMENT = 'some comment'", id.FullyQualifiedName(), tableId.FullyQualifiedName())
	})
}

func TestStreams_CreateOnExternalTable(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	externalTableId := randomSchemaObjectIdentifier()

	// Minimal valid CreateOnExternalTableStreamOptions
	defaultOpts := func() *CreateOnExternalTableStreamOptions {
		return &CreateOnExternalTableStreamOptions{
			name:            id,
			ExternalTableId: externalTableId,
			On: &OnStream{
				At: Bool(true),
				Statement: OnStreamStatement{
					Stream: String("123"),
				},
			},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnExternalTableStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.ExternalTableId]", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalTableId = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnExternalTableStreamOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: exactly one field from [opts.On.At opts.On.Before] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.On.At = Bool(true)
		opts.On.Before = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateOnExternalTableStreamOptions.On", "At", "Before"))
	})

	t.Run("validation: exactly one field from [opts.On.Statement.Timestamp opts.On.Statement.Offset opts.On.Statement.Statement opts.On.Statement.Stream] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.On.Statement = OnStreamStatement{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateOnExternalTableStreamOptions.On.Statement", "Timestamp", "Offset", "Statement", "Stream"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = nil
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON EXTERNAL TABLE %s", id.FullyQualifiedName(), externalTableId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.CopyGrants = Bool(true)
		opts.On = &OnStream{
			At: Bool(true),
			Statement: OnStreamStatement{
				Statement: String("123"),
			},
		}
		opts.InsertOnly = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE STREAM IF NOT EXISTS %s COPY GRANTS ON EXTERNAL TABLE %s AT (STATEMENT => '123') INSERT_ONLY = true COMMENT = 'some comment'`, id.FullyQualifiedName(), externalTableId.FullyQualifiedName())
	})
}

func TestStreams_CreateOnDirectoryTable(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	stageId := randomSchemaObjectIdentifier()

	// Minimal valid CreateOnStageStreamOptions
	defaultOpts := func() *CreateOnDirectoryTableStreamOptions {
		return &CreateOnDirectoryTableStreamOptions{
			name:    id,
			StageId: stageId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnDirectoryTableStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.StageId]", func(t *testing.T) {
		opts := defaultOpts()
		opts.StageId = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnDirectoryTableStreamOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON STAGE %s", id.FullyQualifiedName(), stageId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.CopyGrants = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE STREAM IF NOT EXISTS %s COPY GRANTS ON STAGE %s COMMENT = 'some comment'`, id.FullyQualifiedName(), stageId.FullyQualifiedName())
	})
}

func TestStreams_CreateOnView(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	viewId := randomSchemaObjectIdentifier()

	// Minimal valid CreateOnViewStreamOptions
	defaultOpts := func() *CreateOnViewStreamOptions {
		return &CreateOnViewStreamOptions{
			name:   id,
			ViewId: viewId,
			On: &OnStream{
				At: Bool(true),
				Statement: OnStreamStatement{
					Stream: String("123"),
				},
			},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnViewStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.viewId]", func(t *testing.T) {
		opts := defaultOpts()
		opts.ViewId = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnViewStreamOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: exactly one field from [opts.On.At opts.On.Before] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.On.At = Bool(true)
		opts.On.Before = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateOnViewStreamOptions.On", "At", "Before"))
	})

	t.Run("validation: exactly one field from [opts.On.Statement.Timestamp opts.On.Statement.Offset opts.On.Statement.Statement opts.On.Statement.Stream] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.On.Statement = OnStreamStatement{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateOnViewStreamOptions.On.Statement", "Timestamp", "Offset", "Statement", "Stream"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = nil
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s ON VIEW %s", id.FullyQualifiedName(), viewId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.CopyGrants = Bool(true)
		opts.On = &OnStream{
			Before: Bool(true),
			Statement: OnStreamStatement{
				Stream: String("123"),
			},
		}
		opts.AppendOnly = Bool(true)
		opts.ShowInitialRows = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE STREAM %s COPY GRANTS ON VIEW %s BEFORE (STREAM => '123') APPEND_ONLY = true SHOW_INITIAL_ROWS = true COMMENT = 'some comment'`, id.FullyQualifiedName(), viewId.FullyQualifiedName())
	})
}

func TestStreams_Clone(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	sourceId := randomSchemaObjectIdentifier()

	// Minimal valid CloneStreamOptions
	defaultOpts := func() *CloneStreamOptions {
		return &CloneStreamOptions{
			name:         id,
			sourceStream: sourceId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CloneStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE STREAM %s CLONE %s", id.FullyQualifiedName(), sourceId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.CopyGrants = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE STREAM %s CLONE %s COPY GRANTS", id.FullyQualifiedName(), sourceId.FullyQualifiedName())
	})
}

func TestStreams_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid AlterStreamOptions
	defaultOpts := func() *AlterStreamOptions {
		return &AlterStreamOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.UnsetTags]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.UnsetTags = []ObjectIdentifier{randomAccountObjectIdentifier()}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterStreamOptions", "IfExists", "UnsetTags"))
	})

	t.Run("validation: exactly one field from [opts.SetComment opts.UnsetComment opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterStreamOptions", "SetComment", "UnsetComment", "SetTags", "UnsetTags"))
	})

	t.Run("set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.SetComment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `ALTER STREAM IF EXISTS %s SET COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.UnsetComment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER STREAM IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER STREAM %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})
}

func TestStreams_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid DropStreamOptions
	defaultOpts := func() *DropStreamOptions {
		return &DropStreamOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP STREAM %s`, id.FullyQualifiedName())
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
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW STREAMS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		opts.Like = &Like{Pattern: String("pattern")}
		schemaId := randomDatabaseObjectIdentifier()
		opts.In = &ExtendedIn{In: In{Schema: schemaId}}
		opts.StartsWith = String("starts with pattern")
		opts.Limit = &LimitFrom{Rows: Int(123), From: String("from pattern")}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE STREAMS LIKE 'pattern' IN SCHEMA %s STARTS WITH 'starts with pattern' LIMIT 123 FROM 'from pattern'`, schemaId.FullyQualifiedName())
	})
}

func TestStreams_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid DescribeStreamOptions
	defaultOpts := func() *DescribeStreamOptions {
		return &DescribeStreamOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeStreamOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("valid sql", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE STREAM %s`, id.FullyQualifiedName())
	})
}

func TestToStreamSourceType(t *testing.T) {
	tests := []struct {
		input   string
		want    StreamSourceType
		wantErr string
	}{
		{
			input: "TABLE",
			want:  StreamSourceTypeTable,
		},
		{
			input: "EXTERNAL TABLE",
			want:  StreamSourceTypeExternalTable,
		},
		{
			input: "VIEW",
			want:  StreamSourceTypeView,
		},
		{
			input: "STAGE",
			want:  StreamSourceTypeStage,
		},
		{
			input: "table",
			want:  StreamSourceTypeTable,
		},
		{
			input: "external table",
			want:  StreamSourceTypeExternalTable,
		},
		{
			input: "view",
			want:  StreamSourceTypeView,
		},
		{
			input: "stage",
			want:  StreamSourceTypeStage,
		},
		{
			input:   "",
			wantErr: "invalid stream source type",
		},
		{
			input:   "foo",
			wantErr: "invalid stream source type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ToStreamSourceType(tt.input)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestToStreamMode(t *testing.T) {
	tests := []struct {
		input   string
		want    StreamMode
		wantErr string
	}{
		{
			input: "DEFAULT",
			want:  StreamModeDefault,
		},
		{
			input: "APPEND_ONLY",
			want:  StreamModeAppendOnly,
		},
		{
			input: "INSERT_ONLY",
			want:  StreamModeInsertOnly,
		},
		{
			input: "default",
			want:  StreamModeDefault,
		},
		{
			input: "append_only",
			want:  StreamModeAppendOnly,
		},
		{
			input: "insert_only",
			want:  StreamModeInsertOnly,
		},
		{
			input:   "",
			wantErr: "invalid stream mode",
		},
		{
			input:   "foo",
			wantErr: "invalid stream mode",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ToStreamMode(tt.input)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
