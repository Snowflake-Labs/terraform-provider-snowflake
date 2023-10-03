package sdk

import (
	"testing"
)

func TestPipesCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *CreatePipeOptions {
		return &CreatePipeOptions{
			name:          id,
			copyStatement: "<copy_statement>",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreatePipeOptions = nil
		assertOptsInvalid(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalid(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: copy statement required", func(t *testing.T) {
		opts := defaultOpts()
		opts.copyStatement = ""
		assertOptsInvalid(t, opts, errCopyStatementRequired)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE PIPE %s AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.AutoIngest = Bool(true)
		opts.ErrorIntegration = String("some_error_integration")
		opts.AwsSnsTopic = String("some aws sns topic")
		opts.Integration = String("some integration")
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE PIPE IF NOT EXISTS %s AUTO_INGEST = true ERROR_INTEGRATION = some_error_integration AWS_SNS_TOPIC = 'some aws sns topic' INTEGRATION = 'some integration' COMMENT = 'some comment' AS <copy_statement>`, id.FullyQualifiedName())
	})
}

func TestPipesAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *AlterPipeOptions {
		return &AlterPipeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterPipeOptions = nil
		assertOptsInvalid(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalid(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalid(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PipeSet{
			ErrorIntegration: String("new_error_integration"),
		}
		opts.Unset = &PipeUnset{
			Comment: Bool(true),
		}
		assertOptsInvalid(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: no property to set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PipeSet{}
		assertOptsInvalid(t, opts, errAlterNeedsAtLeastOneProperty)
	})

	t.Run("validation: empty tags slice for set", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = &PipeSetTags{
			Tag: []TagAssociation{},
		}
		assertOptsInvalid(t, opts, errAlterNeedsAtLeastOneProperty)
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &PipeUnset{}
		assertOptsInvalid(t, opts, errAlterNeedsAtLeastOneProperty)
	})

	t.Run("validation: empty tags slice for unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = &PipeUnsetTags{
			Tag: []ObjectIdentifier{},
		}
		assertOptsInvalid(t, opts, errAlterNeedsAtLeastOneProperty)
	})

	t.Run("set tag: single", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = &PipeSetTags{
			Tag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag_name1"),
					Value: "v1",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s SET TAG "tag_name1" = 'v1'`, id.FullyQualifiedName())
	})

	t.Run("set tag: multiple", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = &PipeSetTags{
			Tag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag_name1"),
					Value: "v1",
				},
				{
					Name:  NewAccountObjectIdentifier("tag_name2"),
					Value: "v2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s SET TAG "tag_name1" = 'v1', "tag_name2" = 'v2'`, id.FullyQualifiedName())
	})

	t.Run("set all", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Set = &PipeSet{
			ErrorIntegration:    String("new_error_integration"),
			PipeExecutionPaused: Bool(true),
			Comment:             String("new comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE IF EXISTS %s SET ERROR_INTEGRATION = new_error_integration, PIPE_EXECUTION_PAUSED = true, COMMENT = 'new comment'`, id.FullyQualifiedName())
	})

	t.Run("unset tag: single", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = &PipeUnsetTags{
			Tag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag_name1"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s UNSET TAG "tag_name1"`, id.FullyQualifiedName())
	})

	t.Run("unset tag: multi", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = &PipeUnsetTags{
			Tag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag_name1"),
				NewAccountObjectIdentifier("tag_name2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s UNSET TAG "tag_name1", "tag_name2"`, id.FullyQualifiedName())
	})

	t.Run("unset all", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &PipeUnset{
			PipeExecutionPaused: Bool(true),
			Comment:             Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE IF EXISTS %s UNSET PIPE_EXECUTION_PAUSED, COMMENT`, id.FullyQualifiedName())
	})

	t.Run("refresh", func(t *testing.T) {
		opts := defaultOpts()
		opts.Refresh = &PipeRefresh{}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s REFRESH`, id.FullyQualifiedName())
	})

	t.Run("refresh with all", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Refresh = &PipeRefresh{
			Prefix:        String("/d1"),
			ModifiedAfter: String("2018-07-30T13:56:46-07:00"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE IF EXISTS %s REFRESH PREFIX = '/d1' MODIFIED_AFTER = '2018-07-30T13:56:46-07:00'`, id.FullyQualifiedName())
	})
}

func TestPipesDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *DropPipeOptions {
		return &DropPipeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropPipeOptions = nil
		assertOptsInvalid(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalid(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP PIPE %s`, id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP PIPE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestPipesShow(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)
	databaseIdentifier := NewAccountObjectIdentifier(id.DatabaseName())
	schemaIdentifier := NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())

	defaultOpts := func() *ShowPipeOptions {
		return &ShowPipeOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowPipeOptions = nil
		assertOptsInvalid(t, opts, errNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalid(t, opts, errPatternRequiredForLikeKeyword)
	})

	t.Run("validation: empty in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{}
		assertOptsInvalid(t, opts, errScopeRequiredForInKeyword)
	})

	t.Run("validation: exactly one scope for in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account:  Bool(true),
			Database: databaseIdentifier,
		}
		assertOptsInvalid(t, opts, errScopeRequiredForInKeyword)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES`)
	})

	t.Run("with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES LIKE '%s'`, id.Name())
	})

	t.Run("in account", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES IN ACCOUNT`)
	})

	t.Run("in database", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Database: databaseIdentifier,
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES IN DATABASE %s`, databaseIdentifier.FullyQualifiedName())
	})

	t.Run("in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Schema: schemaIdentifier,
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES IN SCHEMA %s`, schemaIdentifier.FullyQualifiedName())
	})

	t.Run("with like and in account", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Account: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES LIKE '%s' IN ACCOUNT`, id.Name())
	})

	t.Run("with like and in database", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Database: databaseIdentifier,
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES LIKE '%s' IN DATABASE %s`, id.Name(), databaseIdentifier.FullyQualifiedName())
	})

	t.Run("with like and in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Schema: schemaIdentifier,
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES LIKE '%s' IN SCHEMA %s`, id.Name(), schemaIdentifier.FullyQualifiedName())
	})
}

func TestPipesDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *describePipeOptions {
		return &describePipeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *describePipeOptions = nil
		assertOptsInvalid(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalid(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("with name", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE PIPE %s`, id.FullyQualifiedName())
	})
}
