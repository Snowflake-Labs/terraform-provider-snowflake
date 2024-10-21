package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Connections = (*connections)(nil)

type connections struct {
	client *Client
}

func (v *connections) CreateConnection(ctx context.Context, request *CreateConnectionConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *connections) CreateReplicatedConnection(ctx context.Context, request *CreateReplicatedConnectionConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *connections) AlterConnectionFailover(ctx context.Context, request *AlterConnectionFailoverConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *connections) AlterConnection(ctx context.Context, request *AlterConnectionConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *connections) Drop(ctx context.Context, request *DropConnectionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *connections) Show(ctx context.Context, request *ShowConnectionRequest) ([]Connection, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[connectionRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[connectionRow, Connection](dbRows)
	return resultList, nil
}

func (v *connections) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Connection, error) {
	connections, err := v.Show(ctx, NewShowConnectionRequest().WithLike(Like{String(id.Name())}))
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(connections, func(r Connection) bool { return r.Name == id.Name() })
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

func (r *DropConnectionRequest) toOpts() *DropConnectionOptions {
	opts := &DropConnectionOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowConnectionRequest) toOpts() *ShowConnectionOptions {
	opts := &ShowConnectionOptions{
		Like: r.Like,
	}
	return opts
}

func (r connectionRow) convert() *Connection {
	c := &Connection{
		SnowflakeRegion:           r.SnowflakeRegion,
		CreatedOn:                 r.CreatedOn,
		AccountName:               r.AccountName,
		Name:                      r.Name,
		Primary:                   r.Primary,
		FailoverAllowedToAccounts: ParseCommaSeparatedStringArray(r.FailoverAllowedToAccounts, false),
		ConnectionUrl:             r.ConnectionUrl,
		OrgnizationName:           r.OrgnizationName,
		AccountLocator:            r.AccountLocator,
	}
	b, err := parseBooleanParameter("IsPrimary", r.IsPrimary)
	if err != nil {
		return nil
	}
	c.IsPrimary = *b
	if r.Comment.Valid {
		c.Comment = String(r.Comment.String)
	}

	return c
}
