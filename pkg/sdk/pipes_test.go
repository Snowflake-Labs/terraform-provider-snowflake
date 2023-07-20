package sdk

import (
	"testing"
)

func TestPipesCreate(t *testing.T) {
	setUpOpts := func() *PipeCreateOptions {
		return &PipeCreateOptions{
			name:          NewAccountObjectIdentifier("new_pipe"),
			CopyStatement: "<copy_statement>",
		}
	}

	t.Run("basic", func(t *testing.T) {
		opts := setUpOpts()
		assertSqlEquals(t, opts, `CREATE PIPE "new_pipe" AS <copy_statement>`)
	})

	t.Run("if not exists", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfNotExists = Bool(true)
		assertSqlEquals(t, opts, `CREATE PIPE IF NOT EXISTS "new_pipe" AS <copy_statement>`)
	})

	t.Run("auto ingest: true", func(t *testing.T) {
		opts := setUpOpts()
		opts.AutoIngest = Bool(true)
		assertSqlEquals(t, opts, `CREATE PIPE "new_pipe" AUTO_INGEST = true AS <copy_statement>`)
	})

	t.Run("auto ingest: false", func(t *testing.T) {
		opts := setUpOpts()
		opts.AutoIngest = Bool(false)
		assertSqlEquals(t, opts, `CREATE PIPE "new_pipe" AUTO_INGEST = false AS <copy_statement>`)
	})

	t.Run("error integration", func(t *testing.T) {
		opts := setUpOpts()
		opts.ErrorIntegration = String("some error integration")
		assertSqlEquals(t, opts, `CREATE PIPE "new_pipe" ERROR_INTEGRATION = 'some error integration' AS <copy_statement>`)
	})

	t.Run("aws sns topic", func(t *testing.T) {
		opts := setUpOpts()
		opts.AwsSnsTopic = String("some aws sns topic")
		assertSqlEquals(t, opts, `CREATE PIPE "new_pipe" AWS_SNS_TOPIC = 'some aws sns topic' AS <copy_statement>`)
	})

	t.Run("integration", func(t *testing.T) {
		opts := setUpOpts()
		opts.Integration = String("some integration")
		assertSqlEquals(t, opts, `CREATE PIPE "new_pipe" INTEGRATION = 'some integration' AS <copy_statement>`)
	})

	t.Run("comment", func(t *testing.T) {
		opts := setUpOpts()
		opts.Comment = String("some comment")
		assertSqlEquals(t, opts, `CREATE PIPE "new_pipe" COMMENT = 'some comment' AS <copy_statement>`)
	})

	t.Run("all optional", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfNotExists = Bool(true)
		opts.AutoIngest = Bool(true)
		opts.ErrorIntegration = String("some error integration")
		opts.AwsSnsTopic = String("some aws sns topic")
		opts.Integration = String("some integration")
		opts.Comment = String("some comment")
		assertSqlEquals(t, opts, `CREATE PIPE IF NOT EXISTS "new_pipe" AUTO_INGEST = true ERROR_INTEGRATION = 'some error integration' AWS_SNS_TOPIC = 'some aws sns topic' INTEGRATION = 'some integration' COMMENT = 'some comment' AS <copy_statement>`)
	})
}
