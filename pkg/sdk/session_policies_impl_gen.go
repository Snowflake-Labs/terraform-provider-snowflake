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

func (v *sessionPolicies) Alter(ctx context.Context, request *AlterSessionPolicyRequest) error {
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

func (r *AlterSessionPolicyRequest) toOpts() *AlterSessionPolicyOptions {
	opts := &AlterSessionPolicyOptions{
		IfExists: r.IfExists,
		name:     r.name,
		RenameTo: r.RenameTo,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &SessionPolicySet{
			SessionIdleTimeoutMins:   r.Set.SessionIdleTimeoutMins,
			SessionUiIdleTimeoutMins: r.Set.SessionUiIdleTimeoutMins,
			Comment:                  r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &SessionPolicyUnset{
			SessionIdleTimeoutMins:   r.Unset.SessionIdleTimeoutMins,
			SessionUiIdleTimeoutMins: r.Unset.SessionUiIdleTimeoutMins,
			Comment:                  r.Unset.Comment,
		}
	}
	return opts
}
