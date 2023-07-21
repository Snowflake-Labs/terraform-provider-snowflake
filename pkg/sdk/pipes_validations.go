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
	if set := opts.Set; valueSet(set) {
		if !anyValueSet(set.ErrorIntegration, set.PipeExecutionPaused, set.Tag, set.Comment) {
			return errAlterNeedsAtLeastOneProperty
		}
		if valueSet(set.Tag) {
			if !everyValueNil(set.ErrorIntegration, set.PipeExecutionPaused, set.Comment) {
				return errCannotAlterOtherPropertyWithTag
			}
		}
	}
	if unset := opts.Unset; valueSet(unset) {
		if !anyValueSet(unset.PipeExecutionPaused, unset.Tag, unset.Comment) {
			return errAlterNeedsAtLeastOneProperty
		}
		if valueSet(unset.Tag) {
			if !everyValueNil(unset.PipeExecutionPaused, unset.Comment) {
				return errCannotAlterOtherPropertyWithTag
			}
		}
	}
	if valueSet(opts.Refresh) {
		// no validations needed
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
	errCopyStatementRequired           = errors.New("copy statement required")
	errAlterNeedsExactlyOneAction      = errors.New("alter statement needs exactly one action from: set, unset, refresh")
	errAlterNeedsAtLeastOneProperty    = errors.New("alter statement needs at least one property")
	errCannotAlterOtherPropertyWithTag = errors.New("cannot alter both tag and other property in the same statement")
)
