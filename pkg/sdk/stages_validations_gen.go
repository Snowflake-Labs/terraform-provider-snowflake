package sdk

import "errors"

var (
	_ validatable = new(CreateInternalStageOptions)
	_ validatable = new(DropStageOptions)
	_ validatable = new(DescribeStageOptions)
	_ validatable = new(ShowStageOptions)
)

func (opts *CreateInternalStageOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}

func (opts *DropStageOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}

func (opts *DescribeStageOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}

func (opts *ShowStageOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}
