package sdk

import "context"

type SessionPolicies interface {
	Create(ctx context.Context, request *CreateSessionPolicyRequest) error
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
