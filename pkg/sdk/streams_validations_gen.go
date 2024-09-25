package sdk

var (
	_ validatable = new(CreateOnTableStreamOptions)
	_ validatable = new(CreateOnExternalTableStreamOptions)
	_ validatable = new(CreateOnDirectoryTableStreamOptions)
	_ validatable = new(CreateOnViewStreamOptions)
	_ validatable = new(CloneStreamOptions)
	_ validatable = new(AlterStreamOptions)
	_ validatable = new(DropStreamOptions)
	_ validatable = new(ShowStreamOptions)
	_ validatable = new(DescribeStreamOptions)
)

func (opts *CreateOnTableStreamOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.TableId) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateOnTableStreamOptions", "IfNotExists", "OrReplace"))
	}
	if valueSet(opts.On) {
		if !exactlyOneValueSet(opts.On.At, opts.On.Before) {
			errs = append(errs, errExactlyOneOf("CreateOnTableStreamOptions.On", "At", "Before"))
		}
		if valueSet(opts.On.Statement) {
			if !exactlyOneValueSet(opts.On.Statement.Timestamp, opts.On.Statement.Offset, opts.On.Statement.Statement, opts.On.Statement.Stream) {
				errs = append(errs, errExactlyOneOf("CreateOnTableStreamOptions.On.Statement", "Timestamp", "Offset", "Statement", "Stream"))
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *CreateOnExternalTableStreamOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.ExternalTableId) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateOnExternalTableStreamOptions", "IfNotExists", "OrReplace"))
	}
	if valueSet(opts.On) {
		if !exactlyOneValueSet(opts.On.At, opts.On.Before) {
			errs = append(errs, errExactlyOneOf("CreateOnExternalTableStreamOptions.On", "At", "Before"))
		}
		if valueSet(opts.On.Statement) {
			if !exactlyOneValueSet(opts.On.Statement.Timestamp, opts.On.Statement.Offset, opts.On.Statement.Statement, opts.On.Statement.Stream) {
				errs = append(errs, errExactlyOneOf("CreateOnExternalTableStreamOptions.On.Statement", "Timestamp", "Offset", "Statement", "Stream"))
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *CreateOnDirectoryTableStreamOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.StageId) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateOnDirectoryTableStreamOptions", "IfNotExists", "OrReplace"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateOnViewStreamOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.ViewId) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateOnViewStreamOptions", "IfNotExists", "OrReplace"))
	}
	if valueSet(opts.On) {
		if !exactlyOneValueSet(opts.On.At, opts.On.Before) {
			errs = append(errs, errExactlyOneOf("CreateOnViewStreamOptions.On", "At", "Before"))
		}
		if valueSet(opts.On.Statement) {
			if !exactlyOneValueSet(opts.On.Statement.Timestamp, opts.On.Statement.Offset, opts.On.Statement.Statement, opts.On.Statement.Stream) {
				errs = append(errs, errExactlyOneOf("CreateOnViewStreamOptions.On.Statement", "Timestamp", "Offset", "Statement", "Stream"))
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *CloneStreamOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterStreamOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfExists, opts.UnsetTags) {
		errs = append(errs, errOneOf("AlterStreamOptions", "IfExists", "UnsetTags"))
	}
	if !exactlyOneValueSet(opts.SetComment, opts.UnsetComment, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterStreamOptions", "SetComment", "UnsetComment", "SetTags", "UnsetTags"))
	}
	return JoinErrors(errs...)
}

func (opts *DropStreamOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowStreamOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeStreamOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
