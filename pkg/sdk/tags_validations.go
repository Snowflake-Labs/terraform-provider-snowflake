package sdk

import (
	"errors"
	"fmt"
)

var (
	_ validatable = new(createTagOptions)
	_ validatable = new(alterTagOptions)
	_ validatable = new(showTagOptions)
	_ validatable = new(dropTagOptions)
	_ validatable = new(undropTagOptions)
	_ validatable = new(AllowedValues)
	_ validatable = new(TagSet)
	_ validatable = new(TagUnset)
)

func (opts *createTagOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) && *opts.OrReplace && *opts.IfNotExists {
		errs = append(errs, errOneOf("OrReplace", "IfNotExists"))
	}
	if valueSet(opts.Comment) && valueSet(opts.AllowedValues) {
		errs = append(errs, errOneOf("Comment", "AllowedValues"))
	}
	if valueSet(opts.AllowedValues) {
		if err := opts.AllowedValues.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *AllowedValues) validate() error {
	if ok := validateIntInRange(len(v.Values), 1, 50); !ok {
		return fmt.Errorf("Number of the AllowedValues must be between 1 and 50")
	}
	return nil
}

func (v *TagSet) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.MaskingPolicies, v.Comment) {
		errs = append(errs, errOneOf("MaskingPolicies", "Comment"))
	}
	if valueSet(v.MaskingPolicies) {
		if ok := validateIntGreaterThanOrEqual(len(v.MaskingPolicies.MaskingPolicies), 1); !ok {
			errs = append(errs, fmt.Errorf("Number of the MaskingPolicies must be greater than zero"))
		}
	}
	return errors.Join(errs...)
}

func (v *TagUnset) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.MaskingPolicies, v.AllowedValues, v.Comment) {
		errs = append(errs, errOneOf("MaskingPolicies", "AllowedValues", "Comment"))
	}
	if valueSet(v.MaskingPolicies) {
		if ok := validateIntGreaterThanOrEqual(len(v.MaskingPolicies.MaskingPolicies), 1); !ok {
			errs = append(errs, fmt.Errorf("Number of the MaskingPolicies must be greater than zero"))
		}
	}
	return errors.Join(errs...)
}

func (opts *alterTagOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(
		opts.Add,
		opts.Drop,
		opts.Set,
		opts.Unset,
		opts.Rename,
	); !ok {
		errs = append(errs, errAlterNeedsExactlyOneAction)
	}
	if valueSet(opts.Add) && valueSet(opts.Add.AllowedValues) {
		if err := opts.Add.AllowedValues.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Drop) && valueSet(opts.Drop.AllowedValues) {
		if err := opts.Drop.AllowedValues.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if unset := opts.Unset; valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			errs = append(errs, err)
		}
		if !anyValueSet(unset.MaskingPolicies, unset.AllowedValues, unset.Comment) {
			errs = append(errs, errAlterNeedsAtLeastOneProperty)
		}
	}
	if valueSet(opts.Rename) {
		if !validObjectidentifier(opts.Rename.Name) {
			errs = append(errs, errInvalidObjectIdentifier)
		}
		if opts.name.DatabaseName() != opts.Rename.Name.DatabaseName() {
			errs = append(errs, errDifferentDatabase)
		}
	}
	return errors.Join(errs...)
}

func (opts *showTagOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, errPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		errs = append(errs, errScopeRequiredForInKeyword)
	}
	return errors.Join(errs...)
}

func (opts *dropTagOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *undropTagOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
