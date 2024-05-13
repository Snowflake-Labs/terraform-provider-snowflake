package sdk

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ SecurityIntegrations = (*securityIntegrations)(nil)

type securityIntegrations struct {
	client *Client
}

func (v *securityIntegrations) CreateSCIM(ctx context.Context, request *CreateSCIMSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterSCIMIntegration(ctx context.Context, request *AlterSCIMIntegrationSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) Drop(ctx context.Context, request *DropSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) Describe(ctx context.Context, id AccountObjectIdentifier) ([]SecurityIntegrationProperty, error) {
	opts := &DescribeSecurityIntegrationOptions{
		name: id,
	}
	rows, err := validateAndQuery[securityIntegrationDescRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	fmt.Println(rows)
	return convertRows[securityIntegrationDescRow, SecurityIntegrationProperty](rows), nil
}

func (v *securityIntegrations) Show(ctx context.Context, request *ShowSecurityIntegrationRequest) ([]SecurityIntegration, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[securityIntegrationShowRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[securityIntegrationShowRow, SecurityIntegration](dbRows)
	return resultList, nil
}

func (v *securityIntegrations) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*SecurityIntegration, error) {
	// TODO: adjust request if e.g. LIKE is supported for the resource
	securityIntegrations, err := v.Show(ctx, NewShowSecurityIntegrationRequest())
	if err != nil {
		return nil, err
	}
	return collections.FindOne(securityIntegrations, func(r SecurityIntegration) bool { return r.Name == id.Name() })
}

func (r *CreateSCIMSecurityIntegrationRequest) toOpts() *CreateSCIMSecurityIntegrationOptions {
	opts := &CreateSCIMSecurityIntegrationOptions{
		OrReplace:     r.OrReplace,
		IfNotExists:   r.IfNotExists,
		name:          r.name,
		Enabled:       r.Enabled,
		ScimClient:    r.ScimClient,
		RunAsRole:     r.RunAsRole,
		NetworkPolicy: r.NetworkPolicy,
		SyncPassword:  r.SyncPassword,
		Comment:       r.Comment,
	}
	return opts
}

func (r *AlterSCIMIntegrationSecurityIntegrationRequest) toOpts() *AlterSCIMIntegrationSecurityIntegrationOptions {
	opts := &AlterSCIMIntegrationSecurityIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,

		SetTag:   r.SetTag,
		UnsetTag: r.UnsetTag,
	}
	if r.Set != nil {
		opts.Set = &SCIMIntegrationSet{
			Enabled:       r.Set.Enabled,
			NetworkPolicy: r.Set.NetworkPolicy,
			SyncPassword:  r.Set.SyncPassword,
			Comment:       r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &SCIMIntegrationUnset{
			NetworkPolicy: r.Unset.NetworkPolicy,
			SyncPassword:  r.Unset.SyncPassword,
			Comment:       r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropSecurityIntegrationRequest) toOpts() *DropSecurityIntegrationOptions {
	opts := &DropSecurityIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *DescribeSecurityIntegrationRequest) toOpts() *DescribeSecurityIntegrationOptions {
	opts := &DescribeSecurityIntegrationOptions{
		name: r.name,
	}
	return opts
}

func (r securityIntegrationDescRow) convert() *SecurityIntegrationProperty {
	return &SecurityIntegrationProperty{
		Name:    r.Property,
		Type:    r.PropertyType,
		Value:   r.PropertyValue,
		Default: r.PropertyDefault,
	}
}

func (r *ShowSecurityIntegrationRequest) toOpts() *ShowSecurityIntegrationOptions {
	opts := &ShowSecurityIntegrationOptions{
		Like: r.Like,
	}
	return opts
}

func (r securityIntegrationShowRow) convert() *SecurityIntegration {
	s := &SecurityIntegration{
		Name:            r.Name,
		IntegrationType: r.Type,
		Enabled:         r.Enabled,
		CreatedOn:       r.CreatedOn,
	}
	if r.Comment.Valid {
		s.Comment = r.Comment.String
	}
	return s
}
