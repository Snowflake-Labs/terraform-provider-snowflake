package generator

func (f *Field) Identifier(fieldName string, kind string) *Field {
	f.Fields = append(f.Fields, NewField(fieldName, kind, Tags().Identifier(), nil))
	return f
}

//func AccountObjectIdentifier(fieldName string) *Field {
//	return NewField(fieldName, "AccountObjectIdentifier", Tags().Identifier()).WithRequired(true)
//}
//
//func DatabaseObjectIdentifier(fieldName string) *Field {
//	return NewField(fieldName, "DatabaseObjectIdentifier", Tags().Identifier()).WithRequired(true)
//}
