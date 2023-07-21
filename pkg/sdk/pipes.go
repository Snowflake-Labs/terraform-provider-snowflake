package sdk

import (
	"context"
	"time"
)

type Pipes interface {
	// Create creates a pipe.
	Create(ctx context.Context, id SchemaObjectIdentifier, opts *PipeCreateOptions) error
	// Alter modifies an existing pipe.
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *PipeAlterOptions) error
	// Drop removes a pipe.
	Drop(ctx context.Context, id SchemaObjectIdentifier, opts *PipeDropOptions) error
	// Show returns a list of pipes.
	Show(ctx context.Context, opts *PipeShowOptions) ([]*Pipe, error)
	// ShowByID returns a pipe by ID.
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Pipe, error)
	// Describe returns the details of a pipe.
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*Pipe, error)
}

// PipeCreateOptions contains options for creating a new pipe in the system for defining the COPY INTO <table> statement
// used by Snowpipe to load data from an ingestion queue into tables.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/create-pipe.
type PipeCreateOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	pipe        bool                   `ddl:"static" sql:"PIPE"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"` //lint:ignore U1000 This is used in the ddl tag

	AutoIngest       *bool   `ddl:"parameter" sql:"AUTO_INGEST"`
	ErrorIntegration *string `ddl:"parameter,no_quotes" sql:"ERROR_INTEGRATION"`
	AwsSnsTopic      *string `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC"`
	Integration      *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
	Comment          *string `ddl:"parameter,single_quotes" sql:"COMMENT"`

	as            bool   `ddl:"static" sql:"AS"` //lint:ignore U1000 This is used in the ddl tag
	CopyStatement string `ddl:"keyword,no_quotes"`
}

// PipeAlterOptions contains options for modifying a limited set of properties for an existing pipe object.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/alter-pipe.
type PipeAlterOptions struct {
	alter    bool                   `ddl:"static" sql:"ALTER"` //lint:ignore U1000 This is used in the ddl tag
	role     bool                   `ddl:"static" sql:"PIPE"`  //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`

	// One of
	Set     *PipeSet     `ddl:"list,no_parentheses" sql:"SET"`
	Unset   *PipeUnset   `ddl:"list,no_parentheses" sql:"UNSET"`
	Refresh *PipeRefresh `ddl:"keyword" sql:"REFRESH"`
}

type PipeSet struct {
	ErrorIntegration    *string          `ddl:"parameter,no_quotes" sql:"ERROR_INTEGRATION"`
	PipeExecutionPaused *bool            `ddl:"parameter" sql:"PIPE_EXECUTION_PAUSED"`
	Tag                 []TagAssociation `ddl:"keyword" sql:"TAG"`
	Comment             *string          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type PipeUnset struct {
	PipeExecutionPaused *bool              `ddl:"keyword" sql:"PIPE_EXECUTION_PAUSED"`
	Tag                 []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
	Comment             *bool              `ddl:"keyword" sql:"COMMENT"`
}

type PipeRefresh struct {
	Prefix        *string `ddl:"parameter,single_quotes" sql:"PREFIX"`
	ModifiedAfter *string `ddl:"parameter,single_quotes" sql:"MODIFIED_AFTER"`
}

// PipeDropOptions contains options for removing the specified pipe from the current/specified schema.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/drop-pipe.
type PipeDropOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"` //lint:ignore U1000 This is used in the ddl tag
	pipe     bool                   `ddl:"static" sql:"PIPE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// PipeShowOptions contains options for showing pipes which user has access privilege to.
//
// https://docs.snowflake.com/en/sql-reference/sql/show-pipes
type PipeShowOptions struct {
	show  bool  `ddl:"static" sql:"SHOW"`  //lint:ignore U1000 This is used in the ddl tag
	pipes bool  `ddl:"static" sql:"PIPES"` //lint:ignore U1000 This is used in the ddl tag
	Like  *Like `ddl:"keyword" sql:"LIKE"`
	In    *In   `ddl:"keyword" sql:"IN"`
}

// pipeDBRow is used to decode the result of a SHOW PIPES query.
type pipeDBRow struct {
	CreatedOn           time.Time `db:"created_on"`
	Name                string    `db:"name"`
	DatabaseName        string    `db:"database_name"`
	SchemaName          string    `db:"schema_name"`
	Definition          string    `db:"definition"`
	Owner               string    `db:"owner"`
	NotificationChannel string    `db:"notification_channel"`
	Comment             string    `db:"comment"`
	Integration         string    `db:"integration"`
	Pattern             string    `db:"pattern"`
	ErrorIntegration    string    `db:"error_integration"`
	OwnerRoleType       string    `db:"owner_role_type"`
	InvalidReason       string    `db:"invalid_reason"`
}

// Pipe is a user-friendly result for a SHOW PIPES and DESCRIBE PIPE queries.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/show-pipes#output and https://docs.snowflake.com/en/sql-reference/sql/desc-pipe#output.
type Pipe struct {
	CreatedOn           time.Time
	Name                string
	DatabaseName        string
	SchemaName          string
	Definition          string
	Owner               string
	NotificationChannel string
	Comment             string
	Integration         string
	Pattern             string
	ErrorIntegration    string
	OwnerRoleType       string
	InvalidReason       string
}

func (v *Pipe) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Pipe) ObjectType() ObjectType {
	return ObjectTypePipe
}

func (row pipeDBRow) toPipe() *Pipe {
	return &Pipe{
		CreatedOn:           row.CreatedOn,
		Name:                row.Name,
		DatabaseName:        row.DatabaseName,
		SchemaName:          row.SchemaName,
		Definition:          row.Definition,
		Owner:               row.Owner,
		NotificationChannel: row.NotificationChannel,
		Comment:             row.Comment,
		Integration:         row.Integration,
		Pattern:             row.Pattern,
		ErrorIntegration:    row.ErrorIntegration,
		OwnerRoleType:       row.OwnerRoleType,
		InvalidReason:       row.InvalidReason,
	}
}

// describePipeOptions contains options for describing the properties specified for a pipe, as well as the default values of the properties.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/desc-pipe.
type describePipeOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"` //lint:ignore U1000 This is used in the ddl tag
	pipe     bool                   `ddl:"static" sql:"PIPE"`     //lint:ignore U1000 This is used in the ddl tag
	name     SchemaObjectIdentifier `ddl:"identifier"`
}
