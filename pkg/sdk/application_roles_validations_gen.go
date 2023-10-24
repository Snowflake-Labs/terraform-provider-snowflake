package sdk

import "errors"

var _ validatable = new(ShowApplicationRoleOptions)

func (opts *ShowApplicationRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.ApplicationName) {
		errs = append(errs, errInvalidIdentifier("ShowApplicationRoleOptions", "ApplicationName"))
	}
	return errors.Join(errs...)
}
