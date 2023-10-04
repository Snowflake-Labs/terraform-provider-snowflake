package sdk

import (
	"context"
	"database/sql"
)

type SessionPolicies interface {
	Create(ctx context.Context, request *CreateSessionPolicyRequest) error
	Alter(ctx context.Context, request *AlterSessionPolicyRequest) error
	Drop(ctx context.Context, request *DropSessionPolicyRequest) error
	Show(ctx context.Context, request *ShowSessionPolicyRequest) ([]SessionPolicy, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*SessionPolicy, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*SessionPolicyDescription, error)
}

// CreateSessionPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-session-policy.
type CreateSessionPolicyOptions struct {
	create                   bool                   `ddl:"static" sql:"CREATE"`
	OrReplace                *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	sessionPolicy            bool                   `ddl:"static" sql:"SESSION POLICY"`
	IfNotExists              *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                     SchemaObjectIdentifier `ddl:"identifier"`
	SessionIdleTimeoutMins   *int                   `ddl:"parameter,no_quotes" sql:"SESSION_IDLE_TIMEOUT_MINS"`
	SessionUiIdleTimeoutMins *int                   `ddl:"parameter,no_quotes" sql:"SESSION_UI_IDLE_TIMEOUT_MINS"`
	Comment                  *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterSessionPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-session-policy.
type AlterSessionPolicyOptions struct {
	alter         bool                    `ddl:"static" sql:"ALTER"`
	sessionPolicy bool                    `ddl:"static" sql:"SESSION POLICY"`
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier  `ddl:"identifier"`
	RenameTo      *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set           *SessionPolicySet       `ddl:"keyword" sql:"SET"`
	SetTags       []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags     []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
	Unset         *SessionPolicyUnset     `ddl:"keyword" sql:"UNSET"`
}

type SessionPolicySet struct {
	SessionIdleTimeoutMins   *int    `ddl:"parameter,no_quotes" sql:"SESSION_IDLE_TIMEOUT_MINS"`
	SessionUiIdleTimeoutMins *int    `ddl:"parameter,no_quotes" sql:"SESSION_UI_IDLE_TIMEOUT_MINS"`
	Comment                  *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type SessionPolicyUnset struct {
	SessionIdleTimeoutMins   *bool `ddl:"keyword" sql:"SESSION_IDLE_TIMEOUT_MINS"`
	SessionUiIdleTimeoutMins *bool `ddl:"keyword" sql:"SESSION_UI_IDLE_TIMEOUT_MINS"`
	Comment                  *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropSessionPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-session-policy.
type DropSessionPolicyOptions struct {
	drop          bool                   `ddl:"static" sql:"DROP"`
	sessionPolicy bool                   `ddl:"static" sql:"SESSION POLICY"`
	IfExists      *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowSessionPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-session-policies.
type ShowSessionPolicyOptions struct {
	show            bool `ddl:"static" sql:"SHOW"`
	sessionPolicies bool `ddl:"static" sql:"SESSION POLICIES"`
}

type showSessionPolicyDBRow struct {
	CreatedOn    string `db:"created_on"`
	Name         string `db:"name"`
	DatabaseName string `db:"database_name"`
	SchemaName   string `db:"schema_name"`
	Kind         string `db:"kind"`
	Owner        string `db:"owner"`
	Comment      string `db:"comment"`
	Options      string `db:"options"`
}

type SessionPolicy struct {
	CreatedOn    string
	Name         string
	DatabaseName string
	SchemaName   string
	Kind         string
	Owner        string
	Comment      string
	Options      string
}

// DescribeSessionPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-session-policy.
type DescribeSessionPolicyOptions struct {
	describe      bool                   `ddl:"static" sql:"DESCRIBE"`
	sessionPolicy bool                   `ddl:"static" sql:"SESSION POLICY"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

type describeSessionPolicyDBRow struct {
	CreatedOn                string         `db:"created_on"`
	Name                     string         `db:"name"`
	SessionIdleTimeoutMins   int            `db:"session_idle_timeout_mins"`
	SessionUiIdleTimeoutMins int            `db:"session_ui_idle_timeout_mins"`
	Comment                  sql.NullString `db:"comment"`
}

type SessionPolicyDescription struct {
	CreatedOn                string
	Name                     string
	SessionIdleTimeoutMins   int
	SessionUIIdleTimeoutMins int
	Comment                  string
}

func (v *SessionPolicy) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}
