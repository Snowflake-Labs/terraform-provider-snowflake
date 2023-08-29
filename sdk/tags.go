package sdk

// placeholder for the real implementation.
type TagCreateOptions struct{}

type Tag struct {
	DatabaseName string
	SchemaName   string
	Name         string
}

func (v *Tag) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Tag) ObjectType() ObjectType {
	return ObjectTypeTag
}
