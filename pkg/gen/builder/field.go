package builder

func (fb FieldBuilder) IntoFieldBuilder() []FieldBuilder {
	return []FieldBuilder{fb}
}

func Keyword(fieldName string, typer Typer) FieldBuilder {
	return FieldBuilder{
		Name:  fieldName,
		Typer: typer,
		Tags: map[string][]string{
			"ddl": {"keyword"},
			"sql": {},
		},
	}
}

func OptionalSQL(sql string) FieldBuilder {
	kw := Keyword(sqlToFieldName(sql, true), TypeBoolPtr)
	kw.Tags["sql"] = append(kw.Tags["sql"], sql)
	return kw
}

func OptionalText(fieldName string, keywordOpts FieldTransformer) FieldBuilder {
	kw := Keyword(fieldName, TypeStringPtr)
	fb := keywordOpts.transform(&kw)
	return *fb
}

func OptionalValue(fieldName string, typer Typer, keywordOpts FieldTransformer) FieldBuilder {
	kw := Keyword(fieldName, typer)
	fb := keywordOpts.transform(&kw)
	return *fb
}
