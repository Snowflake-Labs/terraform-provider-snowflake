package sdk

import "testing"

func TestStages_CreateInternal(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid CreateInternalStageOptions
	defaultOpts := func() *CreateInternalStageOptions {
		return &CreateInternalStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateInternalStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateInternalStageOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.FileFormat = &StageFileFormat{
			Type: &FileFormatTypeCSV,
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE TEMPORARY STAGE %s FILE_FORMAT = (TYPE = CSV) COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Temporary = Bool(true)
		opts.IfNotExists = Bool(true)
		opts.Encryption = &InternalStageEncryption{
			Type: &InternalStageEncryptionFull,
		}
		opts.DirectoryTableOptions = &InternalDirectoryTableOptions{
			Enable:          Bool(true),
			RefreshOnCreate: Bool(true),
		}
		opts.FileFormat = &StageFileFormat{
			FormatName: String("format name"),
		}
		opts.CopyOptions = &StageCopyOptions{
			OnError: &StageCopyOnErrorOptions{
				Continue: Bool(true),
			},
			SizeLimit:         Int(123),
			Purge:             Bool(true),
			ReturnFailedOnly:  Bool(true),
			MatchByColumnName: &StageCopyColumnMapCaseNone,
			EnforceLength:     Bool(true),
			Truncatecolumns:   Bool(true),
			Force:             Bool(true),
		}
		opts.Comment = String("some comment")
		opts.Tag = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag-name"),
				Value: "tag-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE TEMPORARY STAGE IF NOT EXISTS %s ENCRYPTION = (TYPE = 'SNOWFLAKE_FULL') DIRECTORY = (ENABLE = true REFRESH_ON_CREATE = true) FILE_FORMAT = (FORMAT_NAME = 'format name') COPY_OPTIONS = (ON_ERROR = CONTINUE SIZE_LIMIT = 123 PURGE = true RETURN_FAILED_ONLY = true MATCH_BY_COLUMN_NAME = NONE ENFORCE_LENGTH = true TRUNCATECOLUMNS = true FORCE = true) COMMENT = 'some comment' TAG ("tag-name" = 'tag-value')`, id.FullyQualifiedName())
	})
}

