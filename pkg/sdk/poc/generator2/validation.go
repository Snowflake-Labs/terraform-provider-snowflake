package generator2

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

func (v *Validation) paramsQuoted() []string {
	var params = make([]string, len(v.FieldNames))
	for i, s := range v.FieldNames {
		params[i] = wrapWith(s, `"`)
	}
	return params
}

func (v *Validation) fieldsWithPath() []string {
	var params = make([]string, len(v.FieldNames))
	for i, s := range v.FieldNames {
		params[i] = fmt.Sprintf("opts.%s", s)
	}
	return params
}

func (v *Validation) fieldsWithPathJoined() string {
	return strings.Join(v.fieldsWithPath(), ",")
}

func (v *Validation) Condition(s *Struct) string {
	switch v.Type {
	case ValidIdentifier:
		return fmt.Sprintf("!validObjectidentifier(%s)", v.fieldsWithPathJoined())
	case ConflictingFields:
		return fmt.Sprintf("everyValueSet(%s)", v.fieldsWithPathJoined())
	case ExactlyOneValueSet:
		return fmt.Sprintf("ok := exactlyOneValueSet(%s); !ok", v.fieldsWithPathJoined())
	case AtLeastOneValueSet:
		return fmt.Sprintf("ok := anyValueSet(%s); !ok", v.fieldsWithPathJoined())
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

//func (v *Validation) TodoComment(field *Field) string {
//	switch v.Type {
//	case ValidIdentifier:
//		return fmt.Sprintf("// TODO: validate valid identifier for %v", v.fieldsWithPath(field))
//	case ConflictingFields:
//		return fmt.Sprintf("// TODO: validate conflicting fields for %v", v.fieldsWithPath(field))
//	case ExactlyOneValueSet:
//		return fmt.Sprintf("// TODO: validate exactly one field from %v is present", v.fieldsWithPath(field))
//	case AtLeastOneValueSet:
//		return fmt.Sprintf("// TODO: validate at least one of fields %v set", v.fieldsWithPath(field))
//	}
//	panic("condition for validation unknown")
//}
