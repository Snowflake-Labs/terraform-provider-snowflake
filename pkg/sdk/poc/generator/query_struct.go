package generator

type queryStruct struct {
	name        string
	fields      []*Field
	validations []*Validation
}

func QueryStruct(name string) *queryStruct {
	return &queryStruct{
		name:        name,
		fields:      make([]*Field, 0),
		validations: make([]*Validation, 0),
	}
}

func (v *queryStruct) IntoField() *Field {
	return NewField(v.name, v.name, nil, nil).
		withFields(v.fields...).
		withValidations(v.validations...)
}

func (v *queryStruct) WithValidation(validationType ValidationType, fieldNames ...string) *queryStruct {
	v.validations = append(v.validations, NewValidation(validationType, fieldNames...))
	return v
}

func (v *queryStruct) QueryStructField(queryStruct *queryStruct) *queryStruct {
	v.fields = append(v.fields, queryStruct.IntoField())
	return v
}
