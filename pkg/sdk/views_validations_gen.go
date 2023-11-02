package sdk

var (
	_ validatable = new(CreateViewOptions)
	_ validatable = new(AlterViewOptions)
	_ validatable = new(DropViewOptions)
	_ validatable = new(ShowViewOptions)
	_ validatable = new(DescribeViewOptions)
)

func (opts *CreateViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateViewOptions", "OrReplace", "IfNotExists"))
	}
	if valueSet(opts.RowAccessPolicy) {
		if !ValidObjectIdentifier(opts.RowAccessPolicy.RowAccessPolicy) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.SetComment, opts.UnsetComment, opts.SetSecure, opts.SetChangeTracking, opts.UnsetSecure, opts.SetTags, opts.UnsetTags, opts.AddRowAccessPolicy, opts.DropRowAccessPolicy, opts.DropAndAddRowAccessPolicy, opts.DropAllRowAccessPolicies, opts.SetMaskingPolicyOnColumn, opts.UnsetMaskingPolicyOnColumn, opts.SetTagsOnColumn, opts.UnsetTagsOnColumn) {
		errs = append(errs, errExactlyOneOf("AlterViewOptions", "RenameTo", "SetComment", "UnsetComment", "SetSecure", "SetChangeTracking", "UnsetSecure", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "SetMaskingPolicyOnColumn", "UnsetMaskingPolicyOnColumn", "SetTagsOnColumn", "UnsetTagsOnColumn"))
	}
	if valueSet(opts.AddRowAccessPolicy) {
		if !ValidObjectIdentifier(opts.AddRowAccessPolicy.RowAccessPolicy) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	if valueSet(opts.DropRowAccessPolicy) {
		if !ValidObjectIdentifier(opts.DropRowAccessPolicy.RowAccessPolicy) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	if valueSet(opts.DropAndAddRowAccessPolicy) {
		if valueSet(opts.DropAndAddRowAccessPolicy.Drop) {
			if !ValidObjectIdentifier(opts.DropAndAddRowAccessPolicy.Drop.RowAccessPolicy) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
		if valueSet(opts.DropAndAddRowAccessPolicy.Add) {
			if !ValidObjectIdentifier(opts.DropAndAddRowAccessPolicy.Add.RowAccessPolicy) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
