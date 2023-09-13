package generator

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
)

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

func NewField(name string, kind string, tagBuilder *TagBuilder, transformer FieldTransformer) *Field {
	var tags map[string][]string
	if tagBuilder != nil {
		tags = tagBuilder.Build()
	} else {
		tags = make(map[string][]string)
	}
	f := &Field{
		Name: name,
		Kind: kind,
		Tags: tags,
	}
	if transformer != nil {
		return transformer.Transform(f)
	}
	return f
}

// TODO can be of type querystruct and be converted into field under the hood
func QueryStruct(name string) *Field {
	return NewField(name, name, nil, nil)
}

func (f *Field) QueryStructField(queryStruct *Field, name string, kind string, transformer FieldTransformer) *Field {
	newField := NewField(name, kind, nil, transformer).
		WithFields(queryStruct.Fields...).
		WithValidations(queryStruct.Validations...)
	f.Fields = append(f.Fields, newField)
	return f
}

func (f *Field) WithField(fields *Field) *Field {
	f.Fields = append(f.Fields, fields)
	return f
}

func (f *Field) WithFields(fields ...*Field) *Field {
	f.Fields = append(f.Fields, fields...)
	return f
}

func (f *Field) WithValidations(validations ...*Validation) *Field {
	f.Validations = validations
	return f
}

// TODO find a way to mark certain field as required
//func (f *Field) WithRequired(required bool) *Field {
//	f.Required = required
//	return f
//}

// HasAnyValidationInSubtree checks if any validations are present from current field level downwards
func (f *Field) HasAnyValidationInSubtree() bool {
	if len(f.Validations) > 0 {
		return true
	}
	for _, f := range f.Fields {
		if f.HasAnyValidationInSubtree() {
			return true
		}
	}
	return false
}

// TagsPrintable defines how tags are printed in options structs, it ensures the same order of tags for every field
func (f *Field) TagsPrintable() string {
	var tagNames = []string{"ddl", "sql", "db"}
	var tagParts []string
	for _, tagName := range tagNames {
		var v, ok = f.Tags[tagName]
		if ok {
			tagParts = append(tagParts, fmt.Sprintf(`%s:"%s"`, tagName, strings.Join(v, ",")))
		}
	}
	if len(tagParts) > 0 {
		return fmt.Sprintf("`%s`", strings.Join(tagParts, " "))
	}
	return ""
}

// KindNoPtr return field's Kind but without pointer
func (f *Field) KindNoPtr() string {
	kindWithoutPtr, _ := strings.CutPrefix(f.Kind, "*")
	return kindWithoutPtr
}

// IsStruct checks if field is a struct
func (f *Field) IsStruct() bool {
	return len(f.Fields) > 0
}

// ShouldBeInDto checks if field is not some static SQL field which should not be interacted with by SDK user
// TODO: this is a very naive implementation, consider fixing it with DSL builder connection
func (f *Field) ShouldBeInDto() bool {
	return !slices.Contains(f.Tags["ddl"], "static")
}

// IsRoot checks if field is at the top of field hierarchy, basically it is true for Option structs
func (f *Field) IsRoot() bool {
	return f.Parent == nil
}

// Path returns the way through the tree to the top, with dot separator (e.g. .SomeField.SomeChild)
func (f *Field) Path() string {
	if f.IsRoot() {
		return ""
	} else {
		return fmt.Sprintf("%s.%s", f.Parent.Path(), f.Name)
	}
}

// DtoKind returns what should be fields kind in generated DTO, because it may differ from Kind
func (f *Field) DtoKind() string {
	if f.IsRoot() {
		withoutSuffix, _ := strings.CutSuffix(f.Kind, "Options")
		return fmt.Sprintf("%sRequest", withoutSuffix)
	} else if f.IsStruct() {
		return fmt.Sprintf("%sRequest", f.Kind)
	} else {
		return f.Kind
	}
}

// DtoDecl returns how struct should be declared in generated DTO (e.g. definition is without a pointer)
func (f *Field) DtoDecl() string {
	if f.Parent == nil {
		withoutSuffix, _ := strings.CutSuffix(f.KindNoPtr(), "Options")
		return fmt.Sprintf("%sRequest", withoutSuffix)
	} else if f.IsStruct() {
		return fmt.Sprintf("%sRequest", f.KindNoPtr())
	} else {
		return f.KindNoPtr()
	}
}
