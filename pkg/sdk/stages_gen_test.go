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
		}
		opts.Comment = String("some comment")
		opts.Tag = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag-name"),
				Value: "tag-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE TEMPORARY STAGE IF NOT EXISTS %s ENCRYPTION = (TYPE = 'SNOWFLAKE_FULL') DIRECTORY = (ENABLE = true REFRESH_ON_CREATE = true) FILE_FORMAT = (FORMAT_NAME = 'format name') COPY_OPTIONS = (ON_ERROR = CONTINUE) COMMENT = 'some comment' TAG ("tag-name" = 'tag-value')`, id.FullyQualifiedName())
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnS3StageOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.StorageIntegration opts.ExternalStageParams.Credentials]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams", "StorageIntegration", "Credentials"))
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.Credentials.AwsKeyId opts.ExternalStageParams.Credentials.AwsRole]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams.Credentials", "AwsKeyId", "AwsRole"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnGCSStageOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnAzureStageOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.StorageIntegration opts.ExternalStageParams.Credentials]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams", "StorageIntegration", "Credentials"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOnS3CompatibleStageOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterStageOptions", "RenameTo", "SetTags", "UnsetTags"))
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.StorageIntegration opts.ExternalStageParams.Credentials]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams", "StorageIntegration", "Credentials"))
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.Credentials.AwsKeyId opts.ExternalStageParams.Credentials.AwsRole]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams.Credentials", "AwsKeyId", "AwsRole"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.ExternalStageParams.StorageIntegration opts.ExternalStageParams.Credentials]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf(".ExternalStageParams", "StorageIntegration", "Credentials"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}
