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
	// TODO Add path somewhere to keep track in validations, dtos, etc
}

func NewStruct(name string) *Struct {
	return &Struct{
		Name:        name,
		Fields:      make([]*Field, 0),
		Validations: make([]*Validation, 0),
	}
}

func (s *Struct) WithFields(fields ...*Field) *Struct {
	//fields := make([]*Field, len(into))
	//for i, f := range into {
	//	fields[i] = f.IntoField()
	//}
	//s.Fields = append(s.Fields, fields...)
	s.Fields = fields
	return s
}

func (s *Struct) WithValidations(validations ...*Validation) *Struct {
	s.Validations = validations
	return s
}

func (s *Struct) DtoDecl() string {
	withoutSuffix, _ := strings.CutSuffix(s.Name, "Options")
	return fmt.Sprintf("%sRequest", withoutSuffix)
}
