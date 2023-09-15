package generator

func (v *queryStruct) Assignment(name string, kind string, transformer *ParameterTransformer) *queryStruct {
	v.fields = append(v.fields, NewField(name, kind, Tags().Parameter(), transformer))
	return v
}

func (v *queryStruct) ListAssignment(sqlPrefix string, listItemKind string, transformer *ParameterTransformer) *queryStruct {
	if transformer != nil {
		transformer = transformer.SQL(sqlPrefix)
	} else {
		transformer = ParameterOptions().SQL(sqlPrefix)
	}
	return v.Assignment(sqlToFieldName(sqlPrefix, true), KindOfSlice(listItemKind), transformer)
}

func (v *queryStruct) TextAssignment(sqlPrefix string, transformer *ParameterTransformer) *queryStruct {
	if transformer != nil {
		transformer = transformer.SQL(sqlPrefix)
	} else {
		transformer = ParameterOptions().SQL(sqlPrefix)
	}
	return v.Assignment(sqlToFieldName(sqlPrefix, true), "string", transformer)
}

func (v *queryStruct) OptionalTextAssignment(sqlPrefix string, transformer *ParameterTransformer) *queryStruct {
	if transformer != nil {
		transformer = transformer.SQL(sqlPrefix)
	} else {
		transformer = ParameterOptions().SQL(sqlPrefix)
	}
	return v.Assignment(sqlToFieldName(sqlPrefix, true), "*string", transformer)
}
