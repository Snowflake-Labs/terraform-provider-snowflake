package generator

import (
	"fmt"
	"strings"
)

// Interface groups operations for particular object or objects family (e.g. DATABASE ROLE)
type Interface struct {
	// Name is the interface's name, e.g. "DatabaseRoles"
	Name string
	// nameSingular is the prefix/suffix which can be used to create other structs and methods, e.g. "DatabaseRole"
	nameSingular string
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
	// OptsStructFields defines opts used to create SQL for given operation
	OptsStructFields []*Field
	// Validations are top-level validations for given Opts
	Validations []*Validation
}

// OptsName should create a name for opts in a form of OperationObjectOptions where:
// - Operation is e.g. Create
// - Object is e.g. DatabaseRole (singular)
// which together makes CreateDatabaseRoleOptions
func (o *Operation) OptsName() string {
	return fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.nameSingular)
}

// TODO: handle case where validations are on a deeper level (not the immediate one)
func (o *Operation) AdditionalValidations() []*Field {
	var fieldsWithValidations []*Field
	for _, f := range o.OptsStructFields {
		if len(f.Validations) > 0 {
			fieldsWithValidations = append(fieldsWithValidations, f)
		}
	}
	return fieldsWithValidations
}

type Field struct {
	parent      *Field
	Fields      []*Field
	Validations []*Validation

	Name string
	Kind string
	tags map[string][]string
}

func (f *Field) TagsPrintable() string {
	var tagNames = []string{"ddl", "sql"}
	var tagParts []string
	for _, tagName := range tagNames {
		var v, ok = f.tags[tagName]
		if ok {
			tagParts = append(tagParts, fmt.Sprintf(`%s:"%s"`, tagName, strings.Join(v, ",")))
		}
	}
	return fmt.Sprintf("`%s`", strings.Join(tagParts, " "))
}

func (f *Field) KindNoPtr() string {
	kindWithoutPtr, _ := strings.CutPrefix(f.Kind, "*")
	return kindWithoutPtr
}

func (f *Field) NameLowerCased() string {
	return startingWithLowerCase(f.Name)
}

// ValidationType contains all handled validation types. Below validations are marked to be contained here or not:
// - opts not nil - not present here, handled on template level
// - valid identifier - present here, for now put on level containing given field
// - conflicting fields - present here, put on level containing given fields
// - exactly one value set - present here, put on level containing given fields
// - TODO: nested validation conditionally - not present here, handled by putting validations on lower levels
type ValidationType int64

const (
	ValidIdentifier ValidationType = iota
	ConflictingFields
	ExactlyOneValueSet
)

type Validation struct {
	Type       ValidationType
	fieldNames []string
}

func (v *Validation) paramsQuoted() []string {
	var params = make([]string, len(v.fieldNames))
	for i, s := range v.fieldNames {
		params[i] = wrapWith(s, `"`)
	}
	return params
}

// TODO: handle path to field
func (v *Validation) Condition() string {
	switch v.Type {
	case ValidIdentifier:
		return fmt.Sprintf("!validObjectidentifier(%s)", strings.Join(v.fieldNames, ","))
	case ConflictingFields:
		return fmt.Sprintf("everyValueSet(%s)", strings.Join(v.fieldNames, ","))
	case ExactlyOneValueSet:
		return fmt.Sprintf("ok := exactlyOneValueSet(%s); !ok", strings.Join(v.fieldNames, ","))
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
	}
	panic("condition for validation unknown")
}
