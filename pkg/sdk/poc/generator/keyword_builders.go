package generator

func (v *queryStruct) OptionalSQL(sql string) *queryStruct {
	v.fields = append(v.fields, NewField(sqlToFieldName(sql, true), "*bool", Tags().Keyword().SQL(sql), nil))
	return v
}

func (v *queryStruct) OrReplace() *queryStruct {
	return v.OptionalSQL("OR REPLACE")
}

func (v *queryStruct) IfNotExists() *queryStruct {
	return v.OptionalSQL("IF NOT EXISTS")
}

func (v *queryStruct) IfExists() *queryStruct {
	return v.OptionalSQL("IF EXISTS")
}

func (v *queryStruct) Text(name string, transformer *KeywordTransformer) *queryStruct {
	v.fields = append(v.fields, NewField(name, "string", Tags().Keyword(), transformer))
	return v
}

func (v *queryStruct) OptionalText(name string, transformer *KeywordTransformer) *queryStruct {
	v.fields = append(v.fields, NewField(name, "*string", Tags().Keyword(), transformer))
	return v
}

// SessionParameters *SessionParameters `ddl:"list,no_parentheses"`
func (v *queryStruct) SessionParameters() *queryStruct {
	v.fields = append(v.fields, NewField("SessionParameters", "*SessionParameters", Tags().List().NoParentheses(), nil).withValidations(&Validation{Type: ValidateValue, FieldNames: []string{"SessionParameters"}}))
	return v
}

func (v *queryStruct) OptionalSessionParameters() *queryStruct {
	v.fields = append(v.fields, NewField("SessionParameters", "*SessionParameters", Tags().List().NoParentheses(), nil).withValidations(&Validation{Type: ValidateValue, FieldNames: []string{"SessionParameters"}}))
	return v
}

func (v *queryStruct) WithTags() *queryStruct {
	v.fields = append(v.fields, NewField("Tag", "[]TagAssociation", Tags().Keyword().Parentheses().SQL("TAG"), nil))
	return v
}

func (v *queryStruct) SetTags() *queryStruct {
	v.fields = append(v.fields, NewField("SetTags", "[]TagAssociation", Tags().Keyword().SQL("SET TAG"), nil))
	return v
}

func (v *queryStruct) UnsetTags() *queryStruct {
	v.fields = append(v.fields, NewField("UnsetTags", "[]ObjectIdentifier", Tags().Keyword().SQL("UNSET TAG"), nil))
	return v
}
