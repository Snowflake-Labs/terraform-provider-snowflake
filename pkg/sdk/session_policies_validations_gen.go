package sdk

import "errors"

var (
	_ validatable = new(CreateSessionPolicyOptions)
)

func (opts *CreateSessionPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
