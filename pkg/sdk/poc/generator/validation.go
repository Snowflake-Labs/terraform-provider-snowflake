package generator

import (
	"fmt"
	"strings"
)

// ValidationType contains all handled validation types. Below validations are marked to be contained here or not:
// - opts not nil - not present here, handled on template level
// - valid identifier - present here, for now put on level containing given field
// - conflicting fields - present here, put on level containing given fields
// - exactly one value set - present here, put on level containing given fields
// - at least one value set - present here, put on level containing given fields
// - nested validation conditionally - not present here, handled by putting validations on lower level fields
type ValidationType int64

const (
	ValidIdentifier ValidationType = iota
	ConflictingFields
	ExactlyOneValueSet
	AtLeastOneValueSet
)

type Validation struct {
	Type       ValidationType
	FieldNames []string
}

func NewValidation(validationType ValidationType, fieldNames ...string) *Validation {
	return &Validation{
		Type:       validationType,
		FieldNames: fieldNames,
	}
}

func (f *Field) WithValidation(validationType ValidationType, fieldNames ...string) *Field {
	f.Validations = append(f.Validations, NewValidation(validationType, fieldNames...))
	return f
}

func (v *Validation) paramsQuoted() []string {
	var params = make([]string, len(v.FieldNames))
	for i, s := range v.FieldNames {
		params[i] = wrapWith(s, `"`)
	}
	return params
}

func (v *Validation) fieldsWithPath(field *Field) []string {
	var params = make([]string, len(v.FieldNames))
	for i, s := range v.FieldNames {
		params[i] = fmt.Sprintf("opts%s.%s", field.Path(), s)
	}
	return params
}

func (v *Validation) Condition(field *Field) string {
	switch v.Type {
	case ValidIdentifier:
		return fmt.Sprintf("!validObjectidentifier(%s)", strings.Join(v.fieldsWithPath(field), ","))
	case ConflictingFields:
		return fmt.Sprintf("everyValueSet(%s)", strings.Join(v.fieldsWithPath(field), ","))
	case ExactlyOneValueSet:
		return fmt.Sprintf("ok := exactlyOneValueSet(%s); !ok", strings.Join(v.fieldsWithPath(field), ","))
	case AtLeastOneValueSet:
		return fmt.Sprintf("ok := anyValueSet(%s); !ok", strings.Join(v.fieldsWithPath(field), ","))
	}
	panic("condition for validation unknown")
}

func (v *Validation) Error() string {
	switch v.Type {
	case ValidIdentifier:
		return fmt.Sprintf("ErrInvalidObjectIdentifier")
	case ConflictingFields:
		return fmt.Sprintf("errOneOf(%s)", strings.Join(v.paramsQuoted(), ","))
	case ExactlyOneValueSet:
		return fmt.Sprintf("errExactlyOneOf(%s)", strings.Join(v.paramsQuoted(), ","))
	case AtLeastOneValueSet:
		return fmt.Sprintf("errAtLeastOneOf(%s)", strings.Join(v.paramsQuoted(), ","))
	}
	panic("condition for validation unknown")
}

func (v *Validation) TodoComment(field *Field) string {
	switch v.Type {
	case ValidIdentifier:
		return fmt.Sprintf("validate valid identifier for %v", v.fieldsWithPath(field))
	case ConflictingFields:
		return fmt.Sprintf("validate conflicting fields for %v", v.fieldsWithPath(field))
	case ExactlyOneValueSet:
		return fmt.Sprintf("validate exactly one field from %v is present", v.fieldsWithPath(field))
	case AtLeastOneValueSet:
		return fmt.Sprintf("validate at least one of fields %v set", v.fieldsWithPath(field))
	}
	panic("condition for validation unknown")
}
