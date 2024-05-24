package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ RowAccessPolicies = (*rowAccessPolicies)(nil)

type rowAccessPolicies struct {
	client *Client
}

func (v *rowAccessPolicies) Create(ctx context.Context, request *CreateRowAccessPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *rowAccessPolicies) Alter(ctx context.Context, request *AlterRowAccessPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *rowAccessPolicies) Drop(ctx context.Context, request *DropRowAccessPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *rowAccessPolicies) Show(ctx context.Context, request *ShowRowAccessPolicyRequest) ([]RowAccessPolicy, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[rowAccessPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[rowAccessPolicyDBRow, RowAccessPolicy](dbRows)
	return resultList, nil
}

func (v *rowAccessPolicies) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*RowAccessPolicy, error) {
	request := NewShowRowAccessPolicyRequest().WithIn(&In{Schema: id.SchemaId()}).WithLike(&Like{String(id.Name())})
	rowAccessPolicies, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindOne(rowAccessPolicies, func(r RowAccessPolicy) bool { return r.Name == id.Name() })
}

func (v *rowAccessPolicies) Describe(ctx context.Context, id SchemaObjectIdentifier) (*RowAccessPolicyDescription, error) {
	opts := &DescribeRowAccessPolicyOptions{
		name: id,
	}
	result, err := validateAndQueryOne[describeRowAccessPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (r *CreateRowAccessPolicyRequest) toOpts() *CreateRowAccessPolicyOptions {
	opts := &CreateRowAccessPolicyOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		body:    r.body,
		Comment: r.Comment,
	}
	if r.args != nil {
		s := make([]CreateRowAccessPolicyArgs, len(r.args))
		for i, v := range r.args {
			s[i] = CreateRowAccessPolicyArgs(v)
		}
		opts.args = s
	}
	return opts
}

func (r *AlterRowAccessPolicyRequest) toOpts() *AlterRowAccessPolicyOptions {
	opts := &AlterRowAccessPolicyOptions{
		name:         r.name,
		RenameTo:     r.RenameTo,
		SetBody:      r.SetBody,
		SetTags:      r.SetTags,
		UnsetTags:    r.UnsetTags,
		SetComment:   r.SetComment,
		UnsetComment: r.UnsetComment,
	}
	return opts
}

func (r *DropRowAccessPolicyRequest) toOpts() *DropRowAccessPolicyOptions {
	opts := &DropRowAccessPolicyOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowRowAccessPolicyRequest) toOpts() *ShowRowAccessPolicyOptions {
	opts := &ShowRowAccessPolicyOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r rowAccessPolicyDBRow) convert() *RowAccessPolicy {
	rowAccessPolicy := &RowAccessPolicy{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		DatabaseName:  r.DatabaseName,
		SchemaName:    r.SchemaName,
		Kind:          r.Kind,
		Owner:         r.Owner,
		Options:       r.Options,
		OwnerRoleType: r.OwnerRoleType,
	}
	if r.Comment.Valid {
		rowAccessPolicy.Comment = r.Comment.String
	}
	return rowAccessPolicy
}

func (r *DescribeRowAccessPolicyRequest) toOpts() *DescribeRowAccessPolicyOptions {
	opts := &DescribeRowAccessPolicyOptions{
		name: r.name,
	}
	return opts
}

func (r describeRowAccessPolicyDBRow) convert() *RowAccessPolicyDescription {
	rowAccessPolicyDescription := &RowAccessPolicyDescription{
		Name:       r.Name,
		Signature:  r.Signature,
		ReturnType: r.ReturnType,
		Body:       r.Body,
	}
	return rowAccessPolicyDescription
}
