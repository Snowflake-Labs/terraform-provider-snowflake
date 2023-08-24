package generator

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// Interface groups operations for particular object or objects family (e.g. DATABASE ROLE)
type Interface struct {
	// Name is the interface's name, e.g. "DatabaseRoles"
	Name string
	// NameSingular is the prefix/suffix which can be used to create other structs and methods, e.g. "DatabaseRole"
	NameSingular string
	// Operations contains all operations for given interface
	Operations []*Operation
	// IdentifierKind keeps identifier of the underlying object (e.g. DatabaseObjectIdentifier)
	IdentifierKind string
}

func (i *Interface) NameLowerCased() string {
	return startingWithLowerCase(i.Name)
}

// Operation defines a single operation for given object or objects family (e.g. CREATE DATABASE ROLE)
type Operation struct {
	// Name is the operation's name, e.g. "Create"
	Name string
	// ObjectInterface points to the containing interface
	ObjectInterface *Interface
	// Doc is the URL for the doc used to create given operation, e.g. https://docs.snowflake.com/en/sql-reference/sql/create-database-role
	Doc string
	// Validations are top-level validations for given Opts
	Validations []*Validation
	// OptsField defines opts used to create SQL for given operation
	OptsField *Field
}

type Field struct {
	Parent      *Field
	Fields      []*Field
	Validations []*Validation

	Name     string
	Kind     string
	Tags     map[string][]string
	Required bool
}

// TODO: handle case where validations are on a deeper level (not the immediate one)
func (field *Field) AdditionalValidations() []*Field {
	var fieldsWithValidations []*Field
	for _, f := range field.Fields {
		if len(f.Validations) > 0 {
			fieldsWithValidations = append(fieldsWithValidations, f)
		}
	}
	return fieldsWithValidations
}

func (field *Field) TagsPrintable() string {
	var tagNames = []string{"ddl", "sql"}
	var tagParts []string
	for _, tagName := range tagNames {
		var v, ok = field.Tags[tagName]
		if ok {
			tagParts = append(tagParts, fmt.Sprintf(`%s:"%s"`, tagName, strings.Join(v, ",")))
		}
	}
	return fmt.Sprintf("`%s`", strings.Join(tagParts, " "))
}

func (field *Field) KindNoPtr() string {
	kindWithoutPtr, _ := strings.CutPrefix(field.Kind, "*")
	return kindWithoutPtr
}

func (field *Field) NameLowerCased() string {
	return startingWithLowerCase(field.Name)
}

func (field *Field) IsStruct() bool {
	return len(field.Fields) > 0
}

func (field *Field) ShouldBeInDto() bool {
	return !slices.Contains(field.Tags["ddl"], "static")
}

func (field *Field) DtoKind() string {
	if field.Parent == nil {
		withoutSuffix, _ := strings.CutSuffix(field.Kind, "Options")
		return fmt.Sprintf("%sRequest", withoutSuffix)
	} else if field.IsStruct() {
		return fmt.Sprintf("%sRequest", field.Kind)
	} else {
		return field.Kind
	}
}

func (field *Field) DtoDecl() string {
	if field.Parent == nil {
		withoutSuffix, _ := strings.CutSuffix(field.KindNoPtr(), "Options")
		return fmt.Sprintf("%sRequest", withoutSuffix)
	} else if field.IsStruct() {
		return fmt.Sprintf("%sRequest", field.KindNoPtr())
	} else {
		return field.KindNoPtr()
	}
}

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

func (v *Validation) paramsQuoted() []string {
	var params = make([]string, len(v.FieldNames))
	for i, s := range v.FieldNames {
		params[i] = wrapWith(s, `"`)
	}
	return params
}

// TODO: handle path to field
func (v *Validation) Condition() string {
	switch v.Type {
	case ValidIdentifier:
		return fmt.Sprintf("!validObjectidentifier(%s)", strings.Join(v.FieldNames, ","))
	case ConflictingFields:
		return fmt.Sprintf("everyValueSet(%s)", strings.Join(v.FieldNames, ","))
	case ExactlyOneValueSet:
		return fmt.Sprintf("ok := exactlyOneValueSet(%s); !ok", strings.Join(v.FieldNames, ","))
	case AtLeastOneValueSet:
		return fmt.Sprintf("ok := anyValueSet(%s); !ok", strings.Join(v.FieldNames, ","))
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

func (v *Validation) TodoComment() string {
	switch v.Type {
	case ValidIdentifier:
		return fmt.Sprintf("// TODO: validate valid identifier for %v", v.FieldNames)
	case ConflictingFields:
		return fmt.Sprintf("// TODO: validate conflicting fields for %v", v.FieldNames)
	case ExactlyOneValueSet:
		return fmt.Sprintf("// TODO: validate exactly one field from %v is present", v.FieldNames)
	case AtLeastOneValueSet:
		return fmt.Sprintf("// TODO: validate at least one of fields %v set", v.FieldNames)
	}
	panic("condition for validation unknown")
}
