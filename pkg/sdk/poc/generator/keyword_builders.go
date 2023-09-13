package generator

func (f *Field) OptionalSQL(sql string) *Field {
	f.Fields = append(f.Fields, NewField(sqlToFieldName(sql, true), "*bool", Tags().Keyword().SQL(sql), nil))
	return f
}

func (f *Field) OrReplace() *Field {
	return f.OptionalSQL("OR REPLACE")
}

func (f *Field) IfNotExists() *Field {
	return f.OptionalSQL("IF NOT EXISTS")
}

func (f *Field) IfExists() *Field {
	return f.OptionalSQL("IF EXISTS")
}

func (f *Field) Text(name string, transformer *KeywordTransformer) *Field {
	f.Fields = append(f.Fields, NewField(name, "string", Tags().Keyword(), transformer))
	return f
}
