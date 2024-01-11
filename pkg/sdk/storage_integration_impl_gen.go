package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ StorageIntegrations = (*storageIntegrations)(nil)

type storageIntegrations struct {
	client *Client
}

func (v *storageIntegrations) Create(ctx context.Context, request *CreateStorageIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *storageIntegrations) Alter(ctx context.Context, request *AlterStorageIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *storageIntegrations) Drop(ctx context.Context, request *DropStorageIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *storageIntegrations) Show(ctx context.Context, request *ShowStorageIntegrationRequest) ([]StorageIntegration, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showStorageIntegrationsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[showStorageIntegrationsDbRow, StorageIntegration](dbRows)
	return resultList, nil
}

func (v *storageIntegrations) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*StorageIntegration, error) {
	storageIntegrations, err := v.Show(ctx, NewShowStorageIntegrationRequest().WithLike(&Like{
		Pattern: String(id.Name()),
	}))
	if err != nil {
		return nil, err
	}
	return collections.FindOne(storageIntegrations, func(r StorageIntegration) bool { return r.Name == id.Name() })
}

func (v *storageIntegrations) Describe(ctx context.Context, id AccountObjectIdentifier) ([]StorageIntegrationProperty, error) {
	opts := &DescribeStorageIntegrationOptions{
		name: id,
	}
	rows, err := validateAndQuery[descStorageIntegrationsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[descStorageIntegrationsDbRow, StorageIntegrationProperty](rows), nil
}

func (r *CreateStorageIntegrationRequest) toOpts() *CreateStorageIntegrationOptions {
	opts := &CreateStorageIntegrationOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		Enabled:                 r.Enabled,
		StorageAllowedLocations: r.StorageAllowedLocations,
		StorageBlockedLocations: r.StorageBlockedLocations,
		Comment:                 r.Comment,
	}
	if r.S3StorageProviderParams != nil {
		opts.S3StorageProviderParams = &S3StorageParams{
			StorageAwsRoleArn:   r.S3StorageProviderParams.StorageAwsRoleArn,
			StorageAwsObjectAcl: r.S3StorageProviderParams.StorageAwsObjectAcl,
		}
	}
	if r.GCSStorageProviderParams != nil {
		opts.GCSStorageProviderParams = &GCSStorageParams{}
	}
	if r.AzureStorageProviderParams != nil {
		opts.AzureStorageProviderParams = &AzureStorageParams{
			AzureTenantId: r.AzureStorageProviderParams.AzureTenantId,
		}
	}
	return opts
}

func (r *AlterStorageIntegrationRequest) toOpts() *AlterStorageIntegrationOptions {
	opts := &AlterStorageIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &StorageIntegrationSet{

			Enabled:                 r.Set.Enabled,
			StorageAllowedLocations: r.Set.StorageAllowedLocations,
			StorageBlockedLocations: r.Set.StorageBlockedLocations,
			Comment:                 r.Set.Comment,
		}
		if r.Set.SetS3Params != nil {
			opts.Set.SetS3Params = &SetS3StorageParams{
				StorageAwsRoleArn:   r.Set.SetS3Params.StorageAwsRoleArn,
				StorageAwsObjectAcl: r.Set.SetS3Params.StorageAwsObjectAcl,
			}
		}
		if r.Set.SetAzureParams != nil {
			opts.Set.SetAzureParams = &SetAzureStorageParams{
				AzureTenantId: r.Set.SetAzureParams.AzureTenantId,
			}
		}
	}
	if r.Unset != nil {
		opts.Unset = &StorageIntegrationUnset{
			Enabled:                 r.Unset.Enabled,
			StorageBlockedLocations: r.Unset.StorageBlockedLocations,
			Comment:                 r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropStorageIntegrationRequest) toOpts() *DropStorageIntegrationOptions {
	opts := &DropStorageIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowStorageIntegrationRequest) toOpts() *ShowStorageIntegrationOptions {
	opts := &ShowStorageIntegrationOptions{
		Like: r.Like,
	}
	return opts
}

func (r showStorageIntegrationsDbRow) convert() *StorageIntegration {
	return &StorageIntegration{
		Name:        r.Name,
		StorageType: r.Type,
		Category:    r.Category,
		Enabled:     r.Enabled,
		Comment:     r.Comment,
		CreatedOn:   r.CreatedOn,
	}
}

func (r *DescribeStorageIntegrationRequest) toOpts() *DescribeStorageIntegrationOptions {
	opts := &DescribeStorageIntegrationOptions{
		name: r.name,
	}
	return opts
}

func (r descStorageIntegrationsDbRow) convert() *StorageIntegrationProperty {
	return &StorageIntegrationProperty{
		Name:    r.Property,
		Type:    r.PropertyType,
		Value:   r.PropertyValue,
		Default: r.PropertyDefault,
	}
}
