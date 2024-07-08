package sdk

var (
	_ validatable = new(CreateStreamlitOptions)
	_ validatable = new(AlterStreamlitOptions)
	_ validatable = new(DropStreamlitOptions)
	_ validatable = new(ShowStreamlitOptions)
	_ validatable = new(DescribeStreamlitOptions)
)

func (opts *CreateStreamlitOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.Warehouse != nil && !ValidObjectIdentifier(opts.Warehouse) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateStreamlitOptions", "IfNotExists", "OrReplace"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterStreamlitOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.RenameTo != nil && !ValidObjectIdentifier(opts.RenameTo) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.Set, opts.Unset) {
		errs = append(errs, errExactlyOneOf("AlterStreamlitOptions", "RenameTo", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		if opts.Set.Warehouse != nil && !ValidObjectIdentifier(opts.Set.Warehouse) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
		if !anyValueSet(opts.Set.RootLocation, opts.Set.MainFile, opts.Set.Warehouse, opts.Set.ExternalAccessIntegrations, opts.Set.Comment, opts.Set.Title) {
			errs = append(errs, errAtLeastOneOf("AlterStreamlitOptions.Set", "RootLocation", "MainFile", "Warehouse", "ExternalAccessIntegrations", "Comment", "Title"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.QueryWarehouse, opts.Unset.Title, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterStreamlitOptions.Unset", "QueryWarehouse", "Title", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropStreamlitOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowStreamlitOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeStreamlitOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
