package sdk

import "errors"

var (
	_ validatable = new(CreateOnTableStreamOptions)
	_ validatable = new(CreateOnExternalTableStreamOptions)
	_ validatable = new(CreateOnStageStreamOptions)
	_ validatable = new(CreateOnViewStreamOptions)
	_ validatable = new(CloneStreamOptions)
	_ validatable = new(AlterStreamOptions)
	_ validatable = new(DropStreamOptions)
	_ validatable = new(ShowStreamOptions)
	_ validatable = new(DescribeStreamOptions)
)

func (opts *CreateOnTableStreamOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if !validObjectidentifier(opts.TableId) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateOnTableStreamOptions", "IfNotExists", "OrReplace"))
	}
	return errors.Join(errs...)
}

func (opts *CreateOnExternalTableStreamOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if !validObjectidentifier(opts.ExternalTableId) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateOnExternalTableStreamOptions", "IfNotExists", "OrReplace"))
	}
	return errors.Join(errs...)
}

func (opts *CreateOnStageStreamOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if !validObjectidentifier(opts.StageId) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateOnStageStreamOptions", "IfNotExists", "OrReplace"))
	}
	return errors.Join(errs...)
}

func (opts *CreateOnViewStreamOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if !validObjectidentifier(opts.ViewId) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateOnViewStreamOptions", "IfNotExists", "OrReplace"))
	}
	return errors.Join(errs...)
}

func (opts *CloneStreamOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *AlterStreamOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfExists, opts.UnsetTags) {
		errs = append(errs, errOneOf("AlterStreamOptions", "IfExists", "UnsetTags"))
	}
	if ok := exactlyOneValueSet(opts.SetComment, opts.UnsetComment, opts.SetTags, opts.UnsetTags); !ok {
		errs = append(errs, errExactlyOneOf("SetComment", "UnsetComment", "SetTags", "UnsetTags"))
	}
	return errors.Join(errs...)
}

func (opts *DropStreamOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *ShowStreamOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}

func (opts *DescribeStreamOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
