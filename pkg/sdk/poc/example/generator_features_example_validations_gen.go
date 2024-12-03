package example

import "errors"

var _ validatable = new(AlterFeaturesExamplesOptions)

func (opts *AlterFeaturesExamplesOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return errors.Join(errs...)
}
