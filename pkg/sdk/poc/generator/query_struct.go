package generator

// TODO For Field abstractions use internal Field representation instead of copying only needed fields, e.g.
//
//	type queryStruct struct {
//		internalRepresentation *Field
//		...additional fields that are not present in the Field
//	}
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
	return v.queryStructField(name, queryStruct, "", transformer)
}

func (v *queryStruct) ListQueryStructField(name string, queryStruct *queryStruct, transformer FieldTransformer) *queryStruct {
	return v.queryStructField(name, queryStruct, "[]", transformer)
}

func (v *queryStruct) OptionalQueryStructField(name string, queryStruct *queryStruct, transformer FieldTransformer) *queryStruct {
	return v.queryStructField(name, queryStruct, "*", transformer)
}

func (v *queryStruct) queryStructField(name string, queryStruct *queryStruct, kindPrefix string, transformer FieldTransformer) *queryStruct {
	qs := queryStruct.IntoField()
	qs.Name = name
	qs.Kind = kindPrefix + qs.Kind
	qs = transformer.Transform(qs)
	v.fields = append(v.fields, qs)
	return v
}
