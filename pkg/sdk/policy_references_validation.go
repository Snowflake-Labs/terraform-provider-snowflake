package sdk

import (
	"errors"
)

var _ validatable = new(getForEntityPolicyReferenceOptions)

func (opts *getForEntityPolicyReferenceOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(opts.parameters) {
		errs = append(errs, errNotSet("getForEntityPolicyReferenceOptions", "parameters"))
	} else {
		if !valueSet(opts.parameters.arguments) {
			errs = append(errs, errNotSet("policyReferenceParameters", "arguments"))
		} else {
			if opts.parameters.arguments.refEntityDomain == nil {
				errs = append(errs, errNotSet("policyReferenceFunctionArguments", "refEntityDomain"))
			}
			if opts.parameters.arguments.refEntityName == nil {
				errs = append(errs, errNotSet("policyReferenceFunctionArguments", "refEntityName"))
			}
		}
	}
	return errors.Join(errs...)
}
