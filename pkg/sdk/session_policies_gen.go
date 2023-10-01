package sdk

import "context"

type SessionPolicies interface {
	Create(ctx context.Context, request *CreateSessionPolicyRequest) error
	Alter(ctx context.Context, request *AlterSessionPolicyRequest) error
}

// CreateSessionPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-session-policy.
type CreateSessionPolicyOptions struct {
	create                   bool                    `ddl:"static" sql:"CREATE"`
	OrReplace                *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	sessionPolicy            bool                    `ddl:"static" sql:"SESSION POLICY"`
	IfNotExists              *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                     AccountObjectIdentifier `ddl:"identifier"`
	SessionIdleTimeoutMins   *int                    `ddl:"parameter,no_quotes" sql:"SESSION_IDLE_TIMEOUT_MINS"`
	SessionUiIdleTimeoutMins *int                    `ddl:"parameter,no_quotes" sql:"SESSION_UI_IDLE_TIMEOUT_MINS"`
	Comment                  *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterSessionPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-session-policy.
type AlterSessionPolicyOptions struct {
	alter         bool                     `ddl:"static" sql:"ALTER"`
	sessionPolicy bool                     `ddl:"static" sql:"SESSION POLICY"`
	IfExists      *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name          AccountObjectIdentifier  `ddl:"identifier"`
	RenameTo      *AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set           *SessionPolicySet        `ddl:"keyword" sql:"SET"`
	SetTags       []TagAssociation         `ddl:"keyword" sql:"SET TAG"`
	UnsetTags     []ObjectIdentifier       `ddl:"keyword" sql:"UNSET TAG"`
	Unset         *SessionPolicyUnset      `ddl:"keyword" sql:"UNSET"`
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
