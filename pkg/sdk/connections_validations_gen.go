package sdk

var (
	_ validatable = new(CreateConnectionConnectionOptions)
	_ validatable = new(CreateReplicatedConnectionConnectionOptions)
	_ validatable = new(AlterConnectionFailoverConnectionOptions)
	_ validatable = new(AlterConnectionConnectionOptions)
)

func (opts *CreateConnectionConnectionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *CreateReplicatedConnectionConnectionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.ReplicaOf) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterConnectionFailoverConnectionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !exactlyOneValueSet(opts.EnableConnectionFailover, opts.DisableConnectionFailover, opts.Primary) {
		errs = append(errs, errExactlyOneOf("AlterConnectionFailoverConnectionOptions", "EnableConnectionFailover", "DisableConnectionFailover", "Primary"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterConnectionConnectionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !exactlyOneValueSet(opts.Set, opts.Unset) {
		errs = append(errs, errExactlyOneOf("AlterConnectionConnectionOptions", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterConnectionConnectionOptions.Set", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterConnectionConnectionOptions.Unset", "Comment"))
		}
	}
	return JoinErrors(errs...)
}
