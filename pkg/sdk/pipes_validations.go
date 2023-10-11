package sdk

import (
	"errors"
)

var (
	_ validatable = new(CreatePipeOptions)
	_ validatable = new(AlterPipeOptions)
	_ validatable = new(DropPipeOptions)
	_ validatable = new(ShowPipeOptions)
	_ validatable = new(describePipeOptions)
)

func (opts *CreatePipeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if opts.copyStatement == "" {
		return errCopyStatementRequired
	}
	return nil
}

func (opts *AlterPipeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
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

func (opts *DropPipeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (opts *ShowPipeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		return ErrPatternRequiredForLikeKeyword
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		return errScopeRequiredForInKeyword
	}
	return nil
}

func (opts *describePipeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

var (
	errCopyStatementRequired        = errors.New("copy statement required")
	errScopeRequiredForInKeyword    = errors.New("exactly one scope must be specified for in keyword")
	errAlterNeedsExactlyOneAction   = errors.New("alter statement needs exactly one action from: set, unset, refresh")
	errAlterNeedsAtLeastOneProperty = errors.New("alter statement needs at least one property")
)
