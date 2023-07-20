package sdk

import (
	"context"
)

type Pipes interface {
	// Create creates a pipe.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *PipeCreateOptions) error
}

// PipeCreateOptions contains options for creating a pipe.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/create-pipe.
//
// CREATE [ OR REPLACE ] PIPE [ IF NOT EXISTS ] <name>
// [ AUTO_INGEST = [ TRUE | FALSE ] ]
// [ ERROR_INTEGRATION = <integration_name> ]
// [ AWS_SNS_TOPIC = '<string>' ]
// [ INTEGRATION = '<string>' ]
// [ COMMENT = '<string_literal>' ]
// AS <copy_statement>
type PipeCreateOptions struct {
	create      bool                    `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace   *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	pipe        bool                    `ddl:"static" sql:"PIPE"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier `ddl:"identifier"` //lint:ignore U1000 This is used in the ddl tag

	AutoIngest       *bool   `ddl:"parameter" sql:"AUTO_INGEST"`
	ErrorIntegration *string `ddl:"parameter,single_quotes" sql:"ERROR_INTEGRATION"`
	AwsSnsTopic      *string `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC"`
	Integration      *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
	Comment          *string `ddl:"parameter,single_quotes" sql:"COMMENT"`

	as            bool   `ddl:"static" sql:"AS"` //lint:ignore U1000 This is used in the ddl tag
	CopyStatement string `ddl:"keyword,no_quotes"`
}
