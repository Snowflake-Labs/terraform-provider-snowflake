package generator

func (v *QueryStruct) ListWithParens(sqlPrefix string, listItemKind string, transformer *ListTransformer) *QueryStruct {
	if transformer != nil {
		transformer = transformer.Parentheses().SQL(sqlPrefix)
	} else {
		transformer = ListOptions().Parentheses().SQL(sqlPrefix)
	}
	v.fields = append(v.fields, NewField(sqlToFieldName(sqlPrefix, true), KindOfSlice(listItemKind), Tags().List(), transformer))
	return v
}
