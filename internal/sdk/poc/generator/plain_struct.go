// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package generator

type plainStruct struct {
	name   string
	fields []plainField
}

type plainField struct {
	name string
	kind string
}

func PlainStruct(name string) *plainStruct {
	return &plainStruct{
		name:   name,
		fields: make([]plainField, 0),
	}
}

func (v *plainStruct) Field(name string, kind string) *plainStruct {
	v.fields = append(v.fields, plainField{
		name: name,
		kind: kind,
	})
	return v
}

func (v *plainStruct) IntoField() *Field {
	f := NewField(v.name, v.name, nil, nil)
	for _, field := range v.fields {
		f.withField(NewField(field.name, field.kind, nil, nil))
	}
	return f
}
