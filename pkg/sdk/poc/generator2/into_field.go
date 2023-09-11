package generator2

// TODO Could be used in WithFields(fields ...Into) where we could pass not only fields but also structs or enums
type Into interface {
	IntoField
	IntoStruct
	//IntoEnum
}

type IntoField interface {
	IntoField() *Field
}

type IntoStruct interface {
	IntoStruct() *Struct
}
