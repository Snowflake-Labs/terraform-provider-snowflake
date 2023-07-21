package sdk

import (
	"errors"
)

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
	if ok := exactlyOneValueSet(
		opts.Set,
		opts.Unset,
		opts.Refresh,
	); !ok {
		return errAlterNeedsExactlyOneAction
	}
	if valueSet(opts.Set) {

		return nil
	}
	if valueSet(opts.Unset) {

		return nil
	}
	if valueSet(opts.Refresh) {

		return nil
	}
	return nil
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
	errCopyStatementRequired      = errors.New("copy statement required")
	errAlterNeedsExactlyOneAction = errors.New("alter statement needs exactly one action from: set, unset, refresh")
)
