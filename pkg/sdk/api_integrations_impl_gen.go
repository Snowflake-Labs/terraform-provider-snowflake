package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ ApiIntegrations = (*apiIntegrations)(nil)

type apiIntegrations struct {
	client *Client
}

func (v *apiIntegrations) Create(ctx context.Context, request *CreateApiIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *apiIntegrations) Alter(ctx context.Context, request *AlterApiIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *apiIntegrations) Drop(ctx context.Context, request *DropApiIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *apiIntegrations) Show(ctx context.Context, request *ShowApiIntegrationRequest) ([]ApiIntegration, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showApiIntegrationsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[showApiIntegrationsDbRow, ApiIntegration](dbRows)
	return resultList, nil
}

func (v *apiIntegrations) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegration, error) {
	// TODO: adjust request if e.g. LIKE is supported for the resource
	apiIntegrations, err := v.Show(ctx, NewShowApiIntegrationRequest())
	if err != nil {
		return nil, err
	}
	return collections.FindOne(apiIntegrations, func(r ApiIntegration) bool { return r.Name == id.Name() })
}

func (v *apiIntegrations) Describe(ctx context.Context, id AccountObjectIdentifier) ([]ApiIntegrationProperty, error) {
	opts := &DescribeApiIntegrationOptions{
		name: id,
	}
	rows, err := validateAndQuery[descApiIntegrationsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[descApiIntegrationsDbRow, ApiIntegrationProperty](rows), nil
}

func (r *CreateApiIntegrationRequest) toOpts() *CreateApiIntegrationOptions {
	opts := &CreateApiIntegrationOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		ApiAllowedPrefixes: r.ApiAllowedPrefixes,
		ApiBlockedPrefixes: r.ApiBlockedPrefixes,
		Enabled:            r.Enabled,
		Comment:            r.Comment,
	}
	if r.S3ApiProviderParams != nil {
		opts.S3ApiProviderParams = &S3ApiParams{
			ApiProvider:   r.S3ApiProviderParams.ApiProvider,
			ApiAwsRoleArn: r.S3ApiProviderParams.ApiAwsRoleArn,
			ApiKey:        r.S3ApiProviderParams.ApiKey,
		}
	}
	if r.AzureApiProviderParams != nil {
		opts.AzureApiProviderParams = &AzureApiParams{
			AzureTenantId:        r.AzureApiProviderParams.AzureTenantId,
			AzureAdApplicationId: r.AzureApiProviderParams.AzureAdApplicationId,
			ApiKey:               r.AzureApiProviderParams.ApiKey,
		}
	}
	if r.GCSApiProviderParams != nil {
		opts.GCSApiProviderParams = &GCSApiParams{
			GoogleAudience: r.GCSApiProviderParams.GoogleAudience,
		}
	}
	return opts
}

func (r *AlterApiIntegrationRequest) toOpts() *AlterApiIntegrationOptions {
	opts := &AlterApiIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &ApiIntegrationSet{

			Enabled:            r.Set.Enabled,
			ApiAllowedPrefixes: r.Set.ApiAllowedPrefixes,
			ApiBlockedPrefixes: r.Set.ApiBlockedPrefixes,
			Comment:            r.Set.Comment,
		}
		if r.Set.S3Params != nil {
			opts.Set.S3Params = &SetS3ApiParams{
				ApiAwsRoleArn: r.Set.S3Params.ApiAwsRoleArn,
				ApiKey:        r.Set.S3Params.ApiKey,
			}
		}
		if r.Set.AzureParams != nil {
			opts.Set.AzureParams = &SetAzureApiParams{
				AzureAdApplicationId: r.Set.AzureParams.AzureAdApplicationId,
				ApiKey:               r.Set.AzureParams.ApiKey,
			}
		}
	}
	if r.Unset != nil {
		opts.Unset = &ApiIntegrationUnset{
			ApiKey:             r.Unset.ApiKey,
			Enabled:            r.Unset.Enabled,
			ApiBlockedPrefixes: r.Unset.ApiBlockedPrefixes,
			Comment:            r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropApiIntegrationRequest) toOpts() *DropApiIntegrationOptions {
	opts := &DropApiIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowApiIntegrationRequest) toOpts() *ShowApiIntegrationOptions {
	opts := &ShowApiIntegrationOptions{
		Like: r.Like,
	}
	return opts
}

func (r showApiIntegrationsDbRow) convert() *ApiIntegration {
	// TODO: Mapping
	return &ApiIntegration{}
}

func (r *DescribeApiIntegrationRequest) toOpts() *DescribeApiIntegrationOptions {
	opts := &DescribeApiIntegrationOptions{
		name: r.name,
	}
	return opts
}

func (r descApiIntegrationsDbRow) convert() *ApiIntegrationProperty {
	// TODO: Mapping
	return &ApiIntegrationProperty{}
}
