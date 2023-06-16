package sdk

// placeholder for the real implementation.
type TableCreateOptions struct{}

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
