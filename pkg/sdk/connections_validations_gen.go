package sdk

var (
	_ validatable = new(CreateConnectionOptions)
	_ validatable = new(AlterConnectionOptions)
	_ validatable = new(DropConnectionOptions)
	_ validatable = new(ShowConnectionOptions)
)

func (opts *CreateConnectionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.AsReplicaOf != nil && !ValidObjectIdentifier(opts.AsReplicaOf) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterConnectionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !exactlyOneValueSet(opts.EnableConnectionFailover, opts.DisableConnectionFailover, opts.Primary, opts.Set, opts.Unset) {
		errs = append(errs, errExactlyOneOf("AlterConnectionOptions", "EnableConnectionFailover", "DisableConnectionFailover", "Primary", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterConnectionOptions.Set", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterConnectionOptions.Unset", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropConnectionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowConnectionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
