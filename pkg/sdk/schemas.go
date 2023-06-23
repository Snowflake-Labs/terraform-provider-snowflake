package sdk

// placeholder for the real implementation.
type Schema struct {
	DatabaseName string
	Name         string
}

func (v *Schema) ID() SchemaIdentifier {
	return NewSchemaIdentifier(v.DatabaseName, v.Name)
}

func (v *Schema) ObjectType() ObjectType {
	return ObjectTypeSchema
}
