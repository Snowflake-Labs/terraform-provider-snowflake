package sdk

import (
	"context"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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
	request := NewShowRowAccessPolicyRequest().WithIn(&ExtendedIn{In: In{Schema: id.SchemaId()}}).WithLike(&Like{String(id.Name())})
	rowAccessPolicies, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(rowAccessPolicies, func(r RowAccessPolicy) bool { return r.Name == id.Name() })
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
		Like:  r.Like,
		In:    r.In,
		Limit: r.Limit,
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
		ReturnType: r.ReturnType,
		Body:       r.Body,
	}
	// Format in database is `(column <data_type>)`
	// TODO(SNOW-1596962): Fully support VECTOR data type
	// TODO(SNOW-1660588): Use ParseFunctionArgumentsFromString
	plainSignature := strings.ReplaceAll(r.Signature, "(", "")
	plainSignature = strings.ReplaceAll(plainSignature, ")", "")
	signatureParts := strings.Split(plainSignature, ", ")
	arguments := make([]RowAccessPolicyArgument, len(signatureParts))

	for i, e := range signatureParts {
		parts := strings.Split(e, " ")
		if len(parts) < 2 {
			log.Printf("[DEBUG] parsing policy arguments: expected argument name and type, got %s", e)
			continue
		}
		dataType, err := ToDataType(parts[len(parts)-1])
		if err != nil {
			log.Printf("[DEBUG] converting row access policy db row: invalid data type %s", dataType)
			continue
		}
		arguments[i] = RowAccessPolicyArgument{
			Name: strings.Join(parts[:len(parts)-1], " "),
			Type: dataType,
		}
	}
	rowAccessPolicyDescription.Signature = arguments

	return rowAccessPolicyDescription
}
