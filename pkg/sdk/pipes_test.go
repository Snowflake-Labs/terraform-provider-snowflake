package sdk

import (
	"testing"
)

func TestPipesCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	setUpOpts := func() *PipeCreateOptions {
		return &PipeCreateOptions{
			name:          id,
			CopyStatement: "<copy_statement>",
		}
	}

	t.Run("basic", func(t *testing.T) {
		opts := setUpOpts()
		assertSqlEquals(t, opts, `CREATE PIPE %s AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("if not exists", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfNotExists = Bool(true)
		assertSqlEquals(t, opts, `CREATE PIPE IF NOT EXISTS %s AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("auto ingest: true", func(t *testing.T) {
		opts := setUpOpts()
		opts.AutoIngest = Bool(true)
		assertSqlEquals(t, opts, `CREATE PIPE %s AUTO_INGEST = true AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("auto ingest: false", func(t *testing.T) {
		opts := setUpOpts()
		opts.AutoIngest = Bool(false)
		assertSqlEquals(t, opts, `CREATE PIPE %s AUTO_INGEST = false AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("error integration", func(t *testing.T) {
		opts := setUpOpts()
		opts.ErrorIntegration = String("some_error_integration")
		assertSqlEquals(t, opts, `CREATE PIPE %s ERROR_INTEGRATION = some_error_integration AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("aws sns topic", func(t *testing.T) {
		opts := setUpOpts()
		opts.AwsSnsTopic = String("some aws sns topic")
		assertSqlEquals(t, opts, `CREATE PIPE %s AWS_SNS_TOPIC = 'some aws sns topic' AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("integration", func(t *testing.T) {
		opts := setUpOpts()
		opts.Integration = String("some integration")
		assertSqlEquals(t, opts, `CREATE PIPE %s INTEGRATION = 'some integration' AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("comment", func(t *testing.T) {
		opts := setUpOpts()
		opts.Comment = String("some comment")
		assertSqlEquals(t, opts, `CREATE PIPE %s COMMENT = 'some comment' AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("all optional", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfNotExists = Bool(true)
		opts.AutoIngest = Bool(true)
		opts.ErrorIntegration = String("some_error_integration")
		opts.AwsSnsTopic = String("some aws sns topic")
		opts.Integration = String("some integration")
		opts.Comment = String("some comment")
		assertSqlEquals(t, opts, `CREATE PIPE IF NOT EXISTS %s AUTO_INGEST = true ERROR_INTEGRATION = some_error_integration AWS_SNS_TOPIC = 'some aws sns topic' INTEGRATION = 'some integration' COMMENT = 'some comment' AS <copy_statement>`, id.FullyQualifiedName())
	})
}

func TestPipesAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	setUpOpts := func() *PipeAlterOptions {
		return &PipeAlterOptions{
			name: id,
		}
	}

	t.Run("set error integration", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			ErrorIntegration: String("new_error_integration"),
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s SET ERROR_INTEGRATION = new_error_integration`, id.FullyQualifiedName())
	})

	t.Run("set pipe execution paused: true", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			PipeExecutionPaused: Bool(true),
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s SET PIPE_EXECUTION_PAUSED = true`, id.FullyQualifiedName())
	})

	t.Run("set pipe execution paused: false", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			PipeExecutionPaused: Bool(false),
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s SET PIPE_EXECUTION_PAUSED = false`, id.FullyQualifiedName())
	})

	t.Run("set tag: single", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			Tag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag_name1"),
					Value: "v1",
				},
			},
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s SET TAG "tag_name1" = 'v1'`, id.FullyQualifiedName())
	})

	t.Run("set tag: multiple", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
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
		assertSqlEquals(t, opts, `ALTER PIPE %s SET TAG "tag_name1" = 'v1', "tag_name2" = 'v2'`, id.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			Comment: String("new comment"),
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s SET COMMENT = 'new comment'`, id.FullyQualifiedName())
	})

	t.Run("set all", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		opts.Set = &PipeSet{
			ErrorIntegration:    String("new_error_integration"),
			PipeExecutionPaused: Bool(true),
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
			Comment: String("new comment"),
		}
		assertSqlEquals(t, opts, `ALTER PIPE IF EXISTS %s SET ERROR_INTEGRATION = new_error_integration, PIPE_EXECUTION_PAUSED = true, TAG "tag_name1" = 'v1', "tag_name2" = 'v2', COMMENT = 'new comment'`, id.FullyQualifiedName())
	})

	t.Run("unset pipe execution paused", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{
			PipeExecutionPaused: Bool(true),
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s UNSET PIPE_EXECUTION_PAUSED`, id.FullyQualifiedName())
	})

	t.Run("unset tag: single", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{
			Tag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag_name1"),
			},
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s UNSET TAG "tag_name1"`, id.FullyQualifiedName())
	})

	t.Run("unset tag: single", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{
			Tag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag_name1"),
				NewAccountObjectIdentifier("tag_name2"),
			},
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s UNSET TAG "tag_name1", "tag_name2"`, id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{
			Comment: Bool(true),
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("unset all", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &PipeUnset{
			PipeExecutionPaused: Bool(true),
			Tag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag_name1"),
				NewAccountObjectIdentifier("tag_name2"),
			},
			Comment: Bool(true),
		}
		assertSqlEquals(t, opts, `ALTER PIPE IF EXISTS %s UNSET PIPE_EXECUTION_PAUSED, TAG "tag_name1", "tag_name2", COMMENT`, id.FullyQualifiedName())
	})

	t.Run("refresh", func(t *testing.T) {
		opts := setUpOpts()
		opts.Refresh = &PipeRefresh{}
		assertSqlEquals(t, opts, `ALTER PIPE %s REFRESH`, id.FullyQualifiedName())
	})

	t.Run("refresh with prefix", func(t *testing.T) {
		opts := setUpOpts()
		opts.Refresh = &PipeRefresh{
			Prefix: String("/d1"),
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s REFRESH PREFIX = '/d1'`, id.FullyQualifiedName())
	})

	t.Run("refresh with modify", func(t *testing.T) {
		opts := setUpOpts()
		opts.Refresh = &PipeRefresh{
			ModifiedAfter: String("2018-07-30T13:56:46-07:00"),
		}
		assertSqlEquals(t, opts, `ALTER PIPE %s REFRESH MODIFIED_AFTER = '2018-07-30T13:56:46-07:00'`, id.FullyQualifiedName())
	})

	t.Run("refresh with all", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		opts.Refresh = &PipeRefresh{
			Prefix:        String("/d1"),
			ModifiedAfter: String("2018-07-30T13:56:46-07:00"),
		}
		assertSqlEquals(t, opts, `ALTER PIPE IF EXISTS %s REFRESH PREFIX = '/d1' MODIFIED_AFTER = '2018-07-30T13:56:46-07:00'`, id.FullyQualifiedName())
	})
}

func TestPipesDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	setUpOpts := func() *PipeDropOptions {
		return &PipeDropOptions{
			name: id,
		}
	}

	t.Run("empty options", func(t *testing.T) {
		opts := setUpOpts()
		assertSqlEquals(t, opts, `DROP PIPE %s`, id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		assertSqlEquals(t, opts, `DROP PIPE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestPipesShow(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)
	databaseIdentifier := NewAccountObjectIdentifier(id.DatabaseName())
	schemaIdentifier := NewSchemaIdentifier(id.DatabaseName(), id.SchemaName())

	setUpOpts := func() *PipeShowOptions {
		return &PipeShowOptions{}
	}

	t.Run("empty options", func(t *testing.T) {
		opts := setUpOpts()
		assertSqlEquals(t, opts, `SHOW PIPES`)
	})

	t.Run("with like", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertSqlEquals(t, opts, `SHOW PIPES LIKE '%s'`, id.Name())
	})

	t.Run("in account", func(t *testing.T) {
		opts := setUpOpts()
		opts.In = &In{
			Account: Bool(true),
		}
		assertSqlEquals(t, opts, `SHOW PIPES IN ACCOUNT`)
	})

	t.Run("in database", func(t *testing.T) {
		opts := setUpOpts()
		opts.In = &In{
			Database: databaseIdentifier,
		}
		assertSqlEquals(t, opts, `SHOW PIPES IN DATABASE %s`, databaseIdentifier.FullyQualifiedName())
	})

	t.Run("in schema", func(t *testing.T) {
		opts := setUpOpts()
		opts.In = &In{
			Schema: schemaIdentifier,
		}
		assertSqlEquals(t, opts, `SHOW PIPES IN SCHEMA %s`, schemaIdentifier.FullyQualifiedName())
	})

	t.Run("with like and in account", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Account: Bool(true),
		}
		assertSqlEquals(t, opts, `SHOW PIPES LIKE '%s' IN ACCOUNT`, id.Name())
	})

	t.Run("with like and in database", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Database: databaseIdentifier,
		}
		assertSqlEquals(t, opts, `SHOW PIPES LIKE '%s' IN DATABASE %s`, id.Name(), databaseIdentifier.FullyQualifiedName())
	})

	t.Run("with like and in schema", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Schema: schemaIdentifier,
		}
		assertSqlEquals(t, opts, `SHOW PIPES LIKE '%s' IN SCHEMA %s`, id.Name(), schemaIdentifier.FullyQualifiedName())
	})
}

func TestPipesDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("with name", func(t *testing.T) {
		opts := &describePipeOptions{
			name: id,
		}
		assertSqlEquals(t, opts, `DESCRIBE PIPE %s`, id.FullyQualifiedName())
	})
}
