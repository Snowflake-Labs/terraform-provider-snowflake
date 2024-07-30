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
		if !valueSet(opts.RowAccessPolicy.On) {
			errs = append(errs, errNotSet("CreateViewOptions.RowAccessPolicy", "On"))
		}
	}
	if valueSet(opts.AggregationPolicy) {
		if !ValidObjectIdentifier(opts.AggregationPolicy.AggregationPolicy) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	if valueSet(opts.Columns) {
		for _, columnOption := range opts.Columns {
			if valueSet(columnOption.MaskingPolicy) {
				if !ValidObjectIdentifier(columnOption.MaskingPolicy.MaskingPolicy) {
					errs = append(errs, ErrInvalidObjectIdentifier)
				}
			}
			if valueSet(columnOption.ProjectionPolicy) {
				if !ValidObjectIdentifier(columnOption.ProjectionPolicy.ProjectionPolicy) {
					errs = append(errs, ErrInvalidObjectIdentifier)
				}
			}
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
	if !exactlyOneValueSet(opts.RenameTo, opts.SetComment, opts.UnsetComment, opts.SetSecure, opts.SetChangeTracking, opts.UnsetSecure, opts.SetTags, opts.UnsetTags, opts.AddRowAccessPolicy, opts.DropRowAccessPolicy, opts.DropAndAddRowAccessPolicy, opts.DropAllRowAccessPolicies, opts.SetAggregationPolicy, opts.UnsetAggregationPolicy, opts.SetMaskingPolicyOnColumn, opts.UnsetMaskingPolicyOnColumn, opts.SetProjectionPolicyOnColumn, opts.UnsetProjectionPolicyOnColumn, opts.SetTagsOnColumn, opts.UnsetTagsOnColumn) {
		errs = append(errs, errExactlyOneOf("AlterViewOptions", "RenameTo", "SetComment", "UnsetComment", "SetSecure", "SetChangeTracking", "UnsetSecure", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "SetAggregationPolicy", "UnsetAggregationPolicy", "SetMaskingPolicyOnColumn", "UnsetMaskingPolicyOnColumn", "SetProjectionPolicyOnColumn", "UnsetProjectionPolicyOnColumn", "SetTagsOnColumn", "UnsetTagsOnColumn"))
	}
	if everyValueSet(opts.IfExists, opts.SetSecure) {
		errs = append(errs, errOneOf("AlterViewOptions", "IfExists", "SetSecure"))
	}
	if everyValueSet(opts.IfExists, opts.UnsetSecure) {
		errs = append(errs, errOneOf("AlterViewOptions", "IfExists", "UnsetSecure"))
	}
	if valueSet(opts.AddRowAccessPolicy) {
		if !ValidObjectIdentifier(opts.AddRowAccessPolicy.RowAccessPolicy) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
		if !valueSet(opts.AddRowAccessPolicy.On) {
			errs = append(errs, errNotSet("AlterViewOptions.AddRowAccessPolicy", "On"))
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
			if !valueSet(opts.DropAndAddRowAccessPolicy.Add.On) {
				errs = append(errs, errNotSet("AlterViewOptions.DropAndAddRowAccessPolicy.Add", "On"))
			}
		}
	}
	if valueSet(opts.SetAggregationPolicy) {
		if !ValidObjectIdentifier(opts.SetAggregationPolicy.AggregationPolicy) {
			errs = append(errs, ErrInvalidObjectIdentifier)
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
