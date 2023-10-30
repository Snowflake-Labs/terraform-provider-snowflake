package sdk

var (
	_ validatable = new(CreateInternalStageOptions)
	_ validatable = new(CreateOnS3StageOptions)
	_ validatable = new(CreateOnGCSStageOptions)
	_ validatable = new(CreateOnAzureStageOptions)
	_ validatable = new(CreateOnS3CompatibleStageOptions)
	_ validatable = new(AlterStageOptions)
	_ validatable = new(AlterInternalStageStageOptions)
	_ validatable = new(AlterExternalS3StageStageOptions)
	_ validatable = new(AlterExternalGCSStageStageOptions)
	_ validatable = new(AlterExternalAzureStageStageOptions)
	_ validatable = new(AlterDirectoryTableStageOptions)
	_ validatable = new(DropStageOptions)
	_ validatable = new(DescribeStageOptions)
	_ validatable = new(ShowStageOptions)
)

func (opts *CreateInternalStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *CreateOnS3StageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *CreateOnGCSStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *CreateOnAzureStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *CreateOnS3CompatibleStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *AlterStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *AlterInternalStageStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *AlterExternalS3StageStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *AlterExternalGCSStageStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *AlterExternalAzureStageStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *AlterDirectoryTableStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DropStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *ShowStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
