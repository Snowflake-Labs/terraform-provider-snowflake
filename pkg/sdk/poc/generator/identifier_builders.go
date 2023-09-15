package generator

func (v *queryStruct) Identifier(fieldName string, kind string) *queryStruct {
	v.fields = append(v.fields, NewField(fieldName, kind, Tags().Identifier(), nil))
	return v
}

// func AccountObjectIdentifier(fieldName string) *Field {
//	return NewField(fieldName, "AccountObjectIdentifier", Tags().Identifier()).WithRequired(true)
//}
//
//func DatabaseObjectIdentifier(fieldName string) *Field {
//	return NewField(fieldName, "DatabaseObjectIdentifier", Tags().Identifier()).WithRequired(true)
//}
