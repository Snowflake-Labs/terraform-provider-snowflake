package generator

func (v *queryStruct) assignment(name string, kind string, transformer *ParameterTransformer) *queryStruct {
	v.fields = append(v.fields, NewField(name, kind, Tags().Parameter(), transformer))
	return v
}

func (v *queryStruct) Assignment(sqlPrefix string, kind string, transformer *ParameterTransformer) *queryStruct {
	if transformer != nil {
		transformer = transformer.SQL(sqlPrefix)
	} else {
		transformer = ParameterOptions().SQL(sqlPrefix)
	}
	return v.assignment(sqlToFieldName(sqlPrefix, true), kind, transformer)
}

func (v *queryStruct) OptionalAssignment(sqlPrefix string, kind string, transformer *ParameterTransformer) *queryStruct {
	if len(kind) > 0 && kind[0] != '*' {
		kind = KindOfPointer(kind)
	}
	return v.Assignment(sqlPrefix, kind, transformer)
}

func (v *queryStruct) ListAssignment(sqlPrefix string, listItemKind string, transformer *ParameterTransformer) *queryStruct {
	return v.Assignment(sqlPrefix, KindOfSlice(listItemKind), transformer)
}

func (v *queryStruct) NumberAssignment(sqlPrefix string, transformer *ParameterTransformer) *queryStruct {
	return v.Assignment(sqlPrefix, "int", transformer)
}

func (v *queryStruct) OptionalNumberAssignment(sqlPrefix string, transformer *ParameterTransformer) *queryStruct {
	return v.Assignment(sqlPrefix, "*int", transformer)
}

func (v *queryStruct) TextAssignment(sqlPrefix string, transformer *ParameterTransformer) *queryStruct {
	return v.Assignment(sqlPrefix, "string", transformer)
}

func (v *queryStruct) OptionalTextAssignment(sqlPrefix string, transformer *ParameterTransformer) *queryStruct {
	return v.Assignment(sqlPrefix, "*string", transformer)
}

func (v *queryStruct) BooleanAssignment(sqlPrefix string, transformer *ParameterTransformer) *queryStruct {
	return v.Assignment(sqlPrefix, "bool", transformer)
}

func (v *queryStruct) OptionalBooleanAssignment(sqlPrefix string, transformer *ParameterTransformer) *queryStruct {
	return v.Assignment(sqlPrefix, "*bool", transformer)
}

func (v *queryStruct) OptionalIdentifierAssignment(sqlPrefix string, identifierKind string, transformer *ParameterTransformer) *queryStruct {
	return v.OptionalAssignment(sqlPrefix, identifierKind, transformer)
}

func (v *queryStruct) OptionalComment() *queryStruct {
	return v.OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes())
}

func (v *queryStruct) SetComment() *queryStruct {
	return v.OptionalTextAssignment("SET COMMENT", ParameterOptions().SingleQuotes())
}
