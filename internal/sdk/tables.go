// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

// placeholder for the real implementation.
type CreateTableOptions struct{}

type Table struct {
	DatabaseName string
	SchemaName   string
	Name         string
}

func (v *Table) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Table) ObjectType() ObjectType {
	return ObjectTypeTable
}
