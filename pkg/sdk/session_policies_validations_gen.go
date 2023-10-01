package sdk

import "errors"

var (
	_ validatable = new(CreateSessionPolicyOptions)
	_ validatable = new(AlterSessionPolicyOptions)
)

func (opts *CreateSessionPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *AlterSessionPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.RenameTo, opts.Set, opts.SetTags, opts.UnsetTags, opts.Unset); !ok {
		errs = append(errs, errExactlyOneOf("RenameTo", "Set", "SetTags", "UnsetTags", "Unset"))
	}
	if valueSet(opts.RenameTo) && !validObjectidentifier(opts.RenameTo) {
		errs = append(errs, errInvalidObjectIdentifier)
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
