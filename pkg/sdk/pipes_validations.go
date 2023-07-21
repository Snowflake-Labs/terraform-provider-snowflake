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
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	//TODO implement me
	panic("implement me")
}

func (opts *PipeDropOptions) validateProp() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	//TODO implement me
	panic("implement me")
}

func (p PipeShowOptions) validateProp() error {
	//TODO implement me
	panic("implement me")
}

func (opts *describePipeOptions) validateProp() error {
	//TODO implement me
	panic("implement me")
}

var (
	errCopyStatementRequired = errors.New("copy statement required")
)
