package generator

type queryStruct struct {
	name            string
	fields          []*Field
	identifierField *Field
	validations     []*Validation
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

func (v *queryStruct) QueryStructField(name string, queryStruct *queryStruct, transformer FieldTransformer) *queryStruct {
	qs := queryStruct.IntoField()
	qs.Name = name
	qs = transformer.Transform(qs)
	v.fields = append(v.fields, qs)
	return v
}

func (v *queryStruct) ListQueryStructField(name string, queryStruct *queryStruct, transformer FieldTransformer) *queryStruct {
	qs := queryStruct.IntoField()
	qs.Name = name
	qs.Kind = "[]" + qs.Kind
	qs = transformer.Transform(qs)
	v.fields = append(v.fields, qs)
	return v
}

func (v *queryStruct) OptionalQueryStructField(name string, queryStruct *queryStruct, transformer FieldTransformer) *queryStruct {
	qs := queryStruct.IntoField()
	qs.Name = name
	qs.Kind = "*" + qs.Kind
	qs = transformer.Transform(qs)
	v.fields = append(v.fields, qs)
	return v
}
