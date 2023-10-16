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
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
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
		return fmt.Errorf("number of the AllowedValues must be between 1 and 50")
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
			errs = append(errs, fmt.Errorf("number of the MaskingPolicies must be greater than zero"))
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
			errs = append(errs, fmt.Errorf("number of the MaskingPolicies must be greater than zero"))
		}
	}
	return errors.Join(errs...)
}

func (opts *alterTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.Add, opts.Drop, opts.Set, opts.Unset, opts.Rename); !ok {
		errs = append(errs, errExactlyOneOf("alterTagOptions", "Add", "Drop", "Set", "Unset", "Rename"))
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
	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			errs = append(errs, err)
		}
		if !anyValueSet(opts.Unset.MaskingPolicies, opts.Unset.AllowedValues, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("alterTagOptions.Unset", "MaskingPolicies", "AllowedValues", "Comment"))
		}
	}
	if valueSet(opts.Rename) {
		if !ValidObjectIdentifier(opts.Rename.Name) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	return errors.Join(errs...)
}

func (opts *showTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		errs = append(errs, errExactlyOneOf("showTagOptions.In", "Account", "Database", "Schema"))
	}
	return errors.Join(errs...)
}

func (opts *dropTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *undropTagOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
