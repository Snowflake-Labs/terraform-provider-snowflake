// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Streams interface {
	CreateOnTable(ctx context.Context, request *CreateOnTableStreamRequest) error
	CreateOnExternalTable(ctx context.Context, request *CreateOnExternalTableStreamRequest) error
	CreateOnDirectoryTable(ctx context.Context, request *CreateOnDirectoryTableStreamRequest) error
	CreateOnView(ctx context.Context, request *CreateOnViewStreamRequest) error
	Clone(ctx context.Context, request *CloneStreamRequest) error
	Alter(ctx context.Context, request *AlterStreamRequest) error
	Drop(ctx context.Context, request *DropStreamRequest) error
	Show(ctx context.Context, request *ShowStreamRequest) ([]Stream, error)
	ShowByID(ctx context.Context, request *ShowByIdStreamRequest) (*Stream, error)
	Describe(ctx context.Context, request *DescribeStreamRequest) (*Stream, error)
}

// CreateOnTableStreamOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stream.
type CreateOnTableStreamOptions struct {
	create          bool                   `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	stream          bool                   `ddl:"static" sql:"STREAM"`
	IfNotExists     *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
	CopyGrants      *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	onTable         bool                   `ddl:"static" sql:"ON TABLE"`
	TableId         SchemaObjectIdentifier `ddl:"identifier"`
	On              *OnStream              `ddl:"keyword"`
	AppendOnly      *bool                  `ddl:"parameter" sql:"APPEND_ONLY"`
	ShowInitialRows *bool                  `ddl:"parameter" sql:"SHOW_INITIAL_ROWS"`
	Comment         *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type OnStream struct {
	At        *bool             `ddl:"keyword" sql:"AT"`
	Before    *bool             `ddl:"keyword" sql:"BEFORE"`
	Statement OnStreamStatement `ddl:"list,parentheses"`
}

type OnStreamStatement struct {
	Timestamp *string `ddl:"parameter,arrow_equals" sql:"TIMESTAMP"`
	Offset    *string `ddl:"parameter,arrow_equals" sql:"OFFSET"`
	Statement *string `ddl:"parameter,arrow_equals" sql:"STATEMENT"`
	Stream    *string `ddl:"parameter,single_quotes,arrow_equals" sql:"STREAM"`
}

// CreateOnExternalTableStreamOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stream.
type CreateOnExternalTableStreamOptions struct {
	create          bool                   `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	stream          bool                   `ddl:"static" sql:"STREAM"`
	IfNotExists     *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
	CopyGrants      *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	onExternalTable bool                   `ddl:"static" sql:"ON EXTERNAL TABLE"`
	ExternalTableId SchemaObjectIdentifier `ddl:"identifier"`
	On              *OnStream              `ddl:"keyword"`
	InsertOnly      *bool                  `ddl:"parameter" sql:"INSERT_ONLY"`
	Comment         *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// CreateOnDirectoryTableStreamOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stream.
type CreateOnDirectoryTableStreamOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	stream      bool                   `ddl:"static" sql:"STREAM"`
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	CopyGrants  *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	onStage     bool                   `ddl:"static" sql:"ON STAGE"`
	StageId     SchemaObjectIdentifier `ddl:"identifier"`
	Comment     *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// CreateOnViewStreamOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stream.
type CreateOnViewStreamOptions struct {
	create          bool                   `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	stream          bool                   `ddl:"static" sql:"STREAM"`
	IfNotExists     *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
	CopyGrants      *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	onView          bool                   `ddl:"static" sql:"ON VIEW"`
	ViewId          SchemaObjectIdentifier `ddl:"identifier"`
	On              *OnStream              `ddl:"keyword"`
	AppendOnly      *bool                  `ddl:"parameter" sql:"APPEND_ONLY"`
	ShowInitialRows *bool                  `ddl:"parameter" sql:"SHOW_INITIAL_ROWS"`
	Comment         *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// CloneStreamOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-stream#variant-syntax.
type CloneStreamOptions struct {
	create       bool                   `ddl:"static" sql:"CREATE"`
	OrReplace    *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	stream       bool                   `ddl:"static" sql:"STREAM"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
	sourceStream SchemaObjectIdentifier `ddl:"identifier" sql:"CLONE"`
	CopyGrants   *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
}

// AlterStreamOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-stream.
type AlterStreamOptions struct {
	alter        bool                   `ddl:"static" sql:"ALTER"`
	stream       bool                   `ddl:"static" sql:"STREAM"`
	IfExists     *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
	SetComment   *string                `ddl:"parameter,single_quotes" sql:"SET COMMENT"`
	UnsetComment *bool                  `ddl:"keyword" sql:"UNSET COMMENT"`
	SetTags      []TagAssociation       `ddl:"keyword" sql:"SET TAG"`
	UnsetTags    []ObjectIdentifier     `ddl:"keyword" sql:"UNSET TAG"`
}

// DropStreamOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-stream.
type DropStreamOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	stream   bool                   `ddl:"static" sql:"STREAM"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowStreamOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-streams.
type ShowStreamOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	Terse      *bool      `ddl:"keyword" sql:"TERSE"`
	streams    bool       `ddl:"static" sql:"STREAMS"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	In         *In        `ddl:"keyword" sql:"IN"`
	StartsWith *string    `ddl:"parameter,no_equals,single_quotes" sql:"STARTS WITH"`
	Limit      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type showStreamsDbRow struct {
	CreatedOn     time.Time      `db:"created_on"`
	Name          string         `db:"name"`
	DatabaseName  string         `db:"database_name"`
	SchemaName    string         `db:"schema_name"`
	TableOn       sql.NullString `db:"tableOn"`
	Owner         sql.NullString `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	TableName     sql.NullString `db:"table_name"`
	SourceType    sql.NullString `db:"source_type"`
	BaseTables    sql.NullString `db:"base_tables"`
	Type          sql.NullString `db:"type"`
	Stale         sql.NullString `db:"stale"`
	Mode          sql.NullString `db:"mode"`
	StaleAfter    sql.NullTime   `db:"stale_after"`
	InvalidReason sql.NullString `db:"invalid_reason"`
	OwnerRoleType sql.NullString `db:"owner_role_type"`
}

type Stream struct {
	CreatedOn     time.Time
	Name          string
	DatabaseName  string
	SchemaName    string
	TableOn       *string
	Owner         *string
	Comment       *string
	TableName     *string
	SourceType    *string
	BaseTables    *string
	Type          *string
	Stale         *string
	Mode          *string
	StaleAfter    *time.Time
	InvalidReason *string
	OwnerRoleType *string
}

// DescribeStreamOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-stream.
type DescribeStreamOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	stream   bool                   `ddl:"static" sql:"STREAM"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}
