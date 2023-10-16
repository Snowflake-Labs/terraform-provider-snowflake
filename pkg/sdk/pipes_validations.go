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
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.copyStatement == "" {
		errs = append(errs, errNotSet("CreatePipeOptions", "copyStatement"))
	}
	return errors.Join(errs...)
}

func (opts *AlterPipeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags, opts.Refresh); !ok {
		errs = append(errs, errExactlyOneOf("AlterPipeOptions", "Set", "Unset", "SetTags", "UnsetTags", "Refresh"))
	}
	if set := opts.Set; valueSet(set) {
		if !anyValueSet(set.ErrorIntegration, set.PipeExecutionPaused, set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterPipeOptions.Set", "ErrorIntegration", "PipeExecutionPaused", "Comment"))
		}
	}
	if unset := opts.Unset; valueSet(unset) {
		if !anyValueSet(unset.PipeExecutionPaused, unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterPipeOptions.Unset", "PipeExecutionPaused", "Comment"))
		}
	}
	if setTags := opts.SetTags; valueSet(setTags) {
		if !valueSet(setTags.Tag) {
			errs = append(errs, errNotSet("AlterPipeOptions.SetTags", "Tag"))
		}
	}
	if unsetTags := opts.UnsetTags; valueSet(unsetTags) {
		if !valueSet(unsetTags.Tag) {
			errs = append(errs, errNotSet("AlterPipeOptions.UnsetTags", "Tag"))
		}
	}
	return errors.Join(errs...)
}

func (opts *DropPipeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *ShowPipeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		errs = append(errs, errExactlyOneOf("ShowPipeOptions", "In.Account", "In.Database", "In.Schema"))
	}
	return errors.Join(errs...)
}

func (opts *describePipeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
