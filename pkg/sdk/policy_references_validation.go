package sdk

import (
	"errors"
)

var _ validatable = new(getForEntityPolicyReferenceOptions)

func (opts *getForEntityPolicyReferenceOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !opts.select_ {
		errs = append(errs, errors.New("select_ must be true"))
	}
	if !opts.asterisk {
		errs = append(errs, errors.New("asterisk must be true"))
	}
	if !opts.from {
		errs = append(errs, errors.New("from must be true"))
	}
	if opts.tableFunction == nil {
		errs = append(errs, errors.New("tableFunction must not be nil"))
	}
	return errors.Join(errs...)
}
