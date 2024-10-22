package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Storage location structs for testing
var s3StorageLocationParams = &S3StorageLocationParams{
	Name:              "some s3 name",
	StorageProvider:   S3StorageProviderS3,
	StorageAwsRoleArn: "some s3 role arn",
	StorageBaseUrl:    "some s3 base url",
	Encryption: &ExternalVolumeS3Encryption{
		Type:     S3EncryptionTypeSseS3,
		KmsKeyId: String("some s3 kms key id"),
	},
}

var s3StorageLocationParamsNoneEncryption = &S3StorageLocationParams{
	Name:              "some s3 name",
	StorageProvider:   S3StorageProviderS3,
	StorageAwsRoleArn: "some s3 role arn",
	StorageBaseUrl:    "some s3 base url",
	Encryption: &ExternalVolumeS3Encryption{
		Type: S3EncryptionNone,
	},
}

var s3StorageLocationParamsNoEncryption = &S3StorageLocationParams{
	Name:              "some s3 name",
	StorageProvider:   S3StorageProviderS3,
	StorageAwsRoleArn: "some s3 role arn",
	StorageBaseUrl:    "some s3 base url",
}

var s3StorageLocationParamsGov = &S3StorageLocationParams{
	Name:              "some s3 name",
	StorageProvider:   S3StorageProviderS3GOV,
	StorageAwsRoleArn: "some s3 role arn",
	StorageBaseUrl:    "some s3 base url",
	Encryption: &ExternalVolumeS3Encryption{
		Type:     S3EncryptionTypeSseS3,
		KmsKeyId: String("some s3 kms key id"),
	},
}

var s3StorageLocationParamsWithExternalId = &S3StorageLocationParams{
	Name:                 "some s3 name",
	StorageProvider:      S3StorageProviderS3,
	StorageAwsRoleArn:    "some s3 role arn",
	StorageBaseUrl:       "some s3 base url",
	StorageAwsExternalId: String("some s3 external id"),
	Encryption: &ExternalVolumeS3Encryption{
		Type:     S3EncryptionTypeSseS3,
		KmsKeyId: String("some s3 kms key id"),
	},
}

var gcsStorageLocationParams = &GCSStorageLocationParams{
	Name:           "some gcs name",
	StorageBaseUrl: "some gcs base url",
	Encryption: &ExternalVolumeGCSEncryption{
		Type:     GCSEncryptionTypeSseKms,
		KmsKeyId: String("some gcs kms key id"),
	},
}

var gcsStorageLocationParamsNoneEncryption = &GCSStorageLocationParams{
	Name:           "some gcs name",
	StorageBaseUrl: "some gcs base url",
	Encryption: &ExternalVolumeGCSEncryption{
		Type: GCSEncryptionTypeNone,
	},
}

var gcsStorageLocationParamsNoEncryption = &GCSStorageLocationParams{
	Name:           "some gcs name",
	StorageBaseUrl: "some gcs base url",
}

var azureStorageLocationParams = &AzureStorageLocationParams{
	Name:           "some azure name",
	AzureTenantId:  "some azure tenant id",
	StorageBaseUrl: "some azure base url",
}

