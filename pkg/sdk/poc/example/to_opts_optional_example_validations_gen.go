package example

import "errors"

var _ validatable = new(AlterToOptsOptionalExampleOptions)

func (opts *AlterToOptsOptionalExampleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return errors.Join(errs...)
}
