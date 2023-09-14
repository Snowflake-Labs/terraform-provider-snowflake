package sdk

import "errors"

var (
	_ validatable = new(CreateNetworkPolicyOptions)
	_ validatable = new(DropNetworkPolicyOptions)
	_ validatable = new(ShowNetworkPolicyOptions)
	_ validatable = new(DescribeNetworkPolicyOptions)
)

func (opts *CreateNetworkPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *DropNetworkPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}

func (opts *ShowNetworkPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}

func (opts *DescribeNetworkPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}
