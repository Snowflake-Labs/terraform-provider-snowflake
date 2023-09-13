package generator

func (f *Field) SQL(sql string) *Field {
	//f.Fields = append(f.Fields, NewField(sqlToFieldName(sql, false), "bool", Tags().Static().SQL(sql)))
	return f
}

func (f *Field) Create() *Field {
	return f.SQL("CREATE")
}

func (f *Field) Alter() *Field {
	return f.SQL("ALTER")
}

func (f *Field) Drop() *Field {
	return f.SQL("DROP")
}

func (f *Field) Show() *Field {
	return f.SQL("SHOW")
}

func (f *Field) Describe() *Field {
	return f.SQL("DESCRIBE")
}
