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
