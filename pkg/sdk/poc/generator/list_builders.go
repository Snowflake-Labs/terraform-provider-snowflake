package generator

func (v *queryStruct) List(name string, itemKind string, transformer FieldTransformer) *queryStruct {
	v.fields = append(v.fields, NewField(name, KindOfSlice(itemKind), Tags().List(), transformer))
	return v
}
