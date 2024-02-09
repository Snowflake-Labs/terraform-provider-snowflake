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
	if opts.tableFunction.policyReferenceFunction.arguments.refEntityDomain == nil {
		errs = append(errs, errNotSet("getForEntityPolicyReferenceOptions", "refEntityDomain"))
	}
	if opts.tableFunction.policyReferenceFunction.arguments.refEntityName == nil {
		errs = append(errs, errNotSet("getForEntityPolicyReferenceOptions", "refEntityName"))
	}
	return errors.Join(errs...)
}
