package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ ExternalVolumes = (*externalVolumes)(nil)

type externalVolumes struct {
	client *Client
}

func (v *externalVolumes) Create(ctx context.Context, request *CreateExternalVolumeRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalVolumes) Alter(ctx context.Context, request *AlterExternalVolumeRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalVolumes) Drop(ctx context.Context, request *DropExternalVolumeRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalVolumes) Describe(ctx context.Context, id AccountObjectIdentifier) ([]ExternalVolumeProperty, error) {
	opts := &DescribeExternalVolumeOptions{
		name: id,
	}
	rows, err := validateAndQuery[externalVolumeDescRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[externalVolumeDescRow, ExternalVolumeProperty](rows), nil
}

func (v *externalVolumes) Show(ctx context.Context, request *ShowExternalVolumeRequest) ([]ExternalVolume, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[externalVolumeShowRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[externalVolumeShowRow, ExternalVolume](dbRows)
	return resultList, nil
}

func (v *externalVolumes) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ExternalVolume, error) {
	externalVolumes, err := v.Show(ctx, NewShowExternalVolumeRequest().WithLike(Like{
		Pattern: String(id.Name()),
	}))
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(externalVolumes, func(r ExternalVolume) bool { return r.Name == id.Name() })
}

func (r *CreateExternalVolumeRequest) toOpts() *CreateExternalVolumeOptions {
	opts := &CreateExternalVolumeOptions{
		OrReplace:        r.OrReplace,
		IfNotExists:      r.IfNotExists,
		name:             r.name,
		StorageLocations: r.StorageLocations,
		AllowWrites:      r.AllowWrites,
		Comment:          r.Comment,
	}
	return opts
}

func (r *AlterExternalVolumeRequest) toOpts() *AlterExternalVolumeOptions {
	opts := &AlterExternalVolumeOptions{
		IfExists:              r.IfExists,
		name:                  r.name,
		RemoveStorageLocation: r.RemoveStorageLocation,
	}

	if r.Set != nil {
		opts.Set = &AlterExternalVolumeSet{
			AllowWrites: r.Set.AllowWrites,
			Comment:     r.Set.Comment,
		}
	}

	if r.AddStorageLocation != nil {

		opts.AddStorageLocation = &ExternalVolumeStorageLocation{}

		if r.AddStorageLocation.S3StorageLocationParams != nil {

			opts.AddStorageLocation.S3StorageLocationParams = &S3StorageLocationParams{
				Name:                 r.AddStorageLocation.S3StorageLocationParams.Name,
				StorageProvider:      r.AddStorageLocation.S3StorageLocationParams.StorageProvider,
				StorageAwsRoleArn:    r.AddStorageLocation.S3StorageLocationParams.StorageAwsRoleArn,
				StorageBaseUrl:       r.AddStorageLocation.S3StorageLocationParams.StorageBaseUrl,
				StorageAwsExternalId: r.AddStorageLocation.S3StorageLocationParams.StorageAwsExternalId,
			}

			if r.AddStorageLocation.S3StorageLocationParams.Encryption != nil {
				opts.AddStorageLocation.S3StorageLocationParams.Encryption = &ExternalVolumeS3Encryption{
					Type:     r.AddStorageLocation.S3StorageLocationParams.Encryption.Type,
					KmsKeyId: r.AddStorageLocation.S3StorageLocationParams.Encryption.KmsKeyId,
				}
			}

		}

		if r.AddStorageLocation.GCSStorageLocationParams != nil {

			opts.AddStorageLocation.GCSStorageLocationParams = &GCSStorageLocationParams{
				Name:           r.AddStorageLocation.GCSStorageLocationParams.Name,
				StorageBaseUrl: r.AddStorageLocation.GCSStorageLocationParams.StorageBaseUrl,
			}

			if r.AddStorageLocation.GCSStorageLocationParams.Encryption != nil {
				opts.AddStorageLocation.GCSStorageLocationParams.Encryption = &ExternalVolumeGCSEncryption{
					Type:     r.AddStorageLocation.GCSStorageLocationParams.Encryption.Type,
					KmsKeyId: r.AddStorageLocation.GCSStorageLocationParams.Encryption.KmsKeyId,
				}
			}

		}

		if r.AddStorageLocation.AzureStorageLocationParams != nil {
			opts.AddStorageLocation.AzureStorageLocationParams = &AzureStorageLocationParams{
				Name:           r.AddStorageLocation.AzureStorageLocationParams.Name,
				AzureTenantId:  r.AddStorageLocation.AzureStorageLocationParams.AzureTenantId,
				StorageBaseUrl: r.AddStorageLocation.AzureStorageLocationParams.StorageBaseUrl,
			}
		}

	}

	return opts
}

func (r *DropExternalVolumeRequest) toOpts() *DropExternalVolumeOptions {
	opts := &DropExternalVolumeOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *DescribeExternalVolumeRequest) toOpts() *DescribeExternalVolumeOptions {
	opts := &DescribeExternalVolumeOptions{
		name: r.name,
	}
	return opts
}

func (r externalVolumeDescRow) convert() *ExternalVolumeProperty {
	return &ExternalVolumeProperty{
		Parent:  r.ParentProperty,
		Name:    r.Property,
		Type:    r.PropertyType,
		Value:   r.PropertyValue,
		Default: r.PropertyDefault,
	}
}

func (r *ShowExternalVolumeRequest) toOpts() *ShowExternalVolumeOptions {
	opts := &ShowExternalVolumeOptions{
		Like: r.Like,
	}
	return opts
}

func (r externalVolumeShowRow) convert() *ExternalVolume {
	return &ExternalVolume{
		Name:        r.Name,
		AllowWrites: r.AllowWrites,
		Comment:     r.Comment,
	}
}
