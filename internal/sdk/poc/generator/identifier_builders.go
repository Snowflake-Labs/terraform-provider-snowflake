// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package generator

// Name adds identifier with field name "name" and type will be inferred from interface definition
func (v *queryStruct) Name() *queryStruct {
	identifier := NewField("name", "<will be replaced>", Tags().Identifier(), IdentifierOptions().Required())
	v.identifierField = identifier
	v.fields = append(v.fields, identifier)
	return v
}

func (v *queryStruct) Identifier(fieldName string, kind string, transformer *IdentifierTransformer) *queryStruct {
	v.fields = append(v.fields, NewField(fieldName, kind, Tags().Identifier(), transformer))
	return v
}

func (v *queryStruct) OptionalIdentifier(name string, kind string, transformer *IdentifierTransformer) *queryStruct {
	if len(kind) > 0 && kind[0] != '*' {
		kind = KindOfPointer(kind)
	}
	v.fields = append(v.fields, NewField(name, kind, Tags().Identifier(), transformer))
	return v
}
