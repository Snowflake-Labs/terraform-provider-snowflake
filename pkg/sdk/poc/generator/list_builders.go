package generator

func (v *QueryStruct) List(name string, itemKind string, transformer FieldTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, KindOfSlice(itemKind), Tags().List(), transformer))
	return v
}