func TestExternalVolumes_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid CreateExternalVolumeOptions
	defaultOpts := func() *CreateExternalVolumeOptions {
		return &CreateExternalVolumeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateExternalVolumeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateExternalVolumeOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.StorageLocations[i].S3StorageLocationParams opts.StorageLocations[i].GCSStorageLocationParams opts.StorageLocations[i].AzureStorageLocationParams] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{
			{S3StorageLocationParams: s3StorageLocationParams},
			{},
			{S3StorageLocationParams: s3StorageLocationParams, GCSStorageLocationParams: gcsStorageLocationParams},
			{S3StorageLocationParams: s3StorageLocationParams, AzureStorageLocationParams: azureStorageLocationParams},
			{GCSStorageLocationParams: gcsStorageLocationParams, AzureStorageLocationParams: azureStorageLocationParams},
			{
				S3StorageLocationParams:    s3StorageLocationParams,
				GCSStorageLocationParams:   gcsStorageLocationParams,
				AzureStorageLocationParams: azureStorageLocationParams,
			},
		}
		assertOptsInvalidJoinedErrors(
			t,
			opts,
			errExactlyOneOf("CreateExternalVolumeOptions.StorageLocation[1]", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"),
			errExactlyOneOf("CreateExternalVolumeOptions.StorageLocation[2]", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"),
			errExactlyOneOf("CreateExternalVolumeOptions.StorageLocation[3]", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"),
			errExactlyOneOf("CreateExternalVolumeOptions.StorageLocation[4]", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"),
			errExactlyOneOf("CreateExternalVolumeOptions.StorageLocation[5]", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"),
		)
	})

	t.Run("validation: length of opts.StorageLocations is > 0", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateExternalVolumeOptions", "StorageLocations"))
	})

	t.Run("1 storage location - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{S3StorageLocationParams: s3StorageLocationParams}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'some s3 name' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'some s3 role arn' STORAGE_BASE_URL = 'some s3 base url' ENCRYPTION = (TYPE = 'AWS_SSE_S3' KMS_KEY_ID = 'some s3 kms key id')))`, id.FullyQualifiedName())
	})

	t.Run("1 storage location with comment - s3gov", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{S3StorageLocationParams: s3StorageLocationParamsGov}}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'some s3 name' STORAGE_PROVIDER = 'S3GOV' STORAGE_AWS_ROLE_ARN = 'some s3 role arn' STORAGE_BASE_URL = 'some s3 base url' ENCRYPTION = (TYPE = 'AWS_SSE_S3' KMS_KEY_ID = 'some s3 kms key id'))) COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("1 storage location - s3 none encryption", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{S3StorageLocationParams: s3StorageLocationParamsNoneEncryption}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'some s3 name' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'some s3 role arn' STORAGE_BASE_URL = 'some s3 base url' ENCRYPTION = (TYPE = 'NONE')))`, id.FullyQualifiedName())
	})

	t.Run("1 storage location - s3 no encryption", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{S3StorageLocationParams: s3StorageLocationParamsNoEncryption}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'some s3 name' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'some s3 role arn' STORAGE_BASE_URL = 'some s3 base url'))`, id.FullyQualifiedName())
	})

	t.Run("1 storage location with allow writes - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{GCSStorageLocationParams: gcsStorageLocationParams}}
		opts.AllowWrites = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'some gcs name' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'some gcs base url' ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = 'some gcs kms key id'))) ALLOW_WRITES = true`, id.FullyQualifiedName())
	})

	t.Run("1 storage location with allow writes - gcs none encryption", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{GCSStorageLocationParams: gcsStorageLocationParamsNoneEncryption}}
		opts.AllowWrites = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'some gcs name' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'some gcs base url' ENCRYPTION = (TYPE = 'NONE'))) ALLOW_WRITES = true`, id.FullyQualifiedName())
	})

	t.Run("1 storage location with allow writes - gcs no encryption", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{GCSStorageLocationParams: gcsStorageLocationParamsNoEncryption}}
		opts.AllowWrites = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'some gcs name' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'some gcs base url')) ALLOW_WRITES = true`, id.FullyQualifiedName())
	})

	t.Run("1 storage location - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{AzureStorageLocationParams: azureStorageLocationParams}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'some azure name' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'some azure tenant id' STORAGE_BASE_URL = 'some azure base url'))`, id.FullyQualifiedName())
	})

	t.Run("3 storage locations and all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.StorageLocations = []ExternalVolumeStorageLocation{
			{S3StorageLocationParams: s3StorageLocationParamsWithExternalId},
			{GCSStorageLocationParams: gcsStorageLocationParams},
			{AzureStorageLocationParams: azureStorageLocationParams},
		}
		opts.AllowWrites = Bool(false)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'some s3 name' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'some s3 role arn' STORAGE_BASE_URL = 'some s3 base url' STORAGE_AWS_EXTERNAL_ID = 'some s3 external id' ENCRYPTION = (TYPE = 'AWS_SSE_S3' KMS_KEY_ID = 'some s3 kms key id')), (NAME = 'some gcs name' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'some gcs base url' ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = 'some gcs kms key id')), (NAME = 'some azure name' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'some azure tenant id' STORAGE_BASE_URL = 'some azure base url')) ALLOW_WRITES = false COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestExternalVolumes_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid AlterExternalVolumeOptions
	defaultOpts := func() *AlterExternalVolumeOptions {
		return &AlterExternalVolumeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterExternalVolumeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: exactly one field from [opts.RemoveStorageLocation opts.Set opts.AddStorageLocation] should be present - zero set", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
	})

	t.Run("validation: exactly one field from [opts.RemoveStorageLocation opts.Set opts.AddStorageLocation] should be present - two set", func(t *testing.T) {
		removeAndSetOpts := defaultOpts()
		removeAndAddOpts := defaultOpts()
		setAndAddOpts := defaultOpts()

		removeAndSetOpts.RemoveStorageLocation = String("some storage location")
		removeAndSetOpts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true)}

		removeAndAddOpts.RemoveStorageLocation = String("some storage location")
		removeAndAddOpts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsWithExternalId}

		setAndAddOpts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true)}
		setAndAddOpts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsWithExternalId}

		assertOptsInvalidJoinedErrors(t, removeAndSetOpts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
		assertOptsInvalidJoinedErrors(t, removeAndAddOpts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
		assertOptsInvalidJoinedErrors(t, setAndAddOpts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
	})

	t.Run("validation: exactly one field from [opts.RemoveStorageLocation opts.Set opts.AddStorageLocation] should be present - three set", func(t *testing.T) {
		opts := defaultOpts()
		opts.RemoveStorageLocation = String("some storage location")
		opts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true)}
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsWithExternalId}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.AddStorageLocation.S3StorageLocationParams opts.AddStorageLocation.GCSStorageLocationParams opts.AddStorageLocation.AzureStorageLocationParams] should be present - none set", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"))
	})

	t.Run("validation: exactly one field from [opts.AddStorageLocation.S3StorageLocationParams opts.AddStorageLocation.GCSStorageLocationParams opts.AddStorageLocation.AzureStorageLocationParams] should be present - two set", func(t *testing.T) {
		s3AndGcsOpts := defaultOpts()
		s3AndAzureOpts := defaultOpts()
		gcsAndAzureOpts := defaultOpts()
		s3AndGcsOpts.AddStorageLocation = &ExternalVolumeStorageLocation{
			S3StorageLocationParams:  s3StorageLocationParams,
			GCSStorageLocationParams: gcsStorageLocationParams,
		}
		s3AndAzureOpts.AddStorageLocation = &ExternalVolumeStorageLocation{
			S3StorageLocationParams:    s3StorageLocationParams,
			AzureStorageLocationParams: azureStorageLocationParams,
		}
		gcsAndAzureOpts.AddStorageLocation = &ExternalVolumeStorageLocation{
			GCSStorageLocationParams:   gcsStorageLocationParams,
			AzureStorageLocationParams: azureStorageLocationParams,
		}
		assertOptsInvalidJoinedErrors(t, s3AndGcsOpts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"))
		assertOptsInvalidJoinedErrors(t, s3AndAzureOpts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"))
		assertOptsInvalidJoinedErrors(t, gcsAndAzureOpts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"))
	})

	t.Run("validation: exactly one field from [opts.AddStorageLocation.S3StorageLocationParams opts.AddStorageLocation.GCSStorageLocationParams opts.AddStorageLocation.AzureStorageLocationParams] should be present - three set", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{
			S3StorageLocationParams:    s3StorageLocationParams,
			GCSStorageLocationParams:   gcsStorageLocationParams,
			AzureStorageLocationParams: azureStorageLocationParams,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"))
	})

	t.Run("remove storage location", func(t *testing.T) {
		opts := defaultOpts()
		opts.RemoveStorageLocation = String("some storage location")
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s REMOVE STORAGE_LOCATION 'some storage location'`, id.FullyQualifiedName())
	})

	t.Run("set - allow writes", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s SET ALLOW_WRITES = true`, id.FullyQualifiedName())
	})

	t.Run("set - comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AlterExternalVolumeSet{Comment: String("some comment")}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s SET COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("add storage location - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsWithExternalId}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'some s3 name' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'some s3 role arn' STORAGE_BASE_URL = 'some s3 base url' STORAGE_AWS_EXTERNAL_ID = 'some s3 external id' ENCRYPTION = (TYPE = 'AWS_SSE_S3' KMS_KEY_ID = 'some s3 kms key id'))`, id.FullyQualifiedName())
	})

	t.Run("add storage location - s3 none encryption", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsNoneEncryption}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'some s3 name' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'some s3 role arn' STORAGE_BASE_URL = 'some s3 base url' ENCRYPTION = (TYPE = 'NONE'))`, id.FullyQualifiedName())
	})

	t.Run("add storage location - s3 no encryption", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsNoEncryption}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'some s3 name' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'some s3 role arn' STORAGE_BASE_URL = 'some s3 base url')`, id.FullyQualifiedName())
	})

	t.Run("add storage location - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{GCSStorageLocationParams: gcsStorageLocationParams}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'some gcs name' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'some gcs base url' ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = 'some gcs kms key id'))`, id.FullyQualifiedName())
	})

	t.Run("add storage location - gcs none encryption", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{GCSStorageLocationParams: gcsStorageLocationParamsNoneEncryption}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'some gcs name' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'some gcs base url' ENCRYPTION = (TYPE = 'NONE'))`, id.FullyQualifiedName())
	})

	t.Run("add storage location - gcs no encryption", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{GCSStorageLocationParams: gcsStorageLocationParamsNoEncryption}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'some gcs name' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'some gcs base url')`, id.FullyQualifiedName())
	})

	t.Run("add storage location - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{AzureStorageLocationParams: azureStorageLocationParams}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'some azure name' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'some azure tenant id' STORAGE_BASE_URL = 'some azure base url')`, id.FullyQualifiedName())
	})
}

