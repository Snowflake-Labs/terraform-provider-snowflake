package sdk

import "errors"

func (opts *PipeCreateOptions) validateProp() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if opts.CopyStatement == "" {
		return errCopyStatementRequired
	}
	return nil
}

func (opts *PipeAlterOptions) validateProp() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	//TODO implement me
	panic("implement me")
}

func (opts *PipeDropOptions) validateProp() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (opts *PipeShowOptions) validateProp() error {
	if opts == nil {
		return ErrNilOptions
	}
	//TODO implement me
	panic("implement me")
}

func (opts *describePipeOptions) validateProp() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

var (
	errCopyStatementRequired = errors.New("copy statement required")
)