func TestStages_CreateOnS3(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid CreateOnS3StageOptions
	defaultOpts := func() *CreateOnS3StageOptions {
		return &CreateOnS3StageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnS3StageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnS3StageOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.StorageIntegration opts.ExternalStageParams.Credentials]", func(t *testing.T) {
		opts := defaultOpts()
		integrationId := NewAccountObjectIdentifier("integration")
		opts.ExternalStageParams = &ExternalS3StageParams{
			StorageIntegration: &integrationId,
			Credentials: &ExternalStageS3Credentials{
				AWSRole: String("aws-role"),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams", "StorageIntegration", "Credentials"))
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.Credentials.AWSKeyId opts.ExternalStageParams.Credentials.AWSRole]", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalStageParams = &ExternalS3StageParams{
			Credentials: &ExternalStageS3Credentials{
				AWSKeyId: String("aws-key-id"),
				AWSRole:  String("aws-role"),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams.Credentials", "AWSKeyId", "AWSRole"))
	})

	t.Run("all options - storage integration", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		integrationId := NewAccountObjectIdentifier("integration")
		opts.ExternalStageParams = &ExternalS3StageParams{
			Url:                "some url",
			StorageIntegration: &integrationId,
			Encryption: &ExternalStageS3Encryption{
				Type:      &ExternalStageS3EncryptionCSE,
				MasterKey: String("master-key"),
			},
		}
		opts.FileFormat = &StageFileFormat{
			Type: &FileFormatTypeCSV,
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY STAGE %s URL = 'some url' STORAGE_INTEGRATION = "integration" ENCRYPTION = (TYPE = 'AWS_CSE' MASTER_KEY = 'master-key') FILE_FORMAT = (TYPE = CSV) COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("all options - directory table and credentials", func(t *testing.T) {
		opts := defaultOpts()
		opts.Temporary = Bool(true)
		opts.IfNotExists = Bool(true)
		opts.ExternalStageParams = &ExternalS3StageParams{
			Url: "some url",
			Credentials: &ExternalStageS3Credentials{
				AWSKeyId:     String("aws-key-id"),
				AWSSecretKey: String("aws-secret-key"),
				AWSToken:     String("aws-token"),
			},
		}
		opts.DirectoryTableOptions = &ExternalS3DirectoryTableOptions{
			Enable:          Bool(true),
			RefreshOnCreate: Bool(true),
			AutoRefresh:     Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE TEMPORARY STAGE IF NOT EXISTS %s URL = 'some url' CREDENTIALS = (AWS_KEY_ID = 'aws-key-id' AWS_SECRET_KEY = 'aws-secret-key' AWS_TOKEN = 'aws-token') DIRECTORY = (ENABLE = true REFRESH_ON_CREATE = true AUTO_REFRESH = true)`, id.FullyQualifiedName())
	})
}

func TestStages_CreateOnGCS(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid CreateOnGCSStageOptions
	defaultOpts := func() *CreateOnGCSStageOptions {
		return &CreateOnGCSStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnGCSStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnGCSStageOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		integrationId := NewAccountObjectIdentifier("integration")
		opts.ExternalStageParams = &ExternalGCSStageParams{
			Url:                "some url",
			StorageIntegration: &integrationId,
			Encryption: &ExternalStageGCSEncryption{
				Type:     &ExternalStageGCSEncryptionSSEKMS,
				KmsKeyId: String("kms-key-id"),
			},
		}
		opts.DirectoryTableOptions = &ExternalGCSDirectoryTableOptions{
			Enable:                  Bool(true),
			RefreshOnCreate:         Bool(true),
			AutoRefresh:             Bool(true),
			NotificationIntegration: String("notification-integration"),
		}
		opts.FileFormat = &StageFileFormat{
			Type: &FileFormatTypeCSV,
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY STAGE %s URL = 'some url' STORAGE_INTEGRATION = "integration" ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = 'kms-key-id') DIRECTORY = (ENABLE = true REFRESH_ON_CREATE = true AUTO_REFRESH = true NOTIFICATION_INTEGRATION = 'notification-integration') FILE_FORMAT = (TYPE = CSV) COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestStages_CreateOnAzure(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid CreateOnAzureStageOptions
	defaultOpts := func() *CreateOnAzureStageOptions {
		return &CreateOnAzureStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnAzureStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnAzureStageOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.StorageIntegration opts.ExternalStageParams.Credentials]", func(t *testing.T) {
		opts := defaultOpts()
		integrationId := NewAccountObjectIdentifier("integration")
		opts.ExternalStageParams = &ExternalAzureStageParams{
			StorageIntegration: &integrationId,
			Credentials: &ExternalStageAzureCredentials{
				AzureSasToken: "azure-sas-token",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams", "StorageIntegration", "Credentials"))
	})

	t.Run("all options - storage integration", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		integrationId := NewAccountObjectIdentifier("integration")
		opts.ExternalStageParams = &ExternalAzureStageParams{
			Url:                "some url",
			StorageIntegration: &integrationId,
			Encryption: &ExternalStageAzureEncryption{
				Type:      &ExternalStageAzureEncryptionCSE,
				MasterKey: String("master-key"),
			},
		}
		opts.FileFormat = &StageFileFormat{
			Type: &FileFormatTypeCSV,
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY STAGE %s URL = 'some url' STORAGE_INTEGRATION = "integration" ENCRYPTION = (TYPE = 'AZURE_CSE' MASTER_KEY = 'master-key') FILE_FORMAT = (TYPE = CSV) COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("all options - directory table and credentials", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.DirectoryTableOptions = &ExternalAzureDirectoryTableOptions{
			Enable:                  Bool(true),
			RefreshOnCreate:         Bool(true),
			AutoRefresh:             Bool(true),
			NotificationIntegration: String("notification-integration"),
		}
		opts.ExternalStageParams = &ExternalAzureStageParams{
			Url: "some url",
			Credentials: &ExternalStageAzureCredentials{
				AzureSasToken: "azure-sas-token",
			},
			Encryption: &ExternalStageAzureEncryption{
				Type:      &ExternalStageAzureEncryptionCSE,
				MasterKey: String("master-key"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE STAGE IF NOT EXISTS %s URL = 'some url' CREDENTIALS = (AZURE_SAS_TOKEN = 'azure-sas-token') ENCRYPTION = (TYPE = 'AZURE_CSE' MASTER_KEY = 'master-key') DIRECTORY = (ENABLE = true REFRESH_ON_CREATE = true AUTO_REFRESH = true NOTIFICATION_INTEGRATION = 'notification-integration')`, id.FullyQualifiedName())
	})
}

func TestStages_CreateOnS3Compatible(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid CreateOnS3CompatibleStageOptions
	defaultOpts := func() *CreateOnS3CompatibleStageOptions {
		return &CreateOnS3CompatibleStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOnS3CompatibleStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnS3CompatibleStageOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Temporary = Bool(true)
		opts.IfNotExists = Bool(true)
		opts.Url = "some url"
		opts.Endpoint = "some endpoint"
		opts.Credentials = &ExternalStageS3CompatibleCredentials{
			AWSKeyId:     String("aws-key-id"),
			AWSSecretKey: String("aws-secret-key"),
		}
		opts.FileFormat = &StageFileFormat{
			Type: &FileFormatTypeCSV,
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE TEMPORARY STAGE IF NOT EXISTS %s URL = 'some url' ENDPOINT = 'some endpoint' CREDENTIALS = (AWS_KEY_ID = 'aws-key-id' AWS_SECRET_KEY = 'aws-secret-key') FILE_FORMAT = (TYPE = CSV) COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestStages_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterStageOptions
	defaultOpts := func() *AlterStageOptions {
		return &AlterStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.RenameTo] if set", func(t *testing.T) {
		opts := defaultOpts()
		newId := NewSchemaObjectIdentifier("", "", "")
		opts.RenameTo = &newId
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterStageOptions", "RenameTo", "SetTags", "UnsetTags"))
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.UnsetTags]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("id"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterStageOptions", "IfExists", "UnsetTags"))
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("rename", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		newId := RandomSchemaObjectIdentifier()
		opts.RenameTo = &newId
		assertOptsValidAndSQLEquals(t, opts, "ALTER STAGE IF EXISTS %s RENAME TO %s", id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag-name"),
				Value: "tag-value",
			},
			{
				Name:  NewAccountObjectIdentifier("tag-name2"),
				Value: "tag-value2",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER STAGE IF EXISTS %s SET TAG "tag-name" = 'tag-value', "tag-name2" = 'tag-value2'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag-name"),
			NewAccountObjectIdentifier("tag-name2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER STAGE %s UNSET TAG "tag-name", "tag-name2"`, id.FullyQualifiedName())
	})
}

func TestStages_AlterInternalStage(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterInternalStageStageOptions
	defaultOpts := func() *AlterInternalStageStageOptions {
		return &AlterInternalStageStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterInternalStageStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.FileFormat = &StageFileFormat{
			Type: &FileFormatTypeCSV,
		}
		opts.CopyOptions = &StageCopyOptions{
			OnError: &StageCopyOnErrorOptions{
				AbortStatement: Bool(true),
			},
			SizeLimit:         Int(123),
			Purge:             Bool(true),
			ReturnFailedOnly:  Bool(true),
			MatchByColumnName: &StageCopyColumnMapCaseNone,
			EnforceLength:     Bool(true),
			Truncatecolumns:   Bool(true),
			Force:             Bool(true),
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "ALTER STAGE IF EXISTS %s SET FILE_FORMAT = (TYPE = CSV) COPY_OPTIONS = (ON_ERROR = ABORT_STATEMENT SIZE_LIMIT = 123 PURGE = true RETURN_FAILED_ONLY = true MATCH_BY_COLUMN_NAME = NONE ENFORCE_LENGTH = true TRUNCATECOLUMNS = true FORCE = true) COMMENT = 'some comment'", id.FullyQualifiedName())
	})
}

func TestStages_AlterExternalS3Stage(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterExternalS3StageStageOptions
	defaultOpts := func() *AlterExternalS3StageStageOptions {
		return &AlterExternalS3StageStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterExternalS3StageStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.StorageIntegration opts.ExternalStageParams.Credentials]", func(t *testing.T) {
		opts := defaultOpts()
		integrationId := NewAccountObjectIdentifier("integration")
		opts.ExternalStageParams = &ExternalS3StageParams{
			StorageIntegration: &integrationId,
			Credentials: &ExternalStageS3Credentials{
				AWSRole: String("aws-role"),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams", "StorageIntegration", "Credentials"))
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.Credentials.AWSKeyId opts.ExternalStageParams.Credentials.AWSRole]", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalStageParams = &ExternalS3StageParams{
			Credentials: &ExternalStageS3Credentials{
				AWSKeyId: String("aws-key-id"),
				AWSRole:  String("aws-role"),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams.Credentials", "AWSKeyId", "AWSRole"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		integrationId := NewAccountObjectIdentifier("integration")
		opts.ExternalStageParams = &ExternalS3StageParams{
			Url:                "some url",
			StorageIntegration: &integrationId,
			Encryption: &ExternalStageS3Encryption{
				Type: &ExternalStageS3EncryptionNone,
			},
		}
		opts.FileFormat = &StageFileFormat{
			Type: &FileFormatTypeJSON,
		}
		opts.CopyOptions = &StageCopyOptions{
			OnError: &StageCopyOnErrorOptions{
				Continue: Bool(true),
			},
			SizeLimit:         Int(123),
			Purge:             Bool(true),
			ReturnFailedOnly:  Bool(true),
			MatchByColumnName: &StageCopyColumnMapCaseNone,
			EnforceLength:     Bool(true),
			Truncatecolumns:   Bool(true),
			Force:             Bool(true),
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `ALTER STAGE IF EXISTS %s SET URL = 'some url' STORAGE_INTEGRATION = "integration" ENCRYPTION = (TYPE = 'NONE') FILE_FORMAT = (TYPE = JSON) COPY_OPTIONS = (ON_ERROR = CONTINUE SIZE_LIMIT = 123 PURGE = true RETURN_FAILED_ONLY = true MATCH_BY_COLUMN_NAME = NONE ENFORCE_LENGTH = true TRUNCATECOLUMNS = true FORCE = true) COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestStages_AlterExternalGCSStage(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterExternalGCSStageStageOptions
	defaultOpts := func() *AlterExternalGCSStageStageOptions {
		return &AlterExternalGCSStageStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterExternalGCSStageStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		integrationId := NewAccountObjectIdentifier("integration")
		opts.ExternalStageParams = &ExternalGCSStageParams{
			Url:                "some url",
			StorageIntegration: &integrationId,
			Encryption: &ExternalStageGCSEncryption{
				Type: &ExternalStageGCSEncryptionNone,
			},
		}
		opts.FileFormat = &StageFileFormat{
			Type: &FileFormatTypeJSON,
		}
		opts.CopyOptions = &StageCopyOptions{
			OnError: &StageCopyOnErrorOptions{
				Continue: Bool(true),
			},
			SizeLimit:         Int(123),
			Purge:             Bool(true),
			ReturnFailedOnly:  Bool(true),
			MatchByColumnName: &StageCopyColumnMapCaseNone,
			EnforceLength:     Bool(true),
			Truncatecolumns:   Bool(true),
			Force:             Bool(true),
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `ALTER STAGE IF EXISTS %s SET URL = 'some url' STORAGE_INTEGRATION = "integration" ENCRYPTION = (TYPE = 'NONE') FILE_FORMAT = (TYPE = JSON) COPY_OPTIONS = (ON_ERROR = CONTINUE SIZE_LIMIT = 123 PURGE = true RETURN_FAILED_ONLY = true MATCH_BY_COLUMN_NAME = NONE ENFORCE_LENGTH = true TRUNCATECOLUMNS = true FORCE = true) COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestStages_AlterExternalAzureStage(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterExternalAzureStageStageOptions
	defaultOpts := func() *AlterExternalAzureStageStageOptions {
		return &AlterExternalAzureStageStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterExternalAzureStageStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.StorageIntegration opts.ExternalStageParams.Credentials]", func(t *testing.T) {
		opts := defaultOpts()
		integrationId := NewAccountObjectIdentifier("integrationId")
		opts.ExternalStageParams = &ExternalAzureStageParams{
			StorageIntegration: &integrationId,
			Credentials:        &ExternalStageAzureCredentials{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams", "StorageIntegration", "Credentials"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		integrationId := NewAccountObjectIdentifier("integration")
		opts.ExternalStageParams = &ExternalAzureStageParams{
			Url:                "some url",
			StorageIntegration: &integrationId,
			Encryption: &ExternalStageAzureEncryption{
				Type: &ExternalStageAzureEncryptionNone,
			},
		}
		opts.FileFormat = &StageFileFormat{
			Type: &FileFormatTypeJSON,
		}
		opts.CopyOptions = &StageCopyOptions{
			OnError: &StageCopyOnErrorOptions{
				Continue: Bool(true),
			},
			SizeLimit:         Int(123),
			Purge:             Bool(true),
			ReturnFailedOnly:  Bool(true),
			MatchByColumnName: &StageCopyColumnMapCaseNone,
			EnforceLength:     Bool(true),
			Truncatecolumns:   Bool(true),
			Force:             Bool(true),
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `ALTER STAGE IF EXISTS %s SET URL = 'some url' STORAGE_INTEGRATION = "integration" ENCRYPTION = (TYPE = 'NONE') FILE_FORMAT = (TYPE = JSON) COPY_OPTIONS = (ON_ERROR = CONTINUE SIZE_LIMIT = 123 PURGE = true RETURN_FAILED_ONLY = true MATCH_BY_COLUMN_NAME = NONE ENFORCE_LENGTH = true TRUNCATECOLUMNS = true FORCE = true) COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestStages_AlterDirectoryTable(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterDirectoryTableStageOptions
	defaultOpts := func() *AlterDirectoryTableStageOptions {
		return &AlterDirectoryTableStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterDirectoryTableStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("set directory", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.SetDirectory = &DirectoryTableSet{
			Enable: true,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER STAGE IF EXISTS %s SET DIRECTORY = (ENABLE = true)`, id.FullyQualifiedName())
	})

	t.Run("refresh", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Refresh = &DirectoryTableRefresh{
			Subpath: String("subpath"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER STAGE IF EXISTS %s REFRESH SUBPATH = 'subpath'`, id.FullyQualifiedName())
	})
}

func TestStages_Drop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DropStageOptions
	defaultOpts := func() *DropStageOptions {
		return &DropStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP STAGE IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestStages_Describe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DescribeStageOptions
	defaultOpts := func() *DescribeStageOptions {
		return &DescribeStageOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE STAGE %s", id.FullyQualifiedName())
	})
}

func TestStages_Show(t *testing.T) {
	// Minimal valid ShowStageOptions
	defaultOpts := func() *ShowStageOptions {
		return &ShowStageOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowStageOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW STAGES")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("some pattern"),
		}
		opts.In = &In{
			Schema: NewDatabaseObjectIdentifier("db", "schema"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW STAGES LIKE 'some pattern' IN SCHEMA "db"."schema"`)
	})
}
