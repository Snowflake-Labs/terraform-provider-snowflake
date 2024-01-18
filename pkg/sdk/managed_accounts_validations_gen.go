package sdk

var (
	_ validatable = new(CreateManagedAccountOptions)
	_ validatable = new(DropManagedAccountOptions)
	_ validatable = new(ShowManagedAccountOptions)
)

func (opts *CreateManagedAccountOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.CreateManagedAccountParams) {
		if !valueSet(opts.CreateManagedAccountParams.AdminName) {
			errs = append(errs, errNotSet("CreateManagedAccountOptions.CreateManagedAccountParams", "AdminName"))
		}
		if !valueSet(opts.CreateManagedAccountParams.AdminPassword) {
			errs = append(errs, errNotSet("CreateManagedAccountOptions.CreateManagedAccountParams", "AdminPassword"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropManagedAccountOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowManagedAccountOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
