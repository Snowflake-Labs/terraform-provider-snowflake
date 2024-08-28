package sdk

import "errors"

var _ validatable = new(GetForEntityDataMetricFunctionReferenceOptions)

func (opts *GetForEntityDataMetricFunctionReferenceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !valueSet(opts.parameters) {
		errs = append(errs, errNotSet("GetForEntityDataMetricFunctionReferenceOptions", "parameters"))
	} else {
		if !valueSet(opts.parameters.arguments) {
			errs = append(errs, errNotSet("dataMetricFunctionReferenceParameters", "arguments"))
		} else {
			if opts.parameters.arguments.refEntityDomain == nil {
				errs = append(errs, errNotSet("dataMetricFunctionReferenceFunctionArguments", "refEntityDomain"))
			}
			if opts.parameters.arguments.refEntityName == nil {
				errs = append(errs, errNotSet("dataMetricFunctionReferenceFunctionArguments", "refEntityName"))
			}
		}
	}
	return errors.Join(errs...)
}
