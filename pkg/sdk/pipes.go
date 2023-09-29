package sdk

import (
	"context"
	"database/sql"
)

var _ convertibleRow[Pipe] = new(pipeDBRow)

type Pipes interface {
	// Create creates a pipe.
	Create(ctx context.Context, id SchemaObjectIdentifier, copyStatement string, opts *CreatePipeOptions) error
	// Alter modifies an existing pipe.
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterPipeOptions) error
	// Drop removes a pipe.
	Drop(ctx context.Context, id SchemaObjectIdentifier) error
	// Show returns a list of pipes.
	Show(ctx context.Context, opts *ShowPipeOptions) ([]Pipe, error)
	// ShowByID returns a pipe by ID.
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Pipe, error)
	// Describe returns the details of a pipe.
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*Pipe, error)
}

// CreatePipeOptions contains options for creating a new pipe in the system for defining the COPY INTO <table> statement
// used by Snowpipe to load data from an ingestion queue into tables.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/create-pipe.
type CreatePipeOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	pipe        bool                   `ddl:"static" sql:"PIPE"`
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`

	AutoIngest       *bool   `ddl:"parameter" sql:"AUTO_INGEST"`
	ErrorIntegration *string `ddl:"parameter,no_quotes" sql:"ERROR_INTEGRATION"`
	AwsSnsTopic      *string `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC"`
	Integration      *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
	Comment          *string `ddl:"parameter,single_quotes" sql:"COMMENT"`

	as            bool   `ddl:"static" sql:"AS"`
	copyStatement string `ddl:"keyword,no_quotes"`
}

// AlterPipeOptions contains options for modifying a limited set of properties for an existing pipe object.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/alter-pipe.
type AlterPipeOptions struct {
	alter    bool                   `ddl:"static" sql:"ALTER"`
	role     bool                   `ddl:"static" sql:"PIPE"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`

	// One of
	Set       *PipeSet       `ddl:"list,no_parentheses" sql:"SET"`
	Unset     *PipeUnset     `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTags   *PipeSetTags   `ddl:"list,no_parentheses" sql:"SET TAG"`
	UnsetTags *PipeUnsetTags `ddl:"list,no_parentheses" sql:"UNSET TAG"`
	Refresh   *PipeRefresh   `ddl:"keyword" sql:"REFRESH"`
}

type PipeSet struct {
	ErrorIntegration    *string `ddl:"parameter,no_quotes" sql:"ERROR_INTEGRATION"`
	PipeExecutionPaused *bool   `ddl:"parameter" sql:"PIPE_EXECUTION_PAUSED"`
	Comment             *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type PipeUnset struct {
	PipeExecutionPaused *bool `ddl:"keyword" sql:"PIPE_EXECUTION_PAUSED"`
	Comment             *bool `ddl:"keyword" sql:"COMMENT"`
}

type PipeSetTags struct {
	Tag []TagAssociation `ddl:"keyword"`
}

type PipeUnsetTags struct {
	Tag []ObjectIdentifier `ddl:"keyword"`
}

type PipeRefresh struct {
	Prefix        *string `ddl:"parameter,single_quotes" sql:"PREFIX"`
	ModifiedAfter *string `ddl:"parameter,single_quotes" sql:"MODIFIED_AFTER"`
}

// DropPipeOptions contains options for removing the specified pipe from the current/specified schema.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/drop-pipe.
type DropPipeOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	pipe     bool                   `ddl:"static" sql:"PIPE"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowPipeOptions contains options for showing pipes which user has access privilege to.
//
// https://docs.snowflake.com/en/sql-reference/sql/show-pipes
type ShowPipeOptions struct {
	show  bool  `ddl:"static" sql:"SHOW"`
	pipes bool  `ddl:"static" sql:"PIPES"`
	Like  *Like `ddl:"keyword" sql:"LIKE"`
	In    *In   `ddl:"keyword" sql:"IN"`
}

// pipeDBRow is used to decode the result of a SHOW PIPES query.
type pipeDBRow struct {
	CreatedOn           string         `db:"created_on"`
	Name                string         `db:"name"`
	DatabaseName        string         `db:"database_name"`
	SchemaName          string         `db:"schema_name"`
	Definition          string         `db:"definition"`
	Owner               string         `db:"owner"`
	NotificationChannel sql.NullString `db:"notification_channel"`
	Comment             sql.NullString `db:"comment"`
	Integration         sql.NullString `db:"integration"`
	Pattern             sql.NullString `db:"pattern"`
	ErrorIntegration    sql.NullString `db:"error_integration"`
	OwnerRoleType       sql.NullString `db:"owner_role_type"`
	InvalidReason       sql.NullString `db:"invalid_reason"`
}

// Pipe is a user-friendly result for a SHOW PIPES and DESCRIBE PIPE queries.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/show-pipes#output and https://docs.snowflake.com/en/sql-reference/sql/desc-pipe#output.
type Pipe struct {
	CreatedOn           string
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

func (row pipeDBRow) convert() *Pipe {
	pipe := Pipe{
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		DatabaseName: row.DatabaseName,
		SchemaName:   row.SchemaName,
		Definition:   row.Definition,
		Owner:        row.Owner,
	}
	if row.NotificationChannel.Valid {
		pipe.NotificationChannel = row.NotificationChannel.String
	}
	if row.Comment.Valid {
		pipe.Comment = row.Comment.String
	}
	if row.Integration.Valid {
		pipe.Integration = row.Integration.String
	}
	if row.Pattern.Valid {
		pipe.Pattern = row.Pattern.String
	}
	if row.ErrorIntegration.Valid {
		pipe.ErrorIntegration = row.ErrorIntegration.String
	}
	if row.OwnerRoleType.Valid {
		pipe.OwnerRoleType = row.OwnerRoleType.String
	}
	if row.InvalidReason.Valid {
		pipe.InvalidReason = row.InvalidReason.String
	}
	return &pipe
}

// describePipeOptions contains options for describing the properties specified for a pipe, as well as the default values of the properties.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/desc-pipe.
type describePipeOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	pipe     bool                   `ddl:"static" sql:"PIPE"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}