func TestExternalVolumes_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid DropExternalVolumeOptions
	defaultOpts := func() *DropExternalVolumeOptions {
		return &DropExternalVolumeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropExternalVolumeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP EXTERNAL VOLUME %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP EXTERNAL VOLUME IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestExternalVolumes_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid DescribeExternalVolumeOptions
	defaultOpts := func() *DescribeExternalVolumeOptions {
		return &DescribeExternalVolumeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeExternalVolumeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE EXTERNAL VOLUME %s`, id.FullyQualifiedName())
	})
}

func TestExternalVolumes_Show(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid ShowExternalVolumeOptions
	defaultOpts := func() *ShowExternalVolumeOptions {
		return &ShowExternalVolumeOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowExternalVolumeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW EXTERNAL VOLUMES")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW EXTERNAL VOLUMES LIKE '%s'", id.Name())
	})
}

func Test_ExternalVolumes_ToS3EncryptionType(t *testing.T) {
	type test struct {
		input string
		want  S3EncryptionType
	}

	valid := []test{
		{input: "aws_sse_s3", want: S3EncryptionTypeSseS3},
		{input: "AWS_SSE_S3", want: S3EncryptionTypeSseS3},
		{input: "AWS_SSE_KMS", want: S3EncryptionTypeSseKms},
		{input: "NONE", want: S3EncryptionNone},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToS3EncryptionType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToS3EncryptionType(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ExternalVolumes_ToStorageProvider(t *testing.T) {
	type test struct {
		input string
		want  StorageProvider
	}

	valid := []test{
		{input: "s3", want: StorageProviderS3},
		{input: "S3", want: StorageProviderS3},
		{input: "s3gov", want: StorageProviderS3GOV},
		{input: "S3GOV", want: StorageProviderS3GOV},
		{input: "gcs", want: StorageProviderGCS},
		{input: "GCS", want: StorageProviderGCS},
		{input: "azure", want: StorageProviderAzure},
		{input: "AZURE", want: StorageProviderAzure},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToStorageProvider(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToStorageProvider(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ExternalVolumes_ToS3StorageProvider(t *testing.T) {
	type test struct {
		input string
		want  S3StorageProvider
	}

	valid := []test{
		{input: "s3", want: S3StorageProviderS3},
		{input: "S3", want: S3StorageProviderS3},
		{input: "s3gov", want: S3StorageProviderS3GOV},
		{input: "S3GOV", want: S3StorageProviderS3GOV},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToS3StorageProvider(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToS3StorageProvider(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ExternalVolumes_ToGCSEncryptionType(t *testing.T) {
	type test struct {
		input string
		want  GCSEncryptionType
	}

	valid := []test{
		{input: "gcs_sse_kms", want: GCSEncryptionTypeSseKms},
		{input: "GCS_SSE_KMS", want: GCSEncryptionTypeSseKms},
		{input: "NONE", want: GCSEncryptionTypeNone},
		{input: "none", want: GCSEncryptionTypeNone},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToGCSEncryptionType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToGCSEncryptionType(tc.input)
			require.Error(t, err)
		})
	}
}

// External volume helper tests

func Test_GetStorageLocationName(t *testing.T) {
	s3StorageLocationName := "s3Test"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &ExternalVolumeS3Encryption{
			Type:     S3EncryptionTypeSseKms,
			KmsKeyId: &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			Type:     GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	s3GovStorageLocationA := S3StorageLocationParams{
		Name:              s3StorageLocationName,
		StorageProvider:   S3StorageProviderS3GOV,
		StorageBaseUrl:    s3StorageBaseUrl,
		StorageAwsRoleArn: s3StorageAwsRoleArn,
	}

	testCases := []struct {
		Name            string
		StorageLocation ExternalVolumeStorageLocation
		ExpectedName    string
	}{
		{
			Name:            "S3 storage location name succesfully read",
			StorageLocation: ExternalVolumeStorageLocation{S3StorageLocationParams: &s3StorageLocationA},
			ExpectedName:    s3StorageLocationA.Name,
		},
		{
			Name:            "S3GOV storage location name succesfully read",
			StorageLocation: ExternalVolumeStorageLocation{S3StorageLocationParams: &s3GovStorageLocationA},
			ExpectedName:    s3GovStorageLocationA.Name,
		},
		{
			Name:            "GCS storage location name succesfully read",
			StorageLocation: ExternalVolumeStorageLocation{GCSStorageLocationParams: &gcsStorageLocationA},
			ExpectedName:    gcsStorageLocationA.Name,
		},
		{
			Name:            "Azure storage location name succesfully read",
			StorageLocation: ExternalVolumeStorageLocation{AzureStorageLocationParams: &azureStorageLocationA},
			ExpectedName:    azureStorageLocationA.Name,
		},
	}

	invalidTestCases := []struct {
		Name            string
		StorageLocation ExternalVolumeStorageLocation
	}{
		{
			Name:            "Empty S3 storage location",
			StorageLocation: ExternalVolumeStorageLocation{S3StorageLocationParams: &S3StorageLocationParams{}},
		},
		{
			Name:            "Empty GCS storage location",
			StorageLocation: ExternalVolumeStorageLocation{GCSStorageLocationParams: &GCSStorageLocationParams{}},
		},
		{
			Name:            "Empty Azure storage location",
			StorageLocation: ExternalVolumeStorageLocation{AzureStorageLocationParams: &AzureStorageLocationParams{}},
		},
		{
			Name:            "Empty storage location",
			StorageLocation: ExternalVolumeStorageLocation{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			name, err := GetStorageLocationName(tc.StorageLocation)
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedName, name)
		})
	}
	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := GetStorageLocationName(tc.StorageLocation)
			require.Error(t, err)
		})
	}
}

func Test_GetStorageLocationStorageProvider(t *testing.T) {
	s3StorageLocationName := "s3Test"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &ExternalVolumeS3Encryption{
			Type:     S3EncryptionTypeSseKms,
			KmsKeyId: &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			Type:     GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	s3GovStorageLocationA := S3StorageLocationParams{
		Name:              s3StorageLocationName,
		StorageProvider:   S3StorageProviderS3GOV,
		StorageBaseUrl:    s3StorageBaseUrl,
		StorageAwsRoleArn: s3StorageAwsRoleArn,
	}
	testCases := []struct {
		Name                    string
		StorageLocation         ExternalVolumeStorageLocation
		ExpectedStorageProvider StorageProvider
	}{
		{
			Name:                    "S3 storage provider",
			StorageLocation:         ExternalVolumeStorageLocation{S3StorageLocationParams: &s3StorageLocationA},
			ExpectedStorageProvider: StorageProviderS3,
		},
		{
			Name:                    "S3GOV storage provider",
			StorageLocation:         ExternalVolumeStorageLocation{S3StorageLocationParams: &s3GovStorageLocationA},
			ExpectedStorageProvider: StorageProviderS3GOV,
		},
		{
			Name:                    "GCS storage provider",
			StorageLocation:         ExternalVolumeStorageLocation{GCSStorageLocationParams: &gcsStorageLocationA},
			ExpectedStorageProvider: StorageProviderGCS,
		},
		{
			Name:                    "Azure storage provider",
			StorageLocation:         ExternalVolumeStorageLocation{AzureStorageLocationParams: &azureStorageLocationA},
			ExpectedStorageProvider: StorageProviderAzure,
		},
	}

	invalidTestCases := []struct {
		Name            string
		StorageLocation ExternalVolumeStorageLocation
	}{
		{
			Name:            "Empty S3 storage location",
			StorageLocation: ExternalVolumeStorageLocation{S3StorageLocationParams: &S3StorageLocationParams{}},
		},
		{
			Name:            "Empty GCS storage location",
			StorageLocation: ExternalVolumeStorageLocation{GCSStorageLocationParams: &GCSStorageLocationParams{}},
		},
		{
			Name:            "Empty Azure storage location",
			StorageLocation: ExternalVolumeStorageLocation{AzureStorageLocationParams: &AzureStorageLocationParams{}},
		},
		{
			Name:            "Empty storage location",
			StorageLocation: ExternalVolumeStorageLocation{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			storageProvider, err := GetStorageLocationStorageProvider(tc.StorageLocation)
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedStorageProvider, storageProvider)
		})
	}
	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := GetStorageLocationName(tc.StorageLocation)
			require.Error(t, err)
		})
	}
}

var s3StorageAwsExternalId = "1234567890"

func Test_CopySentinelStorageLocation(t *testing.T) {
	tempStorageLocationName := "terraform_provider_sentinel_storage_location"
	s3StorageLocationName := "s3Test"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &ExternalVolumeS3Encryption{
			Type:     S3EncryptionTypeSseKms,
			KmsKeyId: &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			Type:     GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	t.Run("S3 storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{S3StorageLocationParams: &s3StorageLocationA}
		copiedStorageLocation, err := CopySentinelStorageLocation(storageLocationInput)
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageProvider, s3StorageLocationA.StorageProvider)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageBaseUrl, s3StorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageAwsRoleArn, s3StorageLocationA.StorageAwsRoleArn)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageAwsExternalId, s3StorageLocationA.StorageAwsExternalId)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.Encryption.Type, s3StorageLocationA.Encryption.Type)
		assert.Equal(t, *copiedStorageLocation.S3StorageLocationParams.Encryption.KmsKeyId, *s3StorageLocationA.Encryption.KmsKeyId)
	})

	t.Run("GCS storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{GCSStorageLocationParams: &gcsStorageLocationA}
		copiedStorageLocation, err := CopySentinelStorageLocation(storageLocationInput)
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.StorageBaseUrl, gcsStorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.Encryption.Type, gcsStorageLocationA.Encryption.Type)
		assert.Equal(t, *copiedStorageLocation.GCSStorageLocationParams.Encryption.KmsKeyId, *gcsStorageLocationA.Encryption.KmsKeyId)
	})

	t.Run("Azure storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{AzureStorageLocationParams: &azureStorageLocationA}
		copiedStorageLocation, err := CopySentinelStorageLocation(storageLocationInput)
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.StorageBaseUrl, azureStorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.AzureTenantId, azureStorageLocationA.AzureTenantId)
	})

	invalidTestCases := []struct {
		Name            string
		StorageLocation ExternalVolumeStorageLocation
	}{
		{
			Name:            "Empty S3 storage location",
			StorageLocation: ExternalVolumeStorageLocation{S3StorageLocationParams: &S3StorageLocationParams{}},
		},
		{
			Name:            "Empty GCS storage location",
			StorageLocation: ExternalVolumeStorageLocation{GCSStorageLocationParams: &GCSStorageLocationParams{}},
		},
		{
			Name:            "Empty Azure storage location",
			StorageLocation: ExternalVolumeStorageLocation{AzureStorageLocationParams: &AzureStorageLocationParams{}},
		},
		{
			Name:            "Empty storage location",
			StorageLocation: ExternalVolumeStorageLocation{},
		},
	}

	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := CopySentinelStorageLocation(tc.StorageLocation)
			require.Error(t, err)
		})
	}
}

func Test_CommonPrefixLastIndex(t *testing.T) {
	s3StorageLocationName := "s3Test"
	s3StorageLocationName2 := "s3Test2"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageLocationName2 := "gcsTest2"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageLocationName2 := "azureTest2"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &ExternalVolumeS3Encryption{
			Type:     S3EncryptionTypeSseKms,
			KmsKeyId: &s3EncryptionKmsKeyId,
		},
	}

	s3StorageLocationB := S3StorageLocationParams{
		Name:                 s3StorageLocationName2,
		StorageProvider:      S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &ExternalVolumeS3Encryption{
			Type:     S3EncryptionTypeSseKms,
			KmsKeyId: &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	azureStorageLocationB := AzureStorageLocationParams{
		Name:           azureStorageLocationName2,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			Type:     GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	gcsStorageLocationB := GCSStorageLocationParams{
		Name:           gcsStorageLocationName2,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			Type:     GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	gcsStorageLocationC := GCSStorageLocationParams{
		Name:           "test",
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			Type:     GCSEncryptionTypeSseKms,
			KmsKeyId: &gcsEncryptionKmsKeyId,
		},
	}

	s3GovStorageLocationA := S3StorageLocationParams{
		Name:              s3StorageLocationName,
		StorageProvider:   S3StorageProviderS3GOV,
		StorageBaseUrl:    s3StorageBaseUrl,
		StorageAwsRoleArn: s3StorageAwsRoleArn,
	}

	testCases := []struct {
		Name           string
		ListA          []ExternalVolumeStorageLocation
		ListB          []ExternalVolumeStorageLocation
		ExpectedOutput int
	}{
		{
			Name:           "Two empty lists",
			ListA:          []ExternalVolumeStorageLocation{},
			ListB:          []ExternalVolumeStorageLocation{},
			ExpectedOutput: -1,
		},
		{
			Name:           "First list empty",
			ListA:          []ExternalVolumeStorageLocation{},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Second list empty",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{},
			ExpectedOutput: -1,
		},
		{
			Name:           "Lists with no common prefix - length 1",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationB}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Lists with no common prefix - length 2",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}, {AzureStorageLocationParams: &azureStorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationB}, {AzureStorageLocationParams: &azureStorageLocationB}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Identical lists - length 1",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ExpectedOutput: 0,
		},
		{
			Name:           "Identical lists - length 2",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}, {AzureStorageLocationParams: &azureStorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}, {AzureStorageLocationParams: &azureStorageLocationA}},
			ExpectedOutput: 1,
		},
		{
			Name: "Identical lists - length 3",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3GovStorageLocationA},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3GovStorageLocationA},
			},
			ExpectedOutput: 2,
		},
		{
			Name: "Lists with a common prefix - length 3, matching up to and including index 1",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
			},
			ExpectedOutput: 1,
		},
		{
			Name: "Lists with a common prefix - length 4, matching up to and including index 2",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationC},
			},
			ExpectedOutput: 2,
		},
		{
			Name: "Lists with a common prefix - length 4, matching up to and including index 1",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationC},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
				{GCSStorageLocationParams: &gcsStorageLocationC},
			},
			ExpectedOutput: 1,
		},
		{
			Name: "Lists with a common prefix - different lengths, matching up to and including index 1 (last index of shorter list)",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
			},
			ExpectedOutput: 1,
		},
		{
			Name: "Lists with a common prefix - different lengths, matching up to and including index 2",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3StorageLocationB},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
				{AzureStorageLocationParams: &azureStorageLocationB},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3StorageLocationB},
				{GCSStorageLocationParams: &gcsStorageLocationB},
				{AzureStorageLocationParams: &azureStorageLocationB},
			},
			ExpectedOutput: 2,
		},
		{
			Name:           "Empty S3 storage location",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &S3StorageLocationParams{}}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Empty GCS storage location",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{GCSStorageLocationParams: &GCSStorageLocationParams{}}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Empty Azure storage location",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{AzureStorageLocationParams: &AzureStorageLocationParams{}}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Empty storage location",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{}},
			ExpectedOutput: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			commonPrefixLastIndex, err := CommonPrefixLastIndex(tc.ListA, tc.ListB)
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedOutput, commonPrefixLastIndex)
		})
	}
}
