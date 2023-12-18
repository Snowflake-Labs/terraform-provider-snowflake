package sdk

var (
	_ validatable = new(CreateEventTableOptions)
	_ validatable = new(ShowEventTableOptions)
	_ validatable = new(DescribeEventTableOptions)
	_ validatable = new(DropEventTableOptions)
	_ validatable = new(AlterEventTableOptions)
)

func (opts *CreateEventTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateEventTableOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *ShowEventTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeEventTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DropEventTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterEventTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags, opts.AddRowAccessPolicy, opts.DropRowAccessPolicy, opts.DropAndAddRowAccessPolicy, opts.DropAllRowAccessPolicies, opts.ClusteringAction, opts.SearchOptimizationAction) {
		errs = append(errs, errExactlyOneOf("AlterEventTableOptions", "RenameTo", "Set", "Unset", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "ClusteringAction", "SearchOptimizationAction"))
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
