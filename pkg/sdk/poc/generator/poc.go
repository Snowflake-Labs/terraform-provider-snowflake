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

// NameLowerCased returns interface name starting with a lower case letter
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
	// OptsField defines opts used to create SQL for given operation
	OptsField *Field
}

// Field defines properties of a single field or struct (by defining Fields)
type Field struct {
	// Parent allows to traverse fields hierarchy more easily, nil for root
	Parent *Field
	// Fields defines children, use for struct fields
	Fields []*Field
	// Validations defines validations on given field level (e.g. oneOf for children)
	Validations []*Validation

	// Name is how field is called in parent struct
	Name string
	// Kind is fields type (e.g. string, *bool)
	Kind string
	// Tags should contain ddl and sql tags used for SQL generation
	Tags map[string][]string
	// Required is used to mark fields which are essential (it's used e.g. for DTO builders generation)
	Required bool
}

// HasAnyValidationInSubtree checks if any validations are present from current field level downwards
func (field *Field) HasAnyValidationInSubtree() bool {
	if len(field.Validations) > 0 {
		return true
	}
	for _, f := range field.Fields {
		if f.HasAnyValidationInSubtree() {
			return true
		}
	}
	return false
}

// TagsPrintable defines how tags are printed in options structs, it ensures the same order of tags for every field
func (field *Field) TagsPrintable() string {
	tagNames := []string{"ddl", "sql"}
	var tagParts []string
	for _, tagName := range tagNames {
		v, ok := field.Tags[tagName]
		if ok {
			tagParts = append(tagParts, fmt.Sprintf(`%s:"%s"`, tagName, strings.Join(v, ",")))
		}
	}
	return fmt.Sprintf("`%s`", strings.Join(tagParts, " "))
}

// KindNoPtr return field's Kind but without pointer
func (field *Field) KindNoPtr() string {
	kindWithoutPtr, _ := strings.CutPrefix(field.Kind, "*")
	return kindWithoutPtr
}

// IsStruct checks if field is a struct
func (field *Field) IsStruct() bool {
	return len(field.Fields) > 0
}

// ShouldBeInDto checks if field is not some static SQL field which should not be interacted with by SDK user
// TODO: this is a very naive implementation, consider fixing it with DSL builder connection
func (field *Field) ShouldBeInDto() bool {
	return !slices.Contains(field.Tags["ddl"], "static")
}

// IsRoot checks if field is at the top of field hierarchy, basically it is true for Option structs
func (field *Field) IsRoot() bool {
	return field.Parent == nil
}

// Path returns the way through the tree to the top, with dot separator (e.g. .SomeField.SomeChild)
func (field *Field) Path() string {
	if field.IsRoot() {
		return ""
	} else {
		return fmt.Sprintf("%s.%s", field.Parent.Path(), field.Name)
	}
}

// DtoKind returns what should be fields kind in generated DTO, because it may differ from Kind
func (field *Field) DtoKind() string {
	if field.IsRoot() {
		withoutSuffix, _ := strings.CutSuffix(field.Kind, "Options")
		return fmt.Sprintf("%sRequest", withoutSuffix)
	}
	if field.IsStruct() {
		return fmt.Sprintf("%sRequest", field.Kind)
	}
	return field.Kind
}

// DtoDecl returns how struct should be declared in generated DTO (e.g. definition is without a pointer)
func (field *Field) DtoDecl() string {
	if field.Parent == nil {
		withoutSuffix, _ := strings.CutSuffix(field.KindNoPtr(), "Options")
		return fmt.Sprintf("%sRequest", withoutSuffix)
	}
	if field.IsStruct() {
		return fmt.Sprintf("%sRequest", field.KindNoPtr())
	}
	return field.KindNoPtr()
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

type Validation struct { //nolint
	Type       ValidationType
	Struct     *Field
	FieldNames []string
}

func (v *Validation) paramsQuoted() []string {
	params := make([]string, len(v.FieldNames))
	for i, s := range v.FieldNames {
		params[i] = wrapWith(s, `"`)
	}
	return params
}

func (v *Validation) fieldsWithPath(field *Field) []string {
	params := make([]string, len(v.FieldNames))
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
		return fmt.Sprintf("ErrInvalidObjectIdentifier") //nolint
	case ConflictingFields:
		return fmt.Sprintf("errOneOf(%s, %s)", wrapWith(v.Struct.Name, `"`), strings.Join(v.paramsQuoted(), ","))
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
		return fmt.Sprintf("// TODO: validate valid identifier for %v", v.fieldsWithPath(field))
	case ConflictingFields:
		return fmt.Sprintf("// TODO: validate conflicting fields for %v", v.fieldsWithPath(field))
	case ExactlyOneValueSet:
		return fmt.Sprintf("// TODO: validate exactly one field from %v is present", v.fieldsWithPath(field))
	case AtLeastOneValueSet:
		return fmt.Sprintf("// TODO: validate at least one of fields %v set", v.fieldsWithPath(field))
	}
	panic("condition for validation unknown")
}
