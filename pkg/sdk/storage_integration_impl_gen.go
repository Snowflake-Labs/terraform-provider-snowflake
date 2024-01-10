package sdk

import (
	"context"
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
