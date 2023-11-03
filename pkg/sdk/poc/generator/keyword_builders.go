package generator

func (v *QueryStruct) OptionalSQL(sql string) *QueryStruct {
	v.fields = append(v.fields, NewField(sqlToFieldName(sql, true), "*bool", Tags().Keyword().SQL(sql), nil))
	return v
}

func (v *QueryStruct) OrReplace() *QueryStruct {
	return v.OptionalSQL("OR REPLACE")
}

func (v *QueryStruct) IfNotExists() *QueryStruct {
	return v.OptionalSQL("IF NOT EXISTS")
}

func (v *QueryStruct) IfExists() *QueryStruct {
	return v.OptionalSQL("IF EXISTS")
}

func (v *QueryStruct) Terse() *QueryStruct {
	return v.OptionalSQL("TERSE")
}

func (v *QueryStruct) Text(name string, transformer *KeywordTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, "string", Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) Number(name string, transformer *KeywordTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, "int", Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) OptionalText(name string, transformer *KeywordTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, "*string", Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) OptionalNumber(name string, transformer *KeywordTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, "*int", Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) OptionalLimitFrom() *QueryStruct {
	v.fields = append(v.fields, NewField("Limit", "*LimitFrom", Tags().Keyword().SQL("LIMIT"), nil))
	return v
}

// SessionParameters *SessionParameters `ddl:"list,no_parentheses"`
func (v *QueryStruct) SessionParameters() *QueryStruct {
	v.fields = append(v.fields, NewField("SessionParameters", "*SessionParameters", Tags().List().NoParentheses(), nil).withValidations(NewValidation(ValidateValue, "SessionParameters")))
	return v
}

func (v *QueryStruct) OptionalSessionParameters() *QueryStruct {
	v.fields = append(v.fields, NewField("SessionParameters", "*SessionParameters", Tags().List().NoParentheses(), nil).withValidations(NewValidation(ValidateValue, "SessionParameters")))
	return v
}

func (v *QueryStruct) OptionalSessionParametersUnset() *QueryStruct {
	v.fields = append(v.fields, NewField("SessionParametersUnset", "*SessionParametersUnset", Tags().List().NoParentheses(), nil).withValidations(NewValidation(ValidateValue, "SessionParametersUnset")))
	return v
}

func (v *QueryStruct) List(sqlPrefix string, listItemKind string, transformer *ListTransformer) *QueryStruct {
	if transformer != nil {
		transformer = transformer.Parentheses().SQL(sqlPrefix)
	} else {
		transformer = ListOptions().Parentheses().SQL(sqlPrefix)
	}
	v.fields = append(v.fields, NewField(sqlToFieldName(sqlPrefix, true), KindOfSlice(listItemKind), Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) WithTags() *QueryStruct {
	v.List("TAG", "TagAssociation", nil)
	return v
}

func (v *QueryStruct) SetTags() *QueryStruct {
	v.fields = append(v.fields, NewField("SetTags", "[]TagAssociation", Tags().Keyword().SQL("SET TAG"), nil))
	return v
}

func (v *QueryStruct) UnsetTags() *QueryStruct {
	v.fields = append(v.fields, NewField("UnsetTags", "[]ObjectIdentifier", Tags().Keyword().SQL("UNSET TAG"), nil))
	return v
}

func (v *QueryStruct) OptionalLike() *QueryStruct {
	v.fields = append(v.fields, NewField("Like", "*Like", Tags().Keyword().SQL("LIKE"), nil))
	return v
}

func (v *QueryStruct) OptionalIn() *QueryStruct {
	v.fields = append(v.fields, NewField("In", "*In", Tags().Keyword().SQL("IN"), nil))
	return v
}

func (v *QueryStruct) OptionalStartsWith() *QueryStruct {
	v.fields = append(v.fields, NewField("StartsWith", "*string", Tags().Parameter().NoEquals().SingleQuotes().SQL("STARTS WITH"), nil))
	return v
}

func (v *QueryStruct) OptionalLimit() *QueryStruct {
	v.fields = append(v.fields, NewField("Limit", "*LimitFrom", Tags().Keyword().SQL("LIMIT"), nil))
	return v
}

func (v *QueryStruct) OptionalCopyGrants() *QueryStruct {
	return v.OptionalSQL("COPY GRANTS")
}
