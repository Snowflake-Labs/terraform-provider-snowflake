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
	if opts.tableFunction == nil {
		errs = append(errs, errors.New("tableFunction must not be nil"))
	}
	return errors.Join(errs...)
}
