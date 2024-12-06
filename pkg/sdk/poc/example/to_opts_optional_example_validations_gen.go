package example

import "errors"

var (
	_ validatable = new(CreateToOptsOptionalExampleOptions)
	_ validatable = new(AlterToOptsOptionalExampleOptions)
)

func (opts *CreateToOptsOptionalExampleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return errors.Join(errs...)
}

func (opts *AlterToOptsOptionalExampleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return errors.Join(errs...)
}
