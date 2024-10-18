package sdk

import (
	"context"
)

var _ Conntections = (*conntections)(nil)

type conntections struct {
	client *Client
}

func (v *conntections) CreateConnection(ctx context.Context, request *CreateConnectionConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *conntections) CreateReplicatedConnection(ctx context.Context, request *CreateReplicatedConnectionConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *conntections) AlterConnectionFailover(ctx context.Context, request *AlterConnectionFailoverConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *conntections) AlterConnection(ctx context.Context, request *AlterConnectionConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *CreateConnectionConnectionRequest) toOpts() *CreateConnectionConnectionOptions {
	opts := &CreateConnectionConnectionOptions{
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Comment:     r.Comment,
	}
	return opts
}

func (r *CreateReplicatedConnectionConnectionRequest) toOpts() *CreateReplicatedConnectionConnectionOptions {
	opts := &CreateReplicatedConnectionConnectionOptions{
		IfNotExists: r.IfNotExists,
		name:        r.name,
		ReplicaOf:   r.ReplicaOf,
		Comment:     r.Comment,
	}
	return opts
}

func (r *AlterConnectionFailoverConnectionRequest) toOpts() *AlterConnectionFailoverConnectionOptions {
	opts := &AlterConnectionFailoverConnectionOptions{
		name: r.name,
	}

	if r.EnableConnectionFailover != nil {
		opts.EnableConnectionFailover = &EnableConnectionFailover{
			Accounts:           r.EnableConnectionFailover.Accounts,
			IgnoreEditionCheck: r.EnableConnectionFailover.IgnoreEditionCheck,
		}
	}

	if r.DisableConnectionFailover != nil {
		opts.DisableConnectionFailover = &DisableConnectionFailover{
			ToAccounts: r.DisableConnectionFailover.ToAccounts,
			Accounts:   r.DisableConnectionFailover.Accounts,
		}
	}

	if r.Primary != nil {
		opts.Primary = &Primary{}
	}

	return opts
}

func (r *AlterConnectionConnectionRequest) toOpts() *AlterConnectionConnectionOptions {
	opts := &AlterConnectionConnectionOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}

	if r.Set != nil {
		opts.Set = &Set{
			Comment: r.Set.Comment,
		}
	}

	if r.Unset != nil {
		opts.Unset = &Unset{
			Comment: r.Unset.Comment,
		}
	}

	return opts
}
