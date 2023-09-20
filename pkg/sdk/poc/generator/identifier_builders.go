package generator

// SelfIdentifier
func (v *queryStruct) SelfIdentifier() *queryStruct {
	identifier := NewField("name", "<will be replaced>", Tags().Identifier(), IdentifierOptions().Required())
	v.identifierField = identifier
	v.fields = append(v.fields, identifier)
	return v
}

func (v *queryStruct) Identifier(fieldName string, kind string, transformer *IdentifierTransformer) *queryStruct {
	v.fields = append(v.fields, NewField(fieldName, kind, Tags().Identifier(), transformer))
	return v
}

// func AccountObjectIdentifier(fieldName string) *Field {
//	return NewField(fieldName, "AccountObjectIdentifier", Tags().Identifier()).WithRequired(true)
//}
//
//func DatabaseObjectIdentifier(fieldName string) *Field {
//	return NewField(fieldName, "DatabaseObjectIdentifier", Tags().Identifier()).WithRequired(true)
//}
