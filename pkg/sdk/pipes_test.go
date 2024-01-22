package sdk

import (
	"testing"
)

func TestPipesCreate(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreatePipeOptions {
		return &CreatePipeOptions{
			name:          id,
			copyStatement: "<copy_statement>",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreatePipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: copy statement required", func(t *testing.T) {
		opts := defaultOpts()
		opts.copyStatement = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreatePipeOptions", "copyStatement"))
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
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *AlterPipeOptions {
		return &AlterPipeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterPipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterPipeOptions", "Set", "Unset", "SetTag", "UnsetTag", "Refresh"))
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PipeSet{
			ErrorIntegration: String("new_error_integration"),
		}
		opts.Unset = &PipeUnset{
			Comment: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterPipeOptions", "Set", "Unset", "SetTag", "UnsetTag", "Refresh"))
	})

	t.Run("validation: no property to set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PipeSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterPipeOptions.Set", "ErrorIntegration", "PipeExecutionPaused", "Comment"))
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &PipeUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterPipeOptions.Unset", "ErrorIntegration", "PipeExecutionPaused", "Comment"))
	})

	t.Run("set tag: single", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTag = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag_name1"),
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s SET TAG "tag_name1" = 'v1'`, id.FullyQualifiedName())
	})

	t.Run("set tag: multiple", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTag = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag_name1"),
				Value: "v1",
			},
			{
				Name:  NewAccountObjectIdentifier("tag_name2"),
				Value: "v2",
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
		opts.UnsetTag = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag_name1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s UNSET TAG "tag_name1"`, id.FullyQualifiedName())
	})

	t.Run("unset tag: multi", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTag = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag_name1"),
			NewAccountObjectIdentifier("tag_name2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s UNSET TAG "tag_name1", "tag_name2"`, id.FullyQualifiedName())
	})

	t.Run("unset all", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &PipeUnset{
			ErrorIntegration:    Bool(true),
			PipeExecutionPaused: Bool(true),
			Comment:             Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE IF EXISTS %s UNSET ERROR_INTEGRATION, PIPE_EXECUTION_PAUSED, COMMENT`, id.FullyQualifiedName())
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
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *DropPipeOptions {
		return &DropPipeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropPipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
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
	id := RandomSchemaObjectIdentifier()
	databaseIdentifier := NewAccountObjectIdentifier(id.DatabaseName())
	schemaIdentifier := NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())

	defaultOpts := func() *ShowPipeOptions {
		return &ShowPipeOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowPipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
	})

	t.Run("validation: empty in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ShowPipeOptions.In", "Account", "Database", "Schema"))
	})

	t.Run("validation: exactly one scope for in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account:  Bool(true),
			Database: databaseIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ShowPipeOptions.In", "Account", "Database", "Schema"))
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
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *describePipeOptions {
		return &describePipeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *describePipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("with name", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE PIPE %s`, id.FullyQualifiedName())
	})
}
