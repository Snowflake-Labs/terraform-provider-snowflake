package sdk

import (
	"testing"

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
