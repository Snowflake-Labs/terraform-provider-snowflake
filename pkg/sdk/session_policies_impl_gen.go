package sdk

import "context"

var _ SessionPolicies = (*sessionPolicies)(nil)

type sessionPolicies struct {
	client *Client
}

func (v *sessionPolicies) Create(ctx context.Context, request *CreateSessionPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *CreateSessionPolicyRequest) toOpts() *CreateSessionPolicyOptions {
	opts := &CreateSessionPolicyOptions{
		OrReplace:                r.OrReplace,
		IfNotExists:              r.IfNotExists,
		name:                     r.name,
		SessionIdleTimeoutMins:   r.SessionIdleTimeoutMins,
		SessionUiIdleTimeoutMins: r.SessionUiIdleTimeoutMins,
		Comment:                  r.Comment,
	}
	return opts
}
