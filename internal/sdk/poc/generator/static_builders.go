// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package generator

func (v *queryStruct) SQL(sql string) *queryStruct {
	v.fields = append(v.fields, NewField(sqlToFieldName(sql, false), "bool", Tags().Static().SQL(sql), nil))
	return v
}

func (v *queryStruct) Create() *queryStruct {
	return v.SQL("CREATE")
}

func (v *queryStruct) Alter() *queryStruct {
	return v.SQL("ALTER")
}

func (v *queryStruct) Drop() *queryStruct {
	return v.SQL("DROP")
}

func (v *queryStruct) Show() *queryStruct {
	return v.SQL("SHOW")
}

func (v *queryStruct) Describe() *queryStruct {
	return v.SQL("DESCRIBE")
}
