package generator2

import (
	"fmt"
	"strings"
)

type Struct struct {
	Name   string
	Fields []*Field
	// Validations defines validations on given field level (e.g. oneOf for children)
	Validations []*Validation
}

func NewStruct(name string) *Struct {
	return &Struct{
		Name:        name,
		Fields:      make([]*Field, 0),
		Validations: make([]*Validation, 0),
	}
}

func (s *Struct) WithFields(intoFields ...*Field) *Struct {
	// TODO Can be converted to IntoField
	//fields := make([]*Field, len(intoFields))
	//for i, f := range intoFields {
	//	fields[i] = f.IntoField()
	//}
	s.Fields = intoFields
	return s
}

func (s *Struct) WithValidations(validations ...*Validation) *Struct {
	s.Validations = validations
	return s
}

func (s *Struct) DtoDecl() string {
	withoutSuffix, _ := strings.CutSuffix(s.Name, "Options")
	return fmt.Sprintf("%sRequest", withoutSuffix)
	//if field.Parent == nil {
	//	withoutSuffix, _ := strings.CutSuffix(field.KindNoPtr(), "Options")
	//	return fmt.Sprintf("%sRequest", withoutSuffix)
	//} else if field.IsStruct() {
	//	return fmt.Sprintf("%sRequest", field.KindNoPtr())
	//} else {
	//	return field.KindNoPtr()
	//}
}
