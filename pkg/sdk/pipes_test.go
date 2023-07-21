package sdk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *PipeCreateOptions = nil
		assertOptsInvalid(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: copy statement required", func(t *testing.T) {
		opts := setUpOpts()
		opts.CopyStatement = ""
		assertOptsInvalid(t, opts, errCopyStatementRequired)
	})

	t.Run("basic", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsValidAndSqlEquals(t, opts, `CREATE PIPE %s AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("if not exists", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfNotExists = Bool(true)
		assertOptsValidAndSqlEquals(t, opts, `CREATE PIPE IF NOT EXISTS %s AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("auto ingest: true", func(t *testing.T) {
		opts := setUpOpts()
		opts.AutoIngest = Bool(true)
		assertOptsValidAndSqlEquals(t, opts, `CREATE PIPE %s AUTO_INGEST = true AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("auto ingest: false", func(t *testing.T) {
		opts := setUpOpts()
		opts.AutoIngest = Bool(false)
		assertOptsValidAndSqlEquals(t, opts, `CREATE PIPE %s AUTO_INGEST = false AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("error integration", func(t *testing.T) {
		opts := setUpOpts()
		opts.ErrorIntegration = String("some_error_integration")
		assertOptsValidAndSqlEquals(t, opts, `CREATE PIPE %s ERROR_INTEGRATION = some_error_integration AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("aws sns topic", func(t *testing.T) {
		opts := setUpOpts()
		opts.AwsSnsTopic = String("some aws sns topic")
		assertOptsValidAndSqlEquals(t, opts, `CREATE PIPE %s AWS_SNS_TOPIC = 'some aws sns topic' AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("integration", func(t *testing.T) {
		opts := setUpOpts()
		opts.Integration = String("some integration")
		assertOptsValidAndSqlEquals(t, opts, `CREATE PIPE %s INTEGRATION = 'some integration' AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("comment", func(t *testing.T) {
		opts := setUpOpts()
		opts.Comment = String("some comment")
		assertOptsValidAndSqlEquals(t, opts, `CREATE PIPE %s COMMENT = 'some comment' AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("all optional", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfNotExists = Bool(true)
		opts.AutoIngest = Bool(true)
		opts.ErrorIntegration = String("some_error_integration")
		opts.AwsSnsTopic = String("some aws sns topic")
		opts.Integration = String("some integration")
		opts.Comment = String("some comment")
		assertOptsValidAndSqlEquals(t, opts, `CREATE PIPE IF NOT EXISTS %s AUTO_INGEST = true ERROR_INTEGRATION = some_error_integration AWS_SNS_TOPIC = 'some aws sns topic' INTEGRATION = 'some integration' COMMENT = 'some comment' AS <copy_statement>`, id.FullyQualifiedName())
	})
}

func TestPipesAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	setUpOpts := func() *PipeAlterOptions {
		return &PipeAlterOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *PipeAlterOptions = nil
		assertOptsInvalid(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsInvalid(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			ErrorIntegration: String("new_error_integration"),
		}
		opts.Unset = &PipeUnset{
			Comment: Bool(true),
		}
		assertOptsInvalid(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: no property to set", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{}
		assertOptsInvalid(t, opts, errAlterNeedsAtLeastOneProperty)
	})

	t.Run("validation: tags and other property set", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			Tag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag_name1"),
					Value: "v1",
				},
			},
			Comment: String("new comment"),
		}
		assertOptsInvalid(t, opts, errCannotAlterOtherPropertyWithTag)
	})

	t.Run("validation: empty tags slice for set", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			Tag: []TagAssociation{},
		}
		assertOptsInvalid(t, opts, errAlterNeedsAtLeastOneProperty)
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{}
		assertOptsInvalid(t, opts, errAlterNeedsAtLeastOneProperty)
	})

	t.Run("validation: tags and other property unset", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{
			Tag:     []ObjectIdentifier{NewAccountObjectIdentifier("tag_name1")},
			Comment: Bool(true),
		}
		assertOptsInvalid(t, opts, errCannotAlterOtherPropertyWithTag)
	})

	t.Run("validation: empty tags slice for unset", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{
			Tag: []ObjectIdentifier{},
		}
		assertOptsInvalid(t, opts, errAlterNeedsAtLeastOneProperty)
	})

	t.Run("set error integration", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			ErrorIntegration: String("new_error_integration"),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s SET ERROR_INTEGRATION = new_error_integration`, id.FullyQualifiedName())
	})

	t.Run("set pipe execution paused: true", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			PipeExecutionPaused: Bool(true),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s SET PIPE_EXECUTION_PAUSED = true`, id.FullyQualifiedName())
	})

	t.Run("set pipe execution paused: false", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			PipeExecutionPaused: Bool(false),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s SET PIPE_EXECUTION_PAUSED = false`, id.FullyQualifiedName())
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
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s SET TAG "tag_name1" = 'v1'`, id.FullyQualifiedName())
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
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s SET TAG "tag_name1" = 'v1', "tag_name2" = 'v2'`, id.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &PipeSet{
			Comment: String("new comment"),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s SET COMMENT = 'new comment'`, id.FullyQualifiedName())
	})

	t.Run("set more at the same time", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		opts.Set = &PipeSet{
			ErrorIntegration:    String("new_error_integration"),
			PipeExecutionPaused: Bool(true),
			Comment:             String("new comment"),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE IF EXISTS %s SET ERROR_INTEGRATION = new_error_integration, PIPE_EXECUTION_PAUSED = true, COMMENT = 'new comment'`, id.FullyQualifiedName())
	})

	t.Run("unset pipe execution paused", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{
			PipeExecutionPaused: Bool(true),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s UNSET PIPE_EXECUTION_PAUSED`, id.FullyQualifiedName())
	})

	t.Run("unset tag: single", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{
			Tag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag_name1"),
			},
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s UNSET TAG "tag_name1"`, id.FullyQualifiedName())
	})

	t.Run("unset tag: single", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{
			Tag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag_name1"),
				NewAccountObjectIdentifier("tag_name2"),
			},
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s UNSET TAG "tag_name1", "tag_name2"`, id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &PipeUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("unset more at the same time", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &PipeUnset{
			PipeExecutionPaused: Bool(true),
			Comment:             Bool(true),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE IF EXISTS %s UNSET PIPE_EXECUTION_PAUSED, COMMENT`, id.FullyQualifiedName())
	})

	t.Run("refresh", func(t *testing.T) {
		opts := setUpOpts()
		opts.Refresh = &PipeRefresh{}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s REFRESH`, id.FullyQualifiedName())
	})

	t.Run("refresh with prefix", func(t *testing.T) {
		opts := setUpOpts()
		opts.Refresh = &PipeRefresh{
			Prefix: String("/d1"),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s REFRESH PREFIX = '/d1'`, id.FullyQualifiedName())
	})

	t.Run("refresh with modify", func(t *testing.T) {
		opts := setUpOpts()
		opts.Refresh = &PipeRefresh{
			ModifiedAfter: String("2018-07-30T13:56:46-07:00"),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE %s REFRESH MODIFIED_AFTER = '2018-07-30T13:56:46-07:00'`, id.FullyQualifiedName())
	})

	t.Run("refresh with all", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		opts.Refresh = &PipeRefresh{
			Prefix:        String("/d1"),
			ModifiedAfter: String("2018-07-30T13:56:46-07:00"),
		}
		assertOptsValidAndSqlEquals(t, opts, `ALTER PIPE IF EXISTS %s REFRESH PREFIX = '/d1' MODIFIED_AFTER = '2018-07-30T13:56:46-07:00'`, id.FullyQualifiedName())
	})
}

func TestPipesDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	setUpOpts := func() *PipeDropOptions {
		return &PipeDropOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *PipeDropOptions = nil
		assertOptsInvalid(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsValidAndSqlEquals(t, opts, `DROP PIPE %s`, id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSqlEquals(t, opts, `DROP PIPE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestPipesShow(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)
	databaseIdentifier := NewAccountObjectIdentifier(id.DatabaseName())
	schemaIdentifier := NewSchemaIdentifier(id.DatabaseName(), id.SchemaName())

	setUpOpts := func() *PipeShowOptions {
		return &PipeShowOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *PipeShowOptions = nil
		assertOptsInvalid(t, opts, ErrNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{}
		assertOptsInvalid(t, opts, errPatternRequiredForLikeKeyword)
	})

	t.Run("validation: empty in", func(t *testing.T) {
		opts := setUpOpts()
		opts.In = &In{}
		assertOptsInvalid(t, opts, errScopeRequiredForInKeyword)
	})

	t.Run("validation: exactly one scope for in", func(t *testing.T) {
		opts := setUpOpts()
		opts.In = &In{
			Account:  Bool(true),
			Database: databaseIdentifier,
		}
		assertOptsInvalid(t, opts, errScopeRequiredForInKeyword)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsValidAndSqlEquals(t, opts, `SHOW PIPES`)
	})

	t.Run("with like", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSqlEquals(t, opts, `SHOW PIPES LIKE '%s'`, id.Name())
	})

	t.Run("in account", func(t *testing.T) {
		opts := setUpOpts()
		opts.In = &In{
			Account: Bool(true),
		}
		assertOptsValidAndSqlEquals(t, opts, `SHOW PIPES IN ACCOUNT`)
	})

	t.Run("in database", func(t *testing.T) {
		opts := setUpOpts()
		opts.In = &In{
			Database: databaseIdentifier,
		}
		assertOptsValidAndSqlEquals(t, opts, `SHOW PIPES IN DATABASE %s`, databaseIdentifier.FullyQualifiedName())
	})

	t.Run("in schema", func(t *testing.T) {
		opts := setUpOpts()
		opts.In = &In{
			Schema: schemaIdentifier,
		}
		assertOptsValidAndSqlEquals(t, opts, `SHOW PIPES IN SCHEMA %s`, schemaIdentifier.FullyQualifiedName())
	})

	t.Run("with like and in account", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Account: Bool(true),
		}
		assertOptsValidAndSqlEquals(t, opts, `SHOW PIPES LIKE '%s' IN ACCOUNT`, id.Name())
	})

	t.Run("with like and in database", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Database: databaseIdentifier,
		}
		assertOptsValidAndSqlEquals(t, opts, `SHOW PIPES LIKE '%s' IN DATABASE %s`, id.Name(), databaseIdentifier.FullyQualifiedName())
	})

	t.Run("with like and in schema", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Schema: schemaIdentifier,
		}
		assertOptsValidAndSqlEquals(t, opts, `SHOW PIPES LIKE '%s' IN SCHEMA %s`, id.Name(), schemaIdentifier.FullyQualifiedName())
	})
}

func TestPipesDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	setUpOpts := func() *describePipeOptions {
		return &describePipeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *describePipeOptions = nil
		assertOptsInvalid(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("with name", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsValidAndSqlEquals(t, opts, `DESCRIBE PIPE %s`, id.FullyQualifiedName())
	})
}

// assertOptsInvalid could be reused in tests for other interfaces in sdk package.
func assertOptsInvalid(t *testing.T, opts validatableOpts, expectedError error) {
	err := opts.validateProp()
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

// assertOptsValid could be reused in tests for other interfaces in sdk package.
func assertOptsValid(t *testing.T, opts validatableOpts) {
	err := opts.validateProp()
	assert.NoError(t, err)
}

// assertSqlEquals could be reused in tests for other interfaces in sdk package.
func assertSqlEquals(t *testing.T, opts any, format string, args ...any) {
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(format, args...), actual)
}

// assertOptsValidAndSqlEquals could be reused in tests for other interfaces in sdk package.
// It's a shorthand for assertOptsValid and assertSqlEquals.
func assertOptsValidAndSqlEquals(t *testing.T, opts validatableOpts, format string, args ...any) {
	assertOptsValid(t, opts)
	assertSqlEquals(t, opts, format, args...)
}
