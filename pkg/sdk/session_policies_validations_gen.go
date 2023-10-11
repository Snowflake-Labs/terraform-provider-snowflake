package sdk

import "errors"

var (
	_ validatable = new(CreateSessionPolicyOptions)
	_ validatable = new(AlterSessionPolicyOptions)
	_ validatable = new(DropSessionPolicyOptions)
	_ validatable = new(ShowSessionPolicyOptions)
	_ validatable = new(DescribeSessionPolicyOptions)
)

func (opts *CreateSessionPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateSessionPolicyOptions", "OrReplace", "IfNotExists"))
	}
	return errors.Join(errs...)
}

func (opts *AlterSessionPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.RenameTo, opts.Set, opts.SetTags, opts.UnsetTags, opts.Unset); !ok {
		errs = append(errs, errExactlyOneOf("RenameTo", "Set", "SetTags", "UnsetTags", "Unset"))
	}
	if valueSet(opts.Set) {
		if ok := anyValueSet(opts.Set.SessionIdleTimeoutMins, opts.Set.SessionUiIdleTimeoutMins, opts.Set.Comment); !ok {
			errs = append(errs, errAtLeastOneOf("SessionIdleTimeoutMins", "SessionUiIdleTimeoutMins", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if ok := anyValueSet(opts.Unset.SessionIdleTimeoutMins, opts.Unset.SessionUiIdleTimeoutMins, opts.Unset.Comment); !ok {
			errs = append(errs, errAtLeastOneOf("SessionIdleTimeoutMins", "SessionUiIdleTimeoutMins", "Comment"))
		}
	}
	return errors.Join(errs...)
}

func (opts *DropSessionPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *ShowSessionPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}

func (opts *DescribeSessionPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
