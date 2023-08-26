package sdk

import (
	"errors"
)

var (
	_ validatable = new(PipeCreateOptions)
	_ validatable = new(PipeAlterOptions)
	_ validatable = new(PipeDropOptions)
	_ validatable = new(PipeShowOptions)
	_ validatable = new(describePipeOptions)
)

func (opts *PipeCreateOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if opts.copyStatement == "" {
		return errCopyStatementRequired
	}
	return nil
}

func (opts *PipeAlterOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if ok := exactlyOneValueSet(
		opts.Set,
		opts.Unset,
		opts.SetTags,
		opts.UnsetTags,
		opts.Refresh,
	); !ok {
		return errAlterNeedsExactlyOneAction
	}
	if set := opts.Set; valueSet(set) {
		if !anyValueSet(set.ErrorIntegration, set.PipeExecutionPaused, set.Comment) {
			return errAlterNeedsAtLeastOneProperty
		}
	}
	if unset := opts.Unset; valueSet(unset) {
		if !anyValueSet(unset.PipeExecutionPaused, unset.Comment) {
			return errAlterNeedsAtLeastOneProperty
		}
	}
	if setTags := opts.SetTags; valueSet(setTags) {
		if !valueSet(setTags.Tag) {
			return errAlterNeedsAtLeastOneProperty
		}
	}
	if unsetTags := opts.UnsetTags; valueSet(unsetTags) {
		if !valueSet(unsetTags.Tag) {
			return errAlterNeedsAtLeastOneProperty
		}
	}
	return nil
}

func (opts *PipeDropOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (opts *PipeShowOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		return errPatternRequiredForLikeKeyword
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		return errScopeRequiredForInKeyword
	}
	return nil
}

func (opts *describePipeOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

var (
	errNilOptions                    = errors.New("options cannot be nil")
	errCopyStatementRequired         = errors.New("copy statement required")
	errPatternRequiredForLikeKeyword = errors.New("pattern must be specified for like keyword")
	errScopeRequiredForInKeyword     = errors.New("exactly one scope must be specified for in keyword")
	errAlterNeedsExactlyOneAction    = errors.New("alter statement needs exactly one action from: set, unset, refresh")
	errAlterNeedsAtLeastOneProperty  = errors.New("alter statement needs at least one property")
)
